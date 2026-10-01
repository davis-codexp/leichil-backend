package handlers

import (
	"fmt"
	"leichil-backend/api/models"
	"leichil-backend/api/utils"

	"github.com/gofiber/fiber/v2"
)

var subscription_plan_required_fields = []string{"Name", "Features", "Price", "DurationMonths", "IsActive"}

func subscription_plan_handler(router fiber.Router) {
	router.Get("/", get_subscription_plans)
	router.Get("/:id", get_subscription_plan_by_id)
	router.Post("/", add_subscription_plan)
	router.Put("/:id", update_subscription_plan)
	router.Delete("/:id", delete_subscription_plan)
}

func app_subscription_plan_handler(router fiber.Router) {
	router.Get("/", get_subscription_plans)
}

func get_subscription_plans(ctx *fiber.Ctx) error {
	var subscription_plan models.SubscriptionPlan
	var offset int = utils.GetIntQueryParams(ctx.Query("offset"), 0)
	var limit int = utils.GetIntQueryParams(ctx.Query("limit"), 10)
	var params = map[string]string{
		"name":      ctx.Query("name"),
		"price":     ctx.Query("price"),
		"duration":  ctx.Query("duration"),
		"is_active": ctx.Query("is_active"),
	}
	var key = fmt.Sprintf("subscription-plans-%s", ctx.OriginalURL())
	cache_data, err := utils.GetData(key)
	if err != nil {
		_, result, err := subscription_plan.Get(offset, limit, params)
		if len(result) > 0 {
			utils.SetListData[models.SubscriptionPlan](key, result, 0)
		}
		return utils.GetArrayResponse(ctx, 0, result, "Unable to get Subscription Plans", err)
	}
	return utils.GetCachedArrayResponse(ctx, 0, cache_data)
}

func get_subscription_plan_by_id(ctx *fiber.Ctx) error {
	var subscription_plan models.SubscriptionPlan
	result, err := subscription_plan.GetById(ctx.Params("id"))
	return utils.GetSingleResponse(ctx, result, "Unable to get SubscriptionPlan", err)
}

func add_subscription_plan(ctx *fiber.Ctx) error {
	var subscription_plan = new(models.SubscriptionPlan)
	flag, status, message := utils.ParseAndValidate(subscription_plan, subscription_plan_required_fields, ctx)
	if flag {
		var token = ctx.Locals("user_token").(models.UserToken)
		if subscription_plan.Add(token) {
			status = fiber.StatusOK
			message = "Subscription Plan Added"
			utils.DeleteKeyPattern("subscription-plans*")
		} else {
			status = fiber.StatusInternalServerError
			message = "Unable to Add Subscription Plan"
		}
	}
	return utils.GetMessageResponse(message, message, status, flag, ctx)
}

func update_subscription_plan(ctx *fiber.Ctx) error {
	var subscription_plan = new(models.SubscriptionPlan)
	flag, status, message := utils.ParseAndValidate(subscription_plan, subscription_plan_required_fields, ctx)
	if flag {
		if subscription_plan.Update(ctx.Params("id")) {
			status = fiber.StatusOK
			message = "Subscription Plan Updated"
			utils.DeleteKeyPattern("subscription-plans*")
		} else {
			status = fiber.StatusInternalServerError
			message = "Unable to Update Subscription Plan"
		}
	}
	return utils.GetMessageResponse(message, message, status, flag, ctx)
}

func delete_subscription_plan(ctx *fiber.Ctx) error {
	var subscription_plan models.SubscriptionPlan
	var flag bool = subscription_plan.Delete(ctx.Params("id"))
	if flag {
		utils.DeleteKeyPattern("subscription-plans*")
	}
	return utils.GetMessageResponse("Subscription Plan Deleted", "Unable to Delete Subscription Plan", fiber.StatusOK, flag, ctx)
}

/*---------------------- SUBSCRIPTION HANDLERS ------------------*/
var subscription_required_fields = []string{"CustomerID", "PlanID"}

func subscription_handler(router fiber.Router) {
	router.Get("/", get_subscriptions)
	router.Get("/:id", get_subscription_by_id)
	router.Put("/:id", update_subscription)
	router.Delete("/:id", delete_subscription)
}

func app_subscription_handler(router fiber.Router) {
	router.Get("/", get_customer_subscriptions)
	router.Post("/", add_subscription)
	router.Get("/verify", verify_customer_subscription)
}

func get_subscriptions(ctx *fiber.Ctx) error {
	var subscription models.Subscription
	var offset int = utils.GetIntQueryParams(ctx.Query("offset"), 0)
	var limit int = utils.GetIntQueryParams(ctx.Query("limit"), 10)
	var params = map[string]string{
		"customer": ctx.Query("customer"),
		"plan":     ctx.Query("plan"),
		"date":     ctx.Query("date"),
		"expired":  ctx.Query("expired"),
	}
	count, result, err := subscription.Get(offset, limit, params)
	var count_data int = models.SanitizeData(count.Count)
	return utils.GetArrayResponse(ctx, count_data, result, "Unable to get Subscriptions", err)
}

func get_customer_subscriptions(ctx *fiber.Ctx) error {
	var subscription models.Subscription
	var offset int = utils.GetIntQueryParams(ctx.Query("offset"), 0)
	var limit int = utils.GetIntQueryParams(ctx.Query("limit"), 10)
	var token = ctx.Locals("customer_token").(models.CustomerToken)
	var params = map[string]string{
		"customer": fmt.Sprintf("%d", models.SanitizeData(token.CustomerID)),
		"plan":     ctx.Query("plan"),
		"date":     ctx.Query("date"),
		"expired":  ctx.Query("expired"),
	}
	count, result, err := subscription.Get(offset, limit, params)
	var count_data int = models.SanitizeData(count.Count)
	return utils.GetArrayResponse(ctx, count_data, result, "Unable to get Subscriptions", err)
}

func get_subscription_by_id(ctx *fiber.Ctx) error {
	var subscription models.Subscription
	result, err := subscription.GetById(ctx.Params("id"))
	return utils.GetSingleResponse(ctx, result, "Unable to get Subscription", err)
}

func verify_customer_subscription(ctx *fiber.Ctx) error {
	var subscription models.Subscription
	result, err := subscription.GetCustomerSubscription(ctx.Locals("customer_token").(models.CustomerToken))
	return utils.GetSingleResponse(ctx, result, "Unable to get Customer Subscription", err)
}

func add_subscription(ctx *fiber.Ctx) error {
	var subscription = new(models.Subscription)
	flag, status, message := utils.ParseAndValidate(subscription, subscription_required_fields, ctx)
	if flag {
		var token = ctx.Locals("customer_token").(models.CustomerToken)
		if subscription.Add(token) {
			status = fiber.StatusOK
			message = "Subscription Added"
		} else {
			status = fiber.StatusInternalServerError
			message = "Unable to Add Subscription"
		}
	}
	return utils.GetMessageResponse(message, message, status, flag, ctx)
}

func update_subscription(ctx *fiber.Ctx) error {
	var subscription = new(models.Subscription)
	flag, status, message := utils.ParseAndValidate(subscription, subscription_required_fields, ctx)
	if flag {
		if subscription.Update(ctx.Params("id")) {
			status = fiber.StatusOK
			message = "Subscription Updated"
		} else {
			status = fiber.StatusInternalServerError
			message = "Unable to Update Subscription"
		}
	}
	return utils.GetMessageResponse(message, message, status, flag, ctx)
}

func delete_subscription(ctx *fiber.Ctx) error {
	var subscription models.Subscription
	var flag bool = subscription.Delete(ctx.Params("id"))
	return utils.GetMessageResponse("Subscription Deleted", "Unable to Delete Subscription", fiber.StatusOK, flag, ctx)
}
