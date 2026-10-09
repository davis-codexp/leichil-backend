package models

import (
	"encoding/json"
	"fmt"
	"leichil-backend/api/utils"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"golang.org/x/crypto/bcrypt"
)

type UserToken struct {
	Token     *string
	UserID    *int
	CreatedOn *string
	Name      *string
	Mobile    *string
	Email     *string
	Role      *string
}

type CustomerToken struct {
	Token      *string
	CustomerID *int
	CreatedOn  *string
	Name       *string
	Mobile     *string
}

type CustomerOtp struct {
	Otp        *int    `json:"otp"`
	CustomerID *int    `json:"customer_id"`
	Mobile     *string `json:"mobile"`
	Name       *string `json:"name"`
	CreatedOn  *string `json:"created_on"`
}

type UserLoginBody struct {
	Mobile   *string `json:"mobile"`
	Password *string `json:"password"`
}

type UserLoginResult struct {
	ID       *int    `json:"id,omitempty"`
	Name     *string `json:"fullName,omitempty"`
	Mobile   *string `json:"mobile,omitempty"`
	Role     *string `json:"designationName,omitempty"`
	Email    *string `json:"email,omitempty"`
	Password *string `json:"password,omitempty"`
	Token    *string `json:"token,omitempty"`
}

type CustomerLoginBody struct {
	Mobile *string `json:"mobile,omitempty"`
	Otp    *int    `json:"otp,omitempty"`
}

type CustomerLoginResult struct {
	ID        *int    `json:"id,omitempty"`
	Name      *string `json:"fullName,omitempty"`
	Mobile    *string `json:"mobile,omitempty"`
	Token     *string `json:"token,omitempty"`
	CreatedOn *string `json:"created_on,omitempty"`
}

/*-------------------------USER AUTHENTICATION---------------------------*/
func (body *UserLoginBody) Login() (UserLoginResult, bool) {
	var user UserLoginResult
	var fields = []any{
		&user.ID, &user.Name, &user.Mobile, &user.Email, &user.Role, &user.Password,
	}
	var filter = []any{body.Mobile}
	result, err := utils.RunQuery(user_login_query, fields, &user, filter)
	if err != nil {
		log.Error("Error in Getting User Data for Authentication ", err.Error())
		return user, false
	} else if len(result) == 0 {
		log.Info("No user found for this mobile ", SanitizeData(body.Mobile))
		return user, false
	}
	user = result[0].(UserLoginResult)
	var password = []byte(SanitizeData(body.Password))
	var encrypted = []byte(SanitizeData(user.Password))
	if err = bcrypt.CompareHashAndPassword(encrypted, password); err != nil {
		log.Info("Invalid Password")
		return user, false
	}
	var plain_text = []byte(SanitizeData(body.Password) + time.Now().String())
	token, err := utils.GenerateToken(plain_text)
	if err != nil {
		log.Error("Unable to Generate Token ", err.Error())
		return user, false
	}
	var token_fields = []any{token, user.ID}
	var insert_flag bool = utils.RunInsert(insert_user_token_query, token_fields)
	if !insert_flag {
		log.Error("Error in Insert Token ")
		return user, false
	}
	user.Token = &token
	user.Password = nil
	return user, true
}

func AuthenticateUser(auth_token string) (UserToken, bool) {
	var token UserToken
	var fields = []any{
		&token.Token, &token.UserID, &token.CreatedOn,
		&token.Name, &token.Mobile, &token.Email, &token.Role,
	}
	result, err := utils.RunQuery(get_user_token_query, fields, &token, []any{auth_token})
	if err != nil {
		log.Error("Error in Validate User Token Query ", err.Error())
		return token, false
	}
	if len(result) == 0 {
		return token, false
	}
	token = result[0].(UserToken)
	return token, true
}

/*--------------------CUSTOMER AUTHENTICATION---------------------*/
func (body *CustomerLoginBody) GetOtp() bool {
	var customer Customer
	var fields = []any{&customer.ID, &customer.Name, &customer.Mobile}
	var filter = []any{body.Mobile}
	result, err := utils.RunQuery(customer_login_query, fields, &customer, filter)
	if err != nil {
		log.Error("Error in Getting Customer Data for Authentication ", err.Error())
		return false
	} else if len(result) == 0 {
		log.Info("No customer found for this mobile ", SanitizeData(body.Mobile))
		return false
	}
	customer = result[0].(Customer)
	var otp_str = utils.GenerateOTP()
	otp, _ := strconv.Atoi(otp_str)
	var timestamp string = time.Now().Format("2006-01-02 15:04:05")
	var cus_otp = CustomerOtp{
		Otp:        &otp,
		CustomerID: customer.ID,
		Mobile:     customer.Mobile,
		Name:       customer.Name,
		CreatedOn:  &timestamp,
	}
	if utils.SetSingleData[CustomerOtp](fmt.Sprintf("otp-%s", SanitizeData(customer.Mobile)), []any{cus_otp}, 10*time.Minute) {
		log.Info("OTP: ", otp)
		utils.SendOTP(SanitizeData(customer.Mobile), otp_str)
		return true
	}
	return false
}

func (body *CustomerLoginBody) VerifyOtp() (CustomerLoginResult, bool) {
	var customer CustomerLoginResult
	var otp CustomerOtp
	var key = fmt.Sprintf("otp-%s", SanitizeData(body.Mobile))

	json_otp, err := utils.GetData(key)
	if err != nil {
		log.Info("Invalid Customer OTP")
		return customer, false
	}
	if err := json.Unmarshal([]byte(json_otp), &otp); err != nil {
		log.Error("Unable to verify OTP: ", err.Error())
		return customer, false
	}
	if SanitizeData(otp.Otp) != SanitizeData(body.Otp) {
		log.Info("Invalid Customer OTP")
		return customer, false
	}
	if utils.DeleteData(key) {
		var plain_text = []byte(*body.Mobile + time.Now().String())
		token, err := utils.GenerateToken(plain_text)
		if err != nil {
			log.Error("Error in Generating Customer Token: ", err.Error())
			return customer, false
		}
		var timestamp string = time.Now().Format("2006-01-02 15:04:05")
		customer.ID = otp.CustomerID
		customer.Name = otp.Name
		customer.Mobile = otp.Mobile
		customer.Token = &token
		customer.CreatedOn = &timestamp

		utils.SetSingleData[CustomerLoginResult](fmt.Sprintf("token-%s", token), []any{customer}, 0)
		var session_key = fmt.Sprintf("session-%s", SanitizeData(customer.Mobile))
		old_token, err := utils.GetData(session_key)
		if err != nil {
			log.Error("No existing session found: ", err.Error())
		} else {
			utils.DeleteData(fmt.Sprintf("token-%s", old_token))
			utils.DeleteData(session_key)
			//TODO: Send Push Notification to auto logout other device
		}
		utils.SetStringData(session_key, token, 0)
		return customer, true
	}
	return customer, false
}

func (customer Customer) Logout(token CustomerToken) bool {
	return utils.RunUpdate(delete_customer_token_query, []any{token.Token})
}

func AuthenticateCustomer(auth_token string) (CustomerToken, bool) {
	var token CustomerToken
	var val CustomerLoginResult
	var key = fmt.Sprintf("token-%s", auth_token)
	json_token, err := utils.GetData(key)
	if err != nil {
		return token, false
	}
	if err := json.Unmarshal([]byte(json_token), &val); err != nil {
		log.Error("Unable to verify Token: ", err.Error())
		return token, false
	}
	token.Token = val.Token
	token.CustomerID = val.ID
	token.Mobile = val.Mobile
	token.Name = val.Name
	token.CreatedOn = val.CreatedOn
	return token, true
}

/*----------------------USER TOKEN-----------------------*/
func (token UserToken) Get(authToken string) ([]any, error) {
	var fields = []any{
		&token.Token, &token.UserID, &token.CreatedOn,
		&token.Name, &token.Mobile, &token.Email,
		&token.Role,
	}
	var filter = []any{authToken}
	var query string = get_user_token_query
	return GetSingleData(filter, fields, &token, query)
}

func (token *UserToken) Add() bool {
	var fields = []any{token.Token, token.UserID}
	return utils.RunInsert(insert_user_token_query, fields)
}

func (token UserToken) Delete(auth_token string) bool {
	return utils.RunUpdate(delete_user_token_query, []any{auth_token})
}

/*----------------------CUSTOMER TOKEN-----------------------*/
func (token CustomerToken) Get(authToken string) ([]any, error) {
	var fields = []any{
		&token.Token, &token.CustomerID, &token.CreatedOn,
		&token.Name, &token.Mobile,
	}
	var filter = []any{authToken}
	var query string = get_customer_token_query
	return GetSingleData(filter, fields, &token, query)
}

func (token *CustomerToken) Add() bool {
	var fields = []any{token.Token, token.CustomerID}
	return utils.RunInsert(insert_customer_token_query, fields)
}

func (token CustomerToken) Delete(auth_token string) bool {
	return utils.RunUpdate(delete_customer_token_query, []any{auth_token})
}
