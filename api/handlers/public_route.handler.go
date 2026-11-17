package handlers

import (
	"fmt"
	"leichil-backend/api/models"
	"leichil-backend/api/templates"
	"leichil-backend/api/utils"
	"os"
	"strings"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
)

func public_route_handler(router fiber.Router) {
	router.Get("/browse-contents", public_content_handler)
	router.Get("/browse-contents/tabs", public_content_tab_handler)
	router.Get("/purchase", purchase_handler)
	router.Post("/purchase/verify", purchase_verify_handler)
	router.Post("/subscription/verify", subscription_verify_handler)
	router.Get("/payment-error", public_payment_error_handler)
	router.Get("/payment-success", public_payment_success_handler)
	router.Post("/orders/verify", verify_payment)

	var advance_booking_routes = router.Group("/advance-bookings")
	advance_booking_handler(advance_booking_routes)
}

func public_content_handler(ctx *fiber.Ctx) error {
	var content models.Content
	var params = map[string]string{
		"category":  "PREMIUM",
		"published": "1",
	}
	_, result, _ := content.Get(0, 10000, params)
	var contents = utils.FormatListResult[models.Content](result)
	var is_partial bool = ctx.Query("partial") == "true"
	return render_component(templates.BrowseContents("premium", contents), ctx, is_partial, "Buy & Watch")
}

func public_content_tab_handler(ctx *fiber.Ctx) error {
	var tab = strings.ToLower(ctx.Query("tab"))
	var comp templ.Component
	if tab == "premium" {
		var content models.Content
		var params = map[string]string{
			"category": "PREMIUM",
		}
		_, result, _ := content.Get(0, 10000, params)
		var contents = utils.FormatListResult[models.Content](result)
		comp = templates.PremiumContents(contents)
	} else {
		var plan models.SubscriptionPlan
		_, result, _ := plan.Get(0, 10000, map[string]string{})
		var plan_list = utils.FormatListResult[models.SubscriptionPlan](result)
		comp = templates.SubscriptionSection(plan_list)
	}
	var is_partial bool = ctx.Query("partial") == "true"
	return render_component(comp, ctx, is_partial, "Buy & Watch")
}

func purchase_handler(ctx *fiber.Ctx) error {
	var enc_content_id = ctx.Query("content")
	if len(enc_content_id) > 0 {
		content_id, _ := utils.Decrypt(enc_content_id)
		var content models.Content
		content, _ = content.GetById(content_id, 1)
		var is_partial bool = ctx.Query("partial") == "true"
		return render_component(templates.PremiumPurchase(content), ctx, is_partial, "Purchase Premium Content")
	}
	var enc_plan_id = ctx.Query("subscription")
	if len(enc_plan_id) > 0 {
		plan_id, _ := utils.Decrypt(enc_plan_id)
		var plan models.SubscriptionPlan
		result, _ := plan.GetById(plan_id)
		plan = utils.FormatSingleResult[models.SubscriptionPlan](result)
		var is_partial bool = ctx.Query("partial") == "true"
		return render_component(templates.SubscriptionPurchase(plan), ctx, is_partial, "Subscribe to Leichil TV")
	}
	return ctx.Status(fiber.StatusBadRequest).SendString("Invalid Request")
}

func purchase_verify_handler(ctx *fiber.Ctx) error {
	var body = new(models.WebPurchaseBody)
	flag, _, _ := utils.ParseAndValidate(body, []string{"Mobile", "Price", "ContentID", "Description"}, ctx)
	if flag {
		var obj models.Customer
		var filter = map[string]string{"mobile": models.SanitizeData(body.Mobile)}
		_, result, _ := obj.Get(0, 1, filter)
		if len(result) > 0 {
			obj = utils.FormatSingleResult[models.Customer](result)
			var token = models.CustomerToken{CustomerID: obj.ID}
			var po models.PaymentOrder
			var is_purchased bool = po.VerifyCustomerPurchase(token, fmt.Sprintf("%d", models.SanitizeData(body.ContentID)))
			if is_purchased {
				return render_component(templates.AlreadyPurchased(), ctx, true, "")
			}
			var api_key = os.Getenv("RZP_KEY_ID")
			var order = new(models.PaymentOrder)
			order.ContentID = body.ContentID
			var order_result = order.CreatePurchaseOrder(models.SanitizeData(body.Price)*100, obj)
			order_id, _ := order_result["order_id"].(string)
			amount, _ := order_result["amount"].(int)
			return render_component(
				templates.CustomerDetail(obj, api_key, order_id, models.SanitizeData(body.Description), amount),
				ctx, true, "",
			)
		} else {
			return render_component(templates.CustomerNotFound(), ctx, true, "")
		}
	}
	return ctx.Status(fiber.StatusBadRequest).SendString("Invalid Data")
}

func subscription_verify_handler(ctx *fiber.Ctx) error {
	var body = new(models.WebPurchaseBody)
	flag, _, _ := utils.ParseAndValidate(body, []string{"Mobile", "Price", "PlanID", "Description"}, ctx)
	if flag {
		var obj models.Customer
		var filter = map[string]string{"mobile": models.SanitizeData(body.Mobile)}
		_, result, _ := obj.Get(0, 1, filter)
		if len(result) > 0 {
			obj = utils.FormatSingleResult[models.Customer](result)
			var api_key = os.Getenv("RZP_KEY_ID")
			var order = new(models.PaymentOrder)
			order.PlanID = body.PlanID
			var order_result = order.CreatePurchaseOrder(models.SanitizeData(body.Price)*100, obj)
			order_id, _ := order_result["order_id"].(string)
			amount, _ := order_result["amount"].(int)
			return render_component(
				templates.CustomerDetail(obj, api_key, order_id, models.SanitizeData(body.Description), amount),
				ctx, true, "",
			)
		} else {
			return render_component(templates.CustomerNotFound(), ctx, true, "")
		}
	}
	return ctx.Status(fiber.StatusBadRequest).SendString("Invalid Data")
}

func purchase_subscription_verify_handler(ctx *fiber.Ctx) error {
	var body = new(models.WebPurchaseBody)
	flag, _, _ := utils.ParseAndValidate(body, []string{"Mobile", "Price", "PlanID", "Description"}, ctx)
	if flag {
		var obj models.Customer
		var filter = map[string]string{"mobile": models.SanitizeData(body.Mobile)}
		_, result, _ := obj.Get(0, 1, filter)
		if len(result) > 0 {
			obj = utils.FormatSingleResult[models.Customer](result)
			var order = new(models.PaymentOrder)
			var api_key = os.Getenv("RZP_KEY_ID")
			var order_result = order.CreatePurchaseOrder(models.SanitizeData(body.Price)*100, obj)
			order_id, _ := order_result["order_id"].(string)
			amount, _ := order_result["amount"].(int)
			return render_component(
				templates.CustomerDetail(obj, api_key, order_id, models.SanitizeData(body.Description), amount),
				ctx, true, "",
			)
		} else {
			return render_component(templates.CustomerNotFound(), ctx, true, "")
		}
	}
	return ctx.Status(fiber.StatusBadRequest).SendString("Invalid Data")
}

func public_payment_error_handler(ctx *fiber.Ctx) error {
	var is_partial bool = ctx.Query("partial") == "true"
	return render_component(templates.PaymentFailure(), ctx, is_partial, "Buy & Watch")
}

func public_payment_success_handler(ctx *fiber.Ctx) error {
	var is_partial bool = ctx.Query("partial") == "true"
	return render_component(templates.PaymentSuccess(), ctx, is_partial, "Buy & Watch")
}

func render_component(content templ.Component, ctx *fiber.Ctx, partial bool, title string) error {
	var comp templ.Component
	if partial {
		comp = content
	} else {
		comp = templates.Base(title, content)
	}
	return Render(ctx, comp)
}
