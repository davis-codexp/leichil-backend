package utils

import (
	"reflect"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func ParseAndValidate[T any](body T, fields []string, ctx *fiber.Ctx) (bool, int, string) {
	if err := ctx.BodyParser(body); err != nil {
		log.Error("Error in Parsing Body", err.Error())
		return false, fiber.StatusBadRequest, "Invalid POST Body"
	}
	if ValidateData(body, fields) {
		return true, fiber.StatusOK, ""
	}
	return false, fiber.StatusBadRequest, "Missing Required Fields"
}

func ValidateData[T any](body T, fields []string) bool {
	for _, field := range fields {
		r := reflect.ValueOf(body)
		value := reflect.Indirect(r).FieldByName(field)
		if value.IsNil() {
			log.Info("Missing: ", field)
			return false
		}
		dataType := value.Type().String()
		if dataType == "*string" {
			str := (*string)(value.UnsafePointer())
			if str == nil || len(*str) == 0 {
				return false
			}
		}
	}
	return true
}
