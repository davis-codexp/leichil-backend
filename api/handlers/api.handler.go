package handlers

import (
	"leichil-backend/api/models"
	"leichil-backend/api/utils"

	"github.com/gofiber/fiber/v2"
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

	content_handler(admin_content_routes)
	user_handler(admin_user_routes)
	actor_handler(admin_actor_routes)
	producer_handler(admin_producer_routes)
	customer_handler(admin_customer_routes)
	subscription_plan_handler(admin_subscription_plan_routes)
	subscription_handler(admin_subscription_routes)
}

func AppApiHandler(router fiber.Router) {
	var app_subscription_routes = router.Group("/subscriptions")
	var app_subscription_plan_routes = router.Group("/subscription-plans")
	var app_activity_routes = router.Group("/activities")
	var app_content_routes = router.Group("/contents")
	var app_profile_routes = router.Group("/profile")
	var app_actor_routes = router.Group("/actors")
	var app_order_routes = router.Group("/orders")

	app_subscription_handler(app_subscription_routes)
	app_subscription_plan_handler(app_subscription_plan_routes)
	app_cca_handler(app_activity_routes)
	app_content_handler(app_content_routes)
	app_profile_handler(app_profile_routes)
	app_actor_handler(app_actor_routes)
	app_order_handler(app_order_routes)
}
