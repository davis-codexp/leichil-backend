package handlers

import (
	"leichil-backend/api/models"
	"leichil-backend/api/utils"

	"github.com/gofiber/fiber/v2"
)

var user_required_fields = []string{"Name", "Mobile", "Email", "Role"}

func user_handler(router fiber.Router) {
	router.Post("/", add_user)
	router.Get("/", get_users)
	router.Get("/:id", get_user_by_id)
	router.Put("/:id", update_user)
	router.Delete("/:id", delete_user)
}

/*----------------------API ROUTES-----------------*/
func get_users(ctx *fiber.Ctx) error {
	var user models.User
	var offset int = utils.GetIntQueryParams(ctx.Query("offset"), 0)
	var limit int = utils.GetIntQueryParams(ctx.Query("limit"), 10)
	var params = map[string]string{
		"name":   ctx.Query("name"),
		"mobile": ctx.Query("mobile"),
		"role":   ctx.Query("role"),
		"email":  ctx.Query("email"),
	}
	count, result, err := user.Get(offset, limit, params)
	var count_data int = models.SanitizeData(count.Count)
	return utils.GetArrayResponse(ctx, count_data, result, "Unable to get Users", err)
}

func get_user_by_id(ctx *fiber.Ctx) error {
	var user models.User
	result, err := user.GetById(ctx.Params("id"))
	return utils.GetSingleResponse(ctx, result, "Unable to get User", err)
}

func add_user(ctx *fiber.Ctx) error {
	var user = new(models.User)
	flag, status, message := utils.ParseAndValidate(user, user_required_fields, ctx)
	if flag {
		if user.Add() {
			status = fiber.StatusOK
			message = "User Added"
		} else {
			status = fiber.StatusInternalServerError
			message = "Unable to Add User"
		}
	}
	return utils.GetMessageResponse(message, message, status, flag, ctx)
}

func update_user(ctx *fiber.Ctx) error {
	var user = new(models.User)
	flag, status, message := utils.ParseAndValidate(user, user_required_fields, ctx)
	if flag {
		if user.Update(ctx.Params("id")) {
			status = fiber.StatusOK
			message = "User Updated"
		} else {
			status = fiber.StatusInternalServerError
			message = "Unable to Update User"
		}
	}
	return utils.GetMessageResponse(message, message, status, flag, ctx)
}

func delete_user(ctx *fiber.Ctx) error {
	var user models.User
	var flag bool = user.Delete(ctx.Params("id"))
	return utils.GetMessageResponse("User Deleted", "Unable to Delete User", fiber.StatusOK, flag, ctx)
}
