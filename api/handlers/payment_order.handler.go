package handlers

import (
	"leichil-backend/api/models"
	"leichil-backend/api/utils"

	"github.com/gofiber/fiber/v2"
)

var payment_verification_required_fields = []string{"RazorpayOrderID", "RazorpayPaymentID", "RazorpaySignature"}
var payment_update_required_fields = []string{"Status", "OrderID"}
var order_required_fields = []string{"ContentID"}
var subscription_order_required_fields = []string{"PlanID"}

func app_order_handler(router fiber.Router) {
	router.Post("/", func(ctx *fiber.Ctx) error { return add_payment_order(ctx, "payment") })
	router.Post("/subscriptions", func(ctx *fiber.Ctx) error { return add_payment_order(ctx, "subscription") })
	router.Post("/verify", func(ctx *fiber.Ctx) error { return verify_payment(ctx, "payment") })
	router.Post("/subscriptions/verify", func(ctx *fiber.Ctx) error { return verify_payment(ctx, "subscription") })
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
