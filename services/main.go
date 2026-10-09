package main

import (
	"crypto/subtle"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/joho/godotenv"

	"leichil-backend/services/consumers"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err.Error())
	}
	var token = os.Getenv("TOKEN")
	if token == "" {
		log.Warn("WARNING: AUTH_TOKEN NOT LOADED")
	}
	initListeners()
	var app *fiber.App = fiber.New(fiber.Config{
		BodyLimit: 50 * 1024 * 1024,
	})
	app.Use(cors.New(cors.Config{
		AllowMethods: []string{"GET", "POST"},
	}))
	app.Use("/", static.New("./assets"))
	app.Get("/", func(ctx fiber.Ctx) error {
		return ctx.SendString("LEICHIL MEDIA SERVER")
	})

	app.Post("/api/v1/upload", func(ctx fiber.Ctx) error {
		var auth_token = ctx.Get("Authorization")
		if subtle.ConstantTimeCompare([]byte(auth_token), []byte(token)) != 1 {
			return ctx.Status(fiber.StatusUnauthorized).SendString("Access Denied")
		}
		var response = map[string]any{
			"success": false,
			"result":  nil,
		}
		file, err := ctx.FormFile("file")
		if err != nil {
			log.Error(err.Error())
			return ctx.Status(fiber.StatusBadRequest).JSON(response)
		}
		file_content, err := file.Open()
		if err != nil {
			log.Error("Unable to open uploaded file: ", err.Error())
			return ctx.Status(fiber.StatusInternalServerError).JSON(response)
		}
		defer file_content.Close()
		var buffer = make([]byte, 1024)
		_, err = file_content.Read(buffer)
		if err != nil && err.Error() != "EOF" {
			log.Error("Error in reading uploaded file: ", err.Error())
			return ctx.Status(fiber.StatusInternalServerError).JSON(response)
		}
		var mime_type string = consumers.GetMimeType(buffer, file.Filename)
		var upload_dir = "./assets"
		if err := os.MkdirAll(upload_dir, os.ModePerm); err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(response)
		}
		var current int64 = time.Now().Unix()
		var temp = rand.Int()
		var file_id = int(current) + temp
		var extension string = filepath.Ext(file.Filename)
		var file_name string = fmt.Sprintf(`%d%s`, file_id, extension)
		var file_path string = filepath.Join(upload_dir, file_name)
		if err := ctx.SaveFile(file, file_path); err != nil {
			log.Error("Error in saving original file: ", err.Error())
			return ctx.Status(fiber.StatusInternalServerError).JSON(response)
		}
		if strings.HasPrefix(mime_type, "image/") {
			go consumers.RunImageProcessorV2(file_name)
			file_name = fmt.Sprintf(`%d.%s`, file_id, "webp")
		}
		response["success"] = true
		response["result"] = file_name
		return ctx.Status(fiber.StatusOK).JSON(response)
	})

	/*
		app.Post("/api/v1/upload", func(ctx fiber.Ctx) error {
			var auth_token = ctx.Get("Authorization")
			if subtle.ConstantTimeCompare([]byte(auth_token), []byte(token)) != 1 {
				return ctx.Status(fiber.StatusUnauthorized).SendString("Access Denied")
			}
			var response = map[string]any{
				"success": false,
				"result":  nil,
			}
			file, err := ctx.FormFile("file")
			if err != nil {
				log.Error(err.Error())
				return ctx.Status(fiber.StatusBadRequest).JSON(response)
			}
			file_content, err := file.Open()
			if err != nil {
				log.Error("Unable to open uploaded file: ", err.Error())
				return ctx.Status(fiber.StatusInternalServerError).JSON(response)
			}
			defer file_content.Close()
			var buffer = make([]byte, 1024)
			_, err = file_content.Read(buffer)
			if err != nil && err.Error() != "EOF" {
				log.Error("Error in reading uploaded file: ", err.Error())
				return ctx.Status(fiber.StatusInternalServerError).JSON(response)
			}
			var mime_type string = http.DetectContentType(buffer)
			var upload_dir = "./assets"
			if err := os.MkdirAll(upload_dir, 0755); err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(response)
			}
			var current int64 = time.Now().Unix()
			var temp = rand.Int()
			var file_id = int(current) + temp
			var extension string = filepath.Ext(file.Filename)
			var file_name string = fmt.Sprintf("%d%s", file_id, extension)
			var file_path string = filepath.Join(upload_dir, file_name)
			if err := ctx.SaveFile(file, file_path); err != nil {
				log.Error("Error in saving uploaded file: ", err.Error())
				return ctx.Status(fiber.StatusInternalServerError).JSON(response)
			}
			if strings.HasPrefix(mime_type, "image/") {
				go consumers.RunImageProcessor(file_name)
				file_name = fmt.Sprintf("%d.%s", file_id, "webp")
			}
			response["success"] = true
			response["result"] = file_name
			return ctx.Status(fiber.StatusOK).JSON(response)
		})
	*/

	log.Fatal(app.Listen(":4500").Error())
}

func initListeners() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error("SMS Listener Panicked: %v. Attempting process reset. ", r)
			}
		}()
		consumers.SmsListener()
	}()
}
