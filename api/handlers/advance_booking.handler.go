package handlers

import (
	"leichil-backend/api/models"
	"leichil-backend/api/templates"
	"leichil-backend/api/utils"

	"github.com/gofiber/fiber/v2"
)

var advance_booking_required_fields = []string{"Name", "Mobile"}

func advance_booking_handler(router fiber.Router) {
	router.Get("/", get_advance_bookings)
	router.Get("/:content_id", booking_handler)
}

/*----------------------API ROUTES-----------------*/
func get_advance_bookings(ctx *fiber.Ctx) error {
	var advance_booking models.AdvanceBooking
	_, result, _ := advance_booking.Get(0, 1000, map[string]string{})
	var list = utils.FormatListResult[models.AdvanceBooking](result)
	var is_partial bool = ctx.Query("partial") == "true"
	return render_component(templates.AdvanceBookings(list), ctx, is_partial, "Advance Bookings")
}

func booking_handler(ctx *fiber.Ctx) error {
	var enc_content_id = ctx.Params("content_id")
	if len(enc_content_id) > 0 {
		content_id, _ := utils.Decrypt(enc_content_id)
		var obj models.AdvanceBooking
		result, err := obj.GetByContentId(content_id)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).SendString("An Error Occured")
		}
		obj = utils.FormatSingleResult[models.AdvanceBooking](result)
		var is_partial bool = ctx.Query("partial") == "true"
		return render_component(templates.BookContent(obj), ctx, is_partial, "Book Content")
	}
	return ctx.Status(fiber.StatusBadRequest).SendString("Invalid Request")
}
