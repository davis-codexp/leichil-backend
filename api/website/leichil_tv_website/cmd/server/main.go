package main

import (
	"log"

	"leichil_tv_website/internal/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New(fiber.Config{
		AppName:               "Leichil TV",
		DisableStartupMessage: true,
	})

	app.Static("/static", "./static")

	routes.Register(app)

	log.Println("Server running on http://localhost:3000")

	log.Fatal(app.Listen(":3000"))
}
