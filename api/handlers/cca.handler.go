package handlers

import (
	"leichil-backend/api/models"
	"leichil-backend/api/utils"

	"github.com/gofiber/fiber/v2"
)

var cca_required_fields = []string{"ContentID"}
var favourite_required_fields = []string{"ContentID"}

func app_cca_handler(router fiber.Router) {
	router.Post("/", add_content_activity)
	router.Get("/", get_content_activities)
	router.Post("/favourites", add_customer_favourite)
	router.Get("/favourites", get_customer_favourites)
	router.Delete("/favourites/:content_id", delete_customer_favourite)
}

func get_content_activities(ctx *fiber.Ctx) error {
	var cca models.CustomerContentActivity
	var offset int = utils.GetIntQueryParams(ctx.Query("offset"), 0)
	var limit int = utils.GetIntQueryParams(ctx.Query("limit"), 10)
	var params = map[string]string{
		"content":    ctx.Query("content"),
		"customer":   ctx.Query("customer"),
		"episode":    ctx.Query("episode"),
		"start_date": ctx.Query("start_date"),
		"end_date":   ctx.Query("end_date"),
	}
	_, result, err := cca.Get(offset, limit, params)
	return utils.GetArrayResponse(ctx, 0, result, "Unable to get Content Activities", err)
}

func add_content_activity(ctx *fiber.Ctx) error {
	var cca = new(models.CustomerContentActivity)
	flag, status, message := utils.ParseAndValidate(cca, cca_required_fields, ctx)
	if flag {
		if cca.Add(ctx.Locals("customer_token").(models.CustomerToken)) {
			status = fiber.StatusOK
			message = "Activity Added"
		} else {
			status = fiber.StatusInternalServerError
			message = "Unable to Add Activity"
		}
	}
	return utils.GetMessageResponse(message, message, status, flag, ctx)
}

/*----------------------------CUSTOMER FAVOURITES----------------------------*/
func get_customer_favourites(ctx *fiber.Ctx) error {
	var fav models.CustomerFavourite
	var offset int = utils.GetIntQueryParams(ctx.Query("offset"), 0)
	var limit int = utils.GetIntQueryParams(ctx.Query("limit"), 10)
	var token = ctx.Locals("customer_token").(models.CustomerToken)
	_, result, err := fav.Get(offset, limit, token)
	return utils.GetArrayResponse(ctx, 0, result, "Unable to get Favourites", err)
}

func add_customer_favourite(ctx *fiber.Ctx) error {
	var fav = new(models.CustomerFavourite)
	flag, status, message := utils.ParseAndValidate(fav, favourite_required_fields, ctx)
	if flag {
		if fav.Add(ctx.Locals("customer_token").(models.CustomerToken)) {
			status = fiber.StatusOK
			message = "Added to Favourites"
		} else {
			status = fiber.StatusInternalServerError
			message = "Unable to Add to Favourites"
		}
	}
	return utils.GetMessageResponse(message, message, status, flag, ctx)
}

func delete_customer_favourite(ctx *fiber.Ctx) error {
	var fav models.CustomerFavourite
	var token = ctx.Locals("customer_token").(models.CustomerToken)
	if fav.Delete(ctx.Params("content_id"), token) {
		return utils.GetMessageResponse("Deleted from Favourites", "", fiber.StatusOK, true, ctx)
	}
	return utils.GetMessageResponse("", "Unable to delete from Favourites", fiber.StatusInternalServerError, false, ctx)
}
