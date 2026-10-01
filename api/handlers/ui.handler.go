package handlers

import (
	"leichil-backend/api/templates"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

func UIHandler(router fiber.Router) {
	var admin_ui_routes = router.Group("/admin")

	var producer_ui = router.Group("/producers")
	var actor_ui = router.Group("/actors")

	admin_ui_handler(admin_ui_routes)

	producer_handler(producer_ui)
	actor_handler(actor_ui)
}

func admin_ui_handler(router fiber.Router) {
	router.Get("/", func(ctx *fiber.Ctx) error {
		return Render(ctx, templates.Login(""))
	})
}

func Render(ctx *fiber.Ctx, component templ.Component) error {
	var handler = adaptor.HTTPHandler(templ.Handler(component))
	return handler(ctx)
}
