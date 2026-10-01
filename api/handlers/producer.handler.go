package handlers

import (
	"leichil-backend/api/models"
	"leichil-backend/api/utils"

	"github.com/gofiber/fiber/v2"
)

var producer_required_fields = []string{"Name", "Mobile", "Address"}

func producer_handler(router fiber.Router) {
	router.Get("/", get_producers)
	router.Get("/:id", get_producer_by_id)
	router.Post("/", add_producer)
	router.Put("/:id", update_producer)
	router.Delete("/:id", delete_producer)
}

func get_producers(ctx *fiber.Ctx) error {
	var producer models.Producer
	var offset int = utils.GetIntQueryParams(ctx.Query("offset"), 0)
	var limit int = utils.GetIntQueryParams(ctx.Query("limit"), 10)
	var params = map[string]string{
		"name":   ctx.Query("name"),
		"mobile": ctx.Query("mobile"),
		"email":  ctx.Query("email"),
	}
	count, result, err := producer.Get(offset, limit, params)
	var count_data int = models.SanitizeData(count.Count)
	return utils.GetArrayResponse(ctx, count_data, result, "Unable to get Producers", err)
}

func get_producer_by_id(ctx *fiber.Ctx) error {
	var producer models.Producer
	result, err := producer.GetById(ctx.Params("id"))
	return utils.GetSingleResponse(ctx, result, "Unable to get Producer", err)
}

func add_producer(ctx *fiber.Ctx) error {
	var producer = new(models.Producer)
	flag, status, message := utils.ParseAndValidate(producer, producer_required_fields, ctx)
	if flag {
		var token = ctx.Locals("user_token").(models.UserToken)
		if producer.Add(token) {
			status = fiber.StatusOK
			message = "Producer Added"
		} else {
			status = fiber.StatusInternalServerError
			message = "Unable to Add Producer"
		}
	}
	return utils.GetMessageResponse(message, message, status, flag, ctx)
}

func update_producer(ctx *fiber.Ctx) error {
	var producer = new(models.Producer)
	flag, status, message := utils.ParseAndValidate(producer, producer_required_fields, ctx)
	if flag {
		if producer.Update(ctx.Params("id")) {
			status = fiber.StatusOK
			message = "Producer Updated"
		} else {
			status = fiber.StatusInternalServerError
			message = "Unable to Update Producer"
		}
	}
	return utils.GetMessageResponse(message, message, status, flag, ctx)
}

func delete_producer(ctx *fiber.Ctx) error {
	var producer models.Producer
	var flag bool = producer.Delete(ctx.Params("id"))
	return utils.GetMessageResponse("Producer Deleted", "Unable to Delete Producer", fiber.StatusOK, flag, ctx)
}
