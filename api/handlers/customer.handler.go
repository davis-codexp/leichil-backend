package handlers

import (
	"fmt"
	"leichil-backend/api/models"
	"leichil-backend/api/utils"

	"github.com/gofiber/fiber/v2"
)

var customer_required_fields = []string{"Name", "Mobile"}

func customer_handler(router fiber.Router) {
	router.Get("/", get_customers)
	router.Get("/:id", get_customer_by_id)
	router.Post("/", add_customer)
	router.Put("/:id", update_customer)
	router.Delete("/:id", delete_customer)
}

func app_profile_handler(router fiber.Router) {
	router.Get("/", get_profile)
	router.Put("/", update_profile)
}

/*----------------------API ROUTES-----------------*/
func get_customers(ctx *fiber.Ctx) error {
	var customer models.Customer
	var offset int = utils.GetIntQueryParams(ctx.Query("offset"), 0)
	var limit int = utils.GetIntQueryParams(ctx.Query("limit"), 10)
	var params = map[string]string{
		"name":   ctx.Query("name"),
		"mobile": ctx.Query("mobile"),
	}
	count, result, err := customer.Get(offset, limit, params)
	var count_data int = models.SanitizeData(count.Count)
	return utils.GetArrayResponse(ctx, count_data, result, "Unable to get Customers", err)
}

func get_customer_by_id(ctx *fiber.Ctx) error {
	var customer models.Customer
	result, err := customer.GetById(ctx.Params("id"))
	return utils.GetSingleResponse(ctx, result, "Unable to get Customer", err)
}

func add_customer(ctx *fiber.Ctx) error {
	var customer = new(models.Customer)
	flag, status, message := utils.ParseAndValidate(customer, customer_required_fields, ctx)
	if flag {
		if customer.Add() {
			status = fiber.StatusOK
			message = "Customer Added"
		} else {
			status = fiber.StatusInternalServerError
			message = "Unable to Add Customer"
		}
	}
	return utils.GetMessageResponse(message, message, status, flag, ctx)
}

func update_customer(ctx *fiber.Ctx) error {
	var customer = new(models.Customer)
	flag, status, message := utils.ParseAndValidate(customer, customer_required_fields, ctx)
	if flag {
		if customer.Update(ctx.Params("id")) {
			status = fiber.StatusOK
			message = "Customer Updated"
		} else {
			status = fiber.StatusInternalServerError
			message = "Unable to Update Customer"
		}
	}
	return utils.GetMessageResponse(message, message, status, flag, ctx)
}

func delete_customer(ctx *fiber.Ctx) error {
	var customer models.Customer
	var flag bool = customer.Delete(ctx.Params("id"))
	return utils.GetMessageResponse("Customer Deleted", "Unable to Delete Customer", fiber.StatusOK, flag, ctx)
}

/*-------------------PROFILE ROUTES---------------------*/
func get_profile(ctx *fiber.Ctx) error {
	var profile models.CustomerProfile
	var token = ctx.Locals("customer_token").(models.CustomerToken)
	result, err := profile.Get(token)
	if err != nil {
		utils.GetMessageResponse("", "Unable to get Profile", fiber.StatusInternalServerError, false, ctx)
	}
	var response = map[string]any{
		"success": true,
		"result":  result,
	}
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func update_profile(ctx *fiber.Ctx) error {
	var customer = new(models.Customer)
	flag, status, message := utils.ParseAndValidate(customer, customer_required_fields, ctx)
	if flag {
		var token = ctx.Locals("customer_token").(models.CustomerToken)
		if customer.Update(fmt.Sprintf(`%d`, models.SanitizeData(token.CustomerID))) {
			status = fiber.StatusOK
			message = "Profile Updated"
		} else {
			status = fiber.StatusInternalServerError
			message = "Unable to Update Profile"
		}
	}
	return utils.GetMessageResponse(message, message, status, flag, ctx)
}
