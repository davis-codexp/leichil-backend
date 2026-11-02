package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"leichil-backend/api/models"
	"leichil-backend/api/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func authenticate_user(ctx *fiber.Ctx) error {
	var token []string = ctx.GetReqHeaders()["Authorization"]
	if len(token) == 0 {
		return utils.GetMessageResponse("", "Access Denied", fiber.StatusUnauthorized, false, ctx)
	}
	credentials, flag := models.AuthenticateUser(token[0])
	if flag {
		ctx.Locals("user_token", credentials)
		return ctx.Next()
	}
	return utils.GetMessageResponse("", "Access Denied", fiber.StatusUnauthorized, false, ctx)
}

func authenticate_customer(ctx *fiber.Ctx) error {
	var token []string = ctx.GetReqHeaders()["Authorization"]
	if len(token) == 0 {
		return utils.GetMessageResponse("", "Access Denied", fiber.StatusUnauthorized, false, ctx)
	}
	credentials, flag := models.AuthenticateCustomer(token[0])
	if flag {
		ctx.Locals("customer_token", credentials)
		return ctx.Next()
	}
	return utils.GetMessageResponse("", "Access Denied", fiber.StatusUnauthorized, false, ctx)
}

func WebhookHandler(ctx *fiber.Ctx, webhook_secret string) error {
	var raw_body []byte = ctx.Body()
	if len(raw_body) == 0 {
		return ctx.Status(fiber.StatusBadRequest).SendString("Empty Body")
	}

	var signature string = ctx.Get("X-Razorpay-Signature")
	if signature == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("Missing Signature")
	}

	if !verify_webhook_signature(raw_body, signature, webhook_secret) {
		log.Info("Invalid Razorpay Webhook Signature")
		return ctx.Status(fiber.StatusBadRequest).SendString("Invalid Signature")
	}

	var event models.RazorpayEvent
	if err := json.Unmarshal(raw_body, &event); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString("Invalid JSON")
	}

	//eventID := ctx.Get("X-Razorpay-Event-Id")
	//if eventID != "" {
	// TODO: check if this eventID was already processed in your DB
	// if alreadyProcessed(eventID) {
	//     return c.SendStatus(fiber.StatusOK)
	// }
	//}

	switch event.Event {
	case "payment.captured":
		handle_payment_captured(event)
	case "payment.failed":
		handle_payment_failed(event)
	case "order.paid":
		handle_payment_captured(event)
	default:
		log.Info("Unhandled event: %s", event.Event)
	}
	return ctx.SendStatus(fiber.StatusOK)
}

func verify_webhook_signature(body []byte, signature, secret string) bool {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(body)
	expected := hex.EncodeToString(h.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(signature))
}

func handle_payment_captured(event models.RazorpayEvent) {
	log.Info("Payment Captured Event received on Webhook")
	var payload struct {
		Payment struct {
			Entity models.PaymentEntity `json:"entity"`
		} `json:"payment"`
	}
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		log.Info("failed to parse payment.captured:", err)
		return
	}
	var pe models.PaymentEntity = payload.Payment.Entity
	var po models.PaymentOrder
	if po.UpdateOrder("SUCCESS", pe.ID, pe.OrderID) {
		log.Info("Payment Status Updated to SUCCESS for ORDER ID: ", pe.OrderID)
	} else {
		log.Info("Unable to Update Payment Status for ORDER ID: ", pe.OrderID)
	}
	po.UpdateSubscription(pe.OrderID)
}

func handle_payment_failed(event models.RazorpayEvent) {
	log.Info("Payment Failed Event received on Webhook")
	var payload struct {
		Payment struct {
			Entity models.PaymentEntity `json:"entity"`
		} `json:"payment"`
	}
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		log.Info("failed to parse payment.captured:", err)
		return
	}
	var pe models.PaymentEntity = payload.Payment.Entity
	var po models.PaymentOrder
	if po.UpdateOrder("FAILED", pe.ID, pe.OrderID) {
		log.Info("Payment Status Updated to FAILED for ORDER ID: ", pe.OrderID)
	} else {
		log.Info("Unable to Update Payment Status for ORDER ID: ", pe.OrderID)
	}
}

func APIHandler(router fiber.Router) {
	var admin_routes = router.Group("/admin", authenticate_user)
	var app_routes = router.Group("/app", authenticate_customer)
	var auth_routes = router.Group("/auth")

	AdminApiHandler(admin_routes)
	AuthHandler(auth_routes)
	AppApiHandler(app_routes)
}

func AdminApiHandler(router fiber.Router) {
	var admin_content_routes = router.Group("/contents")
	var admin_user_routes = router.Group("/users")
	var admin_actor_routes = router.Group("/actors")
	var admin_producer_routes = router.Group("/producers")
	var admin_customer_routes = router.Group("/customers")
	var admin_subscription_plan_routes = router.Group("/subscription-plans")
	var admin_subscription_routes = router.Group("/subscriptions")
	var admin_notification_routes = router.Group("/notifications")
	var admin_order_routes = router.Group("/orders")

	content_handler(admin_content_routes)
	user_handler(admin_user_routes)
	actor_handler(admin_actor_routes)
	producer_handler(admin_producer_routes)
	customer_handler(admin_customer_routes)
	subscription_plan_handler(admin_subscription_plan_routes)
	subscription_handler(admin_subscription_routes)
	notification_handler(admin_notification_routes)
	order_handler(admin_order_routes)
}

func AppApiHandler(router fiber.Router) {
	var app_subscription_routes = router.Group("/subscriptions")
	var app_subscription_plan_routes = router.Group("/subscription-plans")
	var app_activity_routes = router.Group("/activities")
	var app_content_routes = router.Group("/contents")
	var app_profile_routes = router.Group("/profile")
	var app_actor_routes = router.Group("/actors")
	var app_order_routes = router.Group("/orders")
	var app_stream_routes = router.Group("/stream")

	app_subscription_handler(app_subscription_routes)
	app_subscription_plan_handler(app_subscription_plan_routes)
	app_cca_handler(app_activity_routes)
	app_content_handler(app_content_routes)
	app_profile_handler(app_profile_routes)
	app_actor_handler(app_actor_routes)
	app_order_handler(app_order_routes)
	app_stream_handler(app_stream_routes)
}
