package utils

import (
	"encoding/json"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetMessageResponse(success string, failure string, status int, flag bool, ctx *fiber.Ctx) error {
	message := failure
	if flag {
		message = success
	}
	response := FormatMessageResponse(message, flag)
	return ctx.Status(status).JSON(response)
}

func FormatMessageResponse(message string, flag bool) map[string]any {
	return map[string]any{
		"success": flag,
		"message": message,
	}
}

func GetArrayResponse(ctx *fiber.Ctx, count int, result []any, message string, err error) error {
	if len(result) == 0 {
		result = []any{}
	}
	if err != nil {
		response := FormatMessageResponse(message, false)
		return ctx.Status(fiber.StatusInternalServerError).JSON(response)
	}
	resultSet := struct {
		Count int   `json:"count"`
		List  []any `json:"list"`
	}{count, result}
	response := map[string]any{
		"success": true,
		"result":  resultSet,
		"source":  "database",
	}
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func GetSingleResponse(ctx *fiber.Ctx, result []any, message string, err error) error {
	if err != nil {
		response := FormatMessageResponse(message, false)
		return ctx.Status(fiber.StatusInternalServerError).JSON(response)
	}
	response := map[string]any{
		"success": false,
		"result":  nil,
		"source":  "database",
	}
	if len(result) > 0 {
		response["success"] = true
		response["result"] = result[0]
	}
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func GetCachedArrayResponse(ctx *fiber.Ctx, count int, result string) error {
	var response = map[string]any{
		"success": true,
		"source":  "cache",
		"result": map[string]any{
			"count": count,
			"list":  json.RawMessage(result),
		},
	}
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func GetCachedSingleResponse(ctx *fiber.Ctx, result string) error {
	var response = map[string]any{
		"success": true,
		"source":  "cache",
		"result":  json.RawMessage(result),
	}
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func GetIntQueryParams(param string, default_value int) int {
	value, err := strconv.ParseInt(param, 10, 0)
	if err != nil {
		return default_value
	}
	return int(value)
}

func FormatListResult[T any](result []any) []T {
	var list []T
	for _, item := range result {
		list = append(list, item.(T))
	}
	return list
}

func FormatSingleResult[T any](result []any) T {
	var data T
	if len(result) > 0 {
		data = result[0].(T)
	}
	return data
}
