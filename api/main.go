package main

import (
	"leichil-backend/api/handlers"
	"leichil-backend/api/utils"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err.Error())
		panic(err.Error())
	}
	var PORT string = os.Getenv("PORT")

	conn, err := utils.GetConnection()
	if err != nil {
		log.Fatal("Unable to establish DB Connection", err.Error())
		panic(err.Error())
	}
	defer conn.Close()

	rdb, err := utils.GetRedisConnection()
	if err != nil {
		log.Fatal("Unable to establish Redis Connection", err.Error())
		panic(err.Error())
	}
	defer rdb.Close()

	var app *fiber.App = fiber.New()
	app.Static("/static", "./static")
	app.Use(cors.New(cors.Config{
		AllowMethods: "GET, POST, PUT, DELETE",
	}))
	var api_router fiber.Router = app.Group("/api/v1")
	var ui_router fiber.Router = app.Group("/")
	handlers.APIHandler(api_router)
	handlers.UIHandler(ui_router)
	if err = app.Listen(":" + PORT); err != nil {
		panic(err.Error())
	}
}
