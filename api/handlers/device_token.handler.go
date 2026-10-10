package handlers

import (
	"leichil-backend/api/models"
	"leichil-backend/api/utils"

	"github.com/gofiber/fiber/v2"
)

var device_token_required_fields = []string{"DeviceToken", "CustomerID"}

func app_device_token_handler(router fiber.Router) {
	router.Get("/", get_device_tokens)
	router.Post("/", add_device_token)
	router.Delete("/:device_token", delete_device_token)
}

func get_device_tokens(ctx *fiber.Ctx) error {
	var dt models.DeviceToken
	var token = ctx.Locals("customer_token").(models.CustomerToken)
	_, result, err := dt.Get(token)
	return utils.GetArrayResponse(ctx, 0, result, "Unable to get Device Tokens", err)
}

func add_device_token(ctx *fiber.Ctx) error {
	var dt = new(models.DeviceToken)
	flag, status, message := utils.ParseAndValidate(dt, device_token_required_fields, ctx)
	if flag {
		var token = ctx.Locals("customer_token").(models.CustomerToken)
		if dt.Add(token) {
			status = fiber.StatusOK
			message = "Device Token Added"
		} else {
			status = fiber.StatusInternalServerError
			message = "Unable to Add Device Token"
		}
	}
	return utils.GetMessageResponse(message, message, status, flag, ctx)
}

func delete_device_token(ctx *fiber.Ctx) error {
	var dt models.DeviceToken
	var status int = fiber.StatusOK
	var flag = dt.Delete(ctx.Params("device_token"))
	if !flag {
		status = fiber.StatusInternalServerError
	}
	return utils.GetMessageResponse("Device Token Deleted", "Unable to delete device token", status, flag, ctx)
}
