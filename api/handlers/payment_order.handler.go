package handlers

import (
	"leichil-backend/api/models"
	"leichil-backend/api/utils"

	"github.com/gofiber/fiber/v2"
)

var payment_verification_required_fields = []string{"RazorpayOrderID", "RazorpayPaymentID", "RazorpaySignature"}
var payment_update_required_fields = []string{"Status", "OrderID"}
var order_required_fields = []string{"ContentID"}

func app_order_handler(router fiber.Router) {
	router.Post("/", add_payment_order)
	router.Post("/verify", verify_payment)
	router.Put("/", update_payment)
	router.Get("/purchases", get_customer_payment_orders)
	router.Get("/purchases/:id", verify_customer_purchase)
}

func get_payment_orders(ctx *fiber.Ctx) error {
	var payment models.PaymentOrder
	var offset int = utils.GetIntQueryParams(ctx.Query("offset"), 0)
	var limit int = utils.GetIntQueryParams(ctx.Query("limit"), 10)
	var params = map[string]string{
		"order_id":    ctx.Query("order_id"),
		"payment_id":  ctx.Query("payment_id"),
		"amount":      ctx.Query("amount"),
		"status":      ctx.Query("content"),
		"content_id":  ctx.Query("content_id"),
		"customer_id": ctx.Query("customer_id"),
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

func add_payment_order(ctx *fiber.Ctx) error {
	var order = new(models.PaymentOrder)
	flag, status, message := utils.ParseAndValidate(order, order_required_fields, ctx)
	if flag {
		var token = ctx.Locals("customer_token").(models.CustomerToken)
		var result = order.Create(token)
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
	flag, status, message := utils.ParseAndValidate(payment, payment_verification_required_fields, ctx)
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
}
