package handlers

import (
	"leichil_tv_website/views/pages"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

func Home(c *fiber.Ctx) error {
	handler := templ.Handler(pages.Home())
	return adaptor.HTTPHandler(handler)(c)
}
