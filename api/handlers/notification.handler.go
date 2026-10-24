package handlers

import (
	"leichil-backend/api/models"
	"leichil-backend/api/utils"

	"github.com/gofiber/fiber/v2"
)

func notification_handler(router fiber.Router) {
	router.Post("/", send_single_notification)
	router.Post("/silent", send_silent_notification)
}

func send_single_notification(ctx *fiber.Ctx) error {
	var body = new(utils.PushNotificationRequest)
	var required_fields = []string{"DeviceID", "Title", "Body"}
	flag, status, message := utils.ParseAndValidate(body, required_fields, ctx)
	if flag {
		var (
			device_id = models.SanitizeData(body.DeviceID)
			title     = models.SanitizeData(body.Title)
			body      = models.SanitizeData(body.Body)
		)
		utils.SendPushNotification(device_id, title, body, map[string]string{})
		message = "Push Notification Sent"
		status = fiber.StatusOK
	}
	return utils.GetMessageResponse(message, message, status, flag, ctx)
}

func send_silent_notification(ctx *fiber.Ctx) error {
	var body = new(utils.PushNotificationRequest)
	var required_fields = []string{"DeviceID", "Data"}
	flag, status, message := utils.ParseAndValidate(body, required_fields, ctx)
	if flag {
		var (
			device_id = models.SanitizeData(body.DeviceID)
			data      = *body.Data
		)
		utils.SendPushNotification(device_id, "", "", data)
		message = "Push Notification Sent"
		status = fiber.StatusOK
	}
	return utils.GetMessageResponse(message, message, status, flag, ctx)
}
