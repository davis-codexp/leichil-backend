package handlers

import (
	"fmt"
	"leichil-backend/api/models"
	"leichil-backend/api/utils"

	"github.com/gofiber/fiber/v2"
)

var actor_required_fields = []string{"Name", "Image"}

func actor_handler(router fiber.Router) {
	router.Get("/", get_actors)
	router.Get("/:id", get_actor_by_id)
	router.Post("/", add_actor)
	router.Put("/:id", update_actor)
	router.Delete("/:id", delete_actor)
}

func app_actor_handler(router fiber.Router) {
	router.Get("/", get_actors)
	router.Get("/:id", get_actor_by_id)
}

/*----------------------API ROUTES-----------------*/
func get_actors(ctx *fiber.Ctx) error {
	var actor models.Actor
	var offset int = utils.GetIntQueryParams(ctx.Query("offset"), 0)
	var limit int = utils.GetIntQueryParams(ctx.Query("limit"), 10)
	var params = map[string]string{
		"name":   ctx.Query("name"),
		"gender": ctx.Query("gender"),
	}

	var key = fmt.Sprintf("actor-%s", ctx.OriginalURL())
	cache_data, err := utils.GetData(key)
	if err != nil {
		_, result, err := actor.Get(offset, limit, params)
		if len(result) > 0 {
			utils.SetListData[models.Actor](key, result, 0)
		}
		return utils.GetArrayResponse(ctx, 0, result, "Unable to get Actors", err)
	}
	return utils.GetCachedArrayResponse(ctx, 0, cache_data)
}

func get_actor_by_id(ctx *fiber.Ctx) error {
	var id = ctx.Params("id")
	var key = fmt.Sprintf("actor-%s", id)
	cache_data, err := utils.GetData(key)
	if err != nil {
		var actor models.Actor
		result, err := actor.GetById(id)
		if err == nil && len(result) > 0 {
			utils.SetSingleData[models.Actor](key, result, 0)
		}
		return utils.GetSingleResponse(ctx, result, "Unable to get Actor", err)
	}
	return utils.GetCachedSingleResponse(ctx, cache_data)
}

func add_actor(ctx *fiber.Ctx) error {
	var actor = new(models.Actor)
	flag, status, message := utils.ParseAndValidate(actor, actor_required_fields, ctx)
	if flag {
		var token = ctx.Locals("user_token").(models.UserToken)
		if actor.Add(token) {
			status = fiber.StatusOK
			message = "Actor Added"
			utils.DeleteKeyPattern("actor*")
		} else {
			status = fiber.StatusInternalServerError
			message = "Unable to Add Actor"
		}
	}
	return utils.GetMessageResponse(message, message, status, flag, ctx)
}

func update_actor(ctx *fiber.Ctx) error {
	var actor = new(models.Actor)
	flag, status, message := utils.ParseAndValidate(actor, actor_required_fields, ctx)
	if flag {
		if actor.Update(ctx.Params("id")) {
			status = fiber.StatusOK
			message = "Actor Updated"
			utils.DeleteKeyPattern("actor*")
		} else {
			status = fiber.StatusInternalServerError
			message = "Unable to Update Actor"
		}
	}
	return utils.GetMessageResponse(message, message, status, flag, ctx)
}

func delete_actor(ctx *fiber.Ctx) error {
	var actor models.Actor
	var flag bool = actor.Delete(ctx.Params("id"))
	if flag {
		utils.DeleteKeyPattern("actor*")
	}
	return utils.GetMessageResponse("Actor Deleted", "Unable to Delete Actor", fiber.StatusOK, flag, ctx)
}
