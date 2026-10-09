package handlers

import (
	"leichil-backend/api/models"
	"leichil-backend/api/utils"

	"github.com/gofiber/fiber/v2"
)

var payment_verification_required_fields = []string{"RazorpayOrderID", "RazorpayPaymentID", "RazorpaySignature"}
var payment_update_required_fields = []string{"Status", "OrderID"}

func app_order_handler(router fiber.Router) {
	router.Post("/", add_payment_order)
	router.Post("/verify", verify_payment)
	router.Put("/", update_payment)
}

func add_payment_order(ctx *fiber.Ctx) error {
	var order models.PaymentOrder
	var token = ctx.Locals("customer_token").(models.CustomerToken)
	var result = order.Create(token)
	var response = map[string]any{
		"success": true,
		"result":  result,
	}
	return ctx.Status(fiber.StatusOK).JSON(response)
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
