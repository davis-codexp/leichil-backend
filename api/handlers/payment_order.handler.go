package handlers

import (
	"encoding/base64"
	"fmt"
	"leichil-backend/api/models"
	"leichil-backend/api/templates"
	"leichil-backend/api/utils"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/skip2/go-qrcode"
)

var payment_verification_required_fields = []string{"RazorpayOrderID", "RazorpayPaymentID", "RazorpaySignature"}
var payment_update_required_fields = []string{"Status", "OrderID"}
var order_required_fields = []string{"ContentID"}
var subscription_order_required_fields = []string{"PlanID"}
var qr_required_fields = []string{"DeviceID", "ContentID"}

func order_handler(router fiber.Router) {
	router.Get("/", get_payment_orders)
}

func app_order_handler(router fiber.Router) {
	router.Post("/", func(ctx *fiber.Ctx) error { return add_payment_order(ctx, "payment") })
	router.Post("/subscriptions", func(ctx *fiber.Ctx) error { return add_payment_order(ctx, "subscription") })
	router.Post("/verify", verify_payment)
	router.Post("/subscriptions/verify", verify_payment)
	router.Put("/", update_payment)
	router.Get("/purchases", get_customer_payment_orders)
	router.Get("/purchases/:id", verify_customer_purchase)
	router.Post("/qr", get_payment_qr)
}

func get_payment_orders(ctx *fiber.Ctx) error {
	var payment models.PaymentOrder
	var offset int = utils.GetIntQueryParams(ctx.Query("offset"), 0)
	var limit int = utils.GetIntQueryParams(ctx.Query("limit"), 10)
	var params = map[string]string{
		"order_id":    ctx.Query("order_id"),
		"payment_id":  ctx.Query("payment_id"),
		"amount":      ctx.Query("amount"),
		"status":      ctx.Query("status"),
		"content_id":  ctx.Query("content_id"),
		"customer_id": ctx.Query("customer_id"),
		"start_date":  ctx.Query("start_date"),
		"end_date":    ctx.Query("end_date"),
	}
	_, result, err := payment.Get(offset, limit, params)
	return utils.GetArrayResponse(ctx, 0, result, "Unable to get Payment Orders", err)
}

func get_customer_payment_orders(ctx *fiber.Ctx) error {
	var payment models.PaymentOrder
	var token = ctx.Locals("customer_token").(models.CustomerToken)
	_, result, err := payment.GetCustomerPurchase(token)
	return utils.GetArrayResponse(ctx, 0, result, "Unable to get Customer Purchases", err)
}

func verify_customer_purchase(ctx *fiber.Ctx) error {
	var payment models.PaymentOrder
	var token = ctx.Locals("customer_token").(models.CustomerToken)
	var flag = payment.VerifyCustomerPurchase(token, ctx.Params("id"))
	var response = map[string]any{
		"success": flag,
	}
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func add_payment_order(ctx *fiber.Ctx, order_type string) error {
	var order = new(models.PaymentOrder)
	var required_fields = order_required_fields
	if order_type == "subscription" {
		required_fields = subscription_order_required_fields
	}
	flag, status, message := utils.ParseAndValidate(order, required_fields, ctx)
	if flag {
		var token = ctx.Locals("customer_token").(models.CustomerToken)
		var result = order.Create(token, order_type)
		var response = map[string]any{
			"success": true,
			"result":  result,
		}
		return ctx.Status(fiber.StatusOK).JSON(response)
	}
	return utils.GetMessageResponse(message, message, status, flag, ctx)
}

func verify_payment(ctx *fiber.Ctx) error {
	var payment = new(models.VerifyPaymentBody)
	var required_fields = payment_verification_required_fields
	flag, status, message := utils.ParseAndValidate(payment, required_fields, ctx)
	if flag {
		var response = map[string]any{
			"success": false,
			"message": "Payment Failed",
		}
		if payment.VerifyPayment() {
			response["success"] = true
			response["message"] = "Payment Successful"
		}
		return ctx.Status(fiber.StatusOK).JSON(response)
	}
	return utils.GetMessageResponse(message, message, status, flag, ctx)
}

func update_payment(ctx *fiber.Ctx) error {
	var message = "Processing Payment"
	return utils.GetMessageResponse(message, message, fiber.StatusOK, true, ctx)
	/*
		var payment = new(models.PaymentOrder)
		flag, status, message := utils.ParseAndValidate(payment, payment_update_required_fields, ctx)
		if flag {
			if payment.Update() {
				flag = true
				message = "Payment Order Updated"
			} else {
				flag = false
				message = "Payment Order Update Failed"
			}
		}
		return utils.GetMessageResponse(message, message, status, flag, ctx)
	*/
}

func get_payment_qr(ctx *fiber.Ctx) error {
	var body = new(models.PaymentQRBody)
	var required_fields = qr_required_fields
	flag, status, message := utils.ParseAndValidate(body, required_fields, ctx)
	if flag {
		var content_id = models.SanitizeData(body.ContentID)
		var timestamp = time.Now().Add(5 * time.Minute).UnixMilli()
		var url string = os.Getenv("PAYMENT_URL")
		var token = ctx.Locals("customer_token").(models.CustomerToken)
		var customer_id int = models.SanitizeData(token.CustomerID)
		var device_id string = models.SanitizeData(body.DeviceID)
		var query_params = fmt.Sprintf(`content_id=%d&customer_id=%d&device_id=%s&timestamp=%d&token=%s`, content_id, customer_id, device_id, timestamp, models.SanitizeData(token.Token))
		encrypted_params, err := utils.Encrypt(query_params)
		if err != nil {
			return utils.GetMessageResponse("", "Unable to Get Payment Link", fiber.StatusInternalServerError, false, ctx)
		}
		url = fmt.Sprintf(`%s?params=%s`, url, encrypted_params)
		return gen_qrcode(url, timestamp, ctx)
	}
	return utils.GetMessageResponse(message, message, status, flag, ctx)
}

func gen_qrcode(url string, timestamp int64, ctx *fiber.Ctx) error {
	png, err := qrcode.Encode(url, qrcode.Medium, 256)
	if err != nil {
		log.Error("Error in generating QR Code: ", err.Error())
		return utils.GetMessageResponse("", "Unable to generate QR CODE", fiber.StatusInternalServerError, false, ctx)
	}
	var base64_str = base64.StdEncoding.EncodeToString(png)
	var img = fmt.Sprintf("data:image/png;base64,%s", base64_str)
	var response = map[string]any{
		"success":    true,
		"result":     img,
		"expires_at": timestamp,
	}
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func HandlePaymentLink(ctx *fiber.Ctx) error {
	var params = ctx.Query("params")
	data, err := utils.Decrypt(params)
	if err != nil {
		log.Error("Error in Decryption: ", err.Error())
		return ctx.Status(fiber.StatusInternalServerError).SendString("Server Error")
	}
	query_params, err := url.ParseQuery(data)
	if err != nil {
		log.Error("Unable to parse decrypted query params: ", err.Error())
		return utils.GetMessageResponse("", "Unable to genrate payment link", fiber.StatusInternalServerError, false, ctx)
	}
	var customer_id = utils.GetIntQueryParams(query_params.Get("customer_id"), 0)
	var content_id = utils.GetIntQueryParams(query_params.Get("content_id"), 0)
	var timestamp = utils.GetIntQueryParams(query_params.Get("timestamp"), 0)
	var token_str string = query_params.Get("token")
	var curr_timestamp = int(time.Now().UnixMilli())
	var order_type = "payment"
	var description string

	if curr_timestamp >= timestamp {
		log.Info("QR Expired")
		return ctx.Status(400).SendString("QR Code is Expired")
	}

	var customer models.Customer
	var content models.Content

	customer_result, err := customer.GetById(fmt.Sprintf(`%d`, customer_id))
	if err != nil {
		log.Error("Error in getting Customer: ", err.Error())
		return Render(ctx, templates.PaymentError())
	}
	customer = utils.FormatSingleResult[models.Customer](customer_result)

	content, err = content.GetById(fmt.Sprintf(`%d`, content_id), 1)
	if err != nil {
		log.Error("Error in getting Content: ", err.Error())
		return Render(ctx, templates.PaymentError())
	}

	if strings.ToUpper(models.SanitizeData(content.Category)) == "SUBSCRIPTION" {
		order_type = "subscription"
		description = "SUBSCRIPTION"
	}

	var order = new(models.PaymentOrder)
	order.ContentID = content.ID
	var token = models.CustomerToken{CustomerID: customer.ID}
	var result = order.Create(token, order_type)
	order_id, _ := result["order_id"].(string)
	amount, _ := result["amount"].(int)
	var api_key = os.Getenv("RZP_KEY_ID")
	return Render(
		ctx,
		templates.Checkout(
			order_id, api_key,
			models.SanitizeData(customer.Name),
			models.SanitizeData(customer.Mobile),
			description, token_str, amount,
		))
}
