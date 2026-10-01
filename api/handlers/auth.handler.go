package handlers

import (
	"leichil-backend/api/models"
	"leichil-backend/api/utils"

	"github.com/gofiber/fiber/v2"
)

func AuthHandler(router fiber.Router) {
	router.Post("/admin/login", user_login)
	router.Post("/signup", customer_sign_up)
	router.Post("/otp", get_customer_otp)
	router.Post("/otp/verify", verify_customer_otp)
	router.Post("/customer/logout", authenticate_customer, customer_logout)
}

func user_login(ctx *fiber.Ctx) error {
	var required_fields = []string{"Mobile", "Password"}
	var login = new(models.UserLoginBody)
	var result models.UserLoginResult
	var success = "Authentication Successful"
	flag, status, failure := utils.ParseAndValidate(login, required_fields, ctx)
	if flag {
		result, flag = login.Login()
		if flag {
			return utils.GetSingleResponse(ctx, []any{result}, "", nil)
		}
		failure = "Invalid Credentials"
	}
	return utils.GetMessageResponse(success, failure, status, flag, ctx)
}

func customer_sign_up(ctx *fiber.Ctx) error {
	var customer = new(models.Customer)
	flag, status, message := utils.ParseAndValidate(customer, customer_required_fields, ctx)
	if flag {
		if customer.IsDuplicateMobile() {
			return utils.GetMessageResponse(
				"", "This Mobile number is already registered",
				fiber.StatusBadRequest, false, ctx,
			)
		}
		if customer.Add() {
			return utils.GetMessageResponse("Successfully Signed Up", "", fiber.StatusOK, true, ctx)
		} else {
			return utils.GetMessageResponse("", "Unable to sign up", fiber.StatusInternalServerError, false, ctx)
		}
	}
	return utils.GetMessageResponse(message, message, status, flag, ctx)
}

func get_customer_otp(ctx *fiber.Ctx) error {
	var body = new(models.CustomerLoginBody)
	var required_fields = []string{"Mobile"}
	flag, status, message := utils.ParseAndValidate(body, required_fields, ctx)
	if flag {
		if body.GetOtp() {
			flag = true
			message = "An OTP has been sent to your registered mobile number"
		} else {
			flag = false
			message = "Invalid Mobile Number"
		}
	}
	return utils.GetMessageResponse(message, message, status, flag, ctx)
}

func verify_customer_otp(ctx *fiber.Ctx) error {
	var body = new(models.CustomerLoginBody)
	var required_fields = []string{"Mobile", "Otp"}
	flag, status, message := utils.ParseAndValidate(body, required_fields, ctx)
	if flag {
		result, otp_flag := body.VerifyOtp()
		if otp_flag {
			return utils.GetSingleResponse(ctx, []any{result}, "Unable to verify customer OTP", nil)
		} else {
			flag = false
			message = "Invalid OTP"
		}
	}
	return utils.GetMessageResponse(message, message, status, flag, ctx)
}

func customer_logout(ctx *fiber.Ctx) error {
	var customer models.Customer
	var token = ctx.Locals("customer_token").(models.CustomerToken)
	var status = fiber.StatusInternalServerError
	var flag = false
	if customer.Logout(token) {
		status = fiber.StatusOK
		flag = true
	}
	return utils.GetMessageResponse("Logged Out", "Unable to Log Out", status, flag, ctx)
}
