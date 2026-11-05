package handlers

import (
	"encoding/base64"
	"fmt"
	"leichil-backend/api/models"
	"leichil-backend/api/utils"
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
var qr_required_fields = []string{"DeviceID", "PaymentType"}

func order_handler(router fiber.Router) {
	router.Get("/", get_payment_orders)
}

func app_order_handler(router fiber.Router) {
	router.Post("/", func(ctx *fiber.Ctx) error { return add_payment_order(ctx, "payment") })
	router.Post("/subscriptions", func(ctx *fiber.Ctx) error { return add_payment_order(ctx, "subscription") })
	router.Post("/verify", func(ctx *fiber.Ctx) error { return verify_payment(ctx, "payment") })
	router.Post("/subscriptions/verify", func(ctx *fiber.Ctx) error { return verify_payment(ctx, "subscription") })
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

func verify_payment(ctx *fiber.Ctx, order_type string) error {
	var payment = new(models.VerifyPaymentBody)
	var required_fields = payment_verification_required_fields
	if order_type == "subscription" {
		required_fields = append(required_fields, "PlanID")
	}
	flag, status, message := utils.ParseAndValidate(payment, required_fields, ctx)
	if flag {
		var response = map[string]any{
			"success": false,
			"message": "Payment Failed",
		}
		var token = ctx.Locals("customer_token").(models.CustomerToken)
		if payment.VerifyPayment(token, order_type) {
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
		var payment_type string = strings.ToUpper(models.SanitizeData(body.PaymentType))
		if payment_type == "PURCHASE" {
			var content_id = models.SanitizeData(body.ContentID)
			if content_id > 0 {
				var url string = os.Getenv("PAYMENT_URL")
				var device_id string = models.SanitizeData(body.DeviceID)
				var token = ctx.Locals("customer_token").(models.CustomerToken)
				var customer_id int = models.SanitizeData(token.CustomerID)
				var timestamp = time.Now().UnixMilli()
				var query_params = fmt.Sprintf(`content_id=%d&customer_id=%d&device_id=%s&timestamp=%d`, content_id, customer_id, device_id, timestamp)
				encrypted_params, err := utils.Encrypt(query_params)
				if err != nil {
					return utils.GetMessageResponse("", "Unable to Get Payment Link", fiber.StatusInternalServerError, false, ctx)
				}
				url = fmt.Sprintf(`%s?params=%s`, url, encrypted_params)
				return gen_qrcode(url, ctx)
			}
			return utils.GetMessageResponse("", "MISSING CONTENT ID", fiber.StatusBadRequest, false, ctx)
		} else if payment_type == "SUBSCRIPTION" {
			if models.SanitizeData(body.PlanID) > 0 {
				var url string = os.Getenv("PAYMENT_URL")
				return gen_qrcode(url, ctx)
			}
			return utils.GetMessageResponse("", "MISSING PLAN ID", fiber.StatusBadRequest, false, ctx)
		}
		status = fiber.StatusBadRequest
		flag = false
		message = "INVALID PAYMENT TYPE"
	}
	return utils.GetMessageResponse(message, message, status, flag, ctx)
}

func gen_qrcode(url string, ctx *fiber.Ctx) error {
	png, err := qrcode.Encode(url, qrcode.Medium, 256)
	if err != nil {
		log.Error("Error in generating QR Code: ", err.Error())
		return utils.GetMessageResponse("", "Unable to generate QR CODE", fiber.StatusInternalServerError, false, ctx)
	}
	var base64_str = base64.StdEncoding.EncodeToString(png)
	var img = fmt.Sprintf("data:image/png;base64,%s", base64_str)
	var response = map[string]any{
		"success": true,
		"result":  img,
	}
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func HandlePaymentLink(ctx *fiber.Ctx) error {
	var params = ctx.Query("params")
	log.Info(params)
	data, err := utils.Decrypt(params)
	if err != nil {
		log.Error("Error in Decryption: ", err.Error())
		return ctx.Status(fiber.StatusInternalServerError).SendString("Server Error")
	}
	return ctx.Status(fiber.StatusOK).SendString(data)
}
