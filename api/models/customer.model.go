package models

import (
	"fmt"
	"leichil-backend/api/utils"

	"github.com/gofiber/fiber/v2/log"
)

type Customer struct {
	ID        *int    `json:"id"`
	Name      *string `json:"name,omitempty"`
	Mobile    *string `json:"mobile,omitempty"`
	Avatar    *string `json:"avatar,omitempty"`
	CreatedOn *string `json:"created_on,omitempty"`
}

type CustomerProfile struct {
	Profile      *Customer     `json:"profile"`
	Subscription *Subscription `json:"subscription"`
}

func (customer Customer) Get(offset, limit int, params map[string]string) (CountResult, []any, error) {
	var fields = get_customer_fields(&customer)
	var query string = get_customers_query
	var count_query string = count_customers_query
	filter_query, filter := customer.QueryBuilder(params)
	query += filter_query + order_filter
	count_query += filter_query
	return GetListData(filter, fields, &customer, count_query, query, offset, limit)
}

func (customer Customer) GetById(id string) ([]any, error) {
	var fields = get_customer_fields(&customer)
	var query string = get_customers_query + id_filter
	var filter = []any{id}
	return GetSingleData(filter, fields, &customer, query)
}

func (profile CustomerProfile) Get(token CustomerToken) (CustomerProfile, error) {
	var cus Customer
	var sub Subscription
	profile_result, err := cus.GetById(fmt.Sprintf(`%d`, SanitizeData(token.CustomerID)))
	if err != nil {
		log.Error("Error in getting customer profile: ", err.Error())
		return profile, err
	}
	cus = utils.FormatSingleResult[Customer](profile_result)
	profile.Profile = &cus
	subscription_result, err := sub.GetCustomerSubscription(token)
	if err != nil {
		log.Error("Error in getting subscription: ", err.Error())
		return profile, err
	}
	if len(subscription_result) > 0 {
		sub = utils.FormatSingleResult[Subscription](subscription_result)
		profile.Subscription = &sub
	}
	return profile, nil
}

func (customer Customer) IsDuplicateMobile() bool {
	var fields = get_customer_fields(&customer)
	var query string = get_customers_query + mobile_filter
	var filter = []any{customer.Mobile}
	result, err := GetSingleData(filter, fields, &customer, query)
	if err != nil {
		return false
	}
	return len(result) > 0
}

func (customer *Customer) Add() bool {
	var fields = get_customer_add_fields(customer)
	return utils.RunInsert(insert_customer_query, fields)
}

func (customer *Customer) Update(id string) bool {
	var fields = []any{customer.Name, customer.Avatar, id}
	return utils.RunUpdate(update_customer_query, fields)
}

func (customer Customer) Delete(id string) bool {
	return utils.RunUpdate(delete_customer_query, []any{id})
}

func (customer Customer) QueryBuilder(params map[string]string) (string, []any) {
	var query string
	var filter []any
	if len(params["name"]) > 0 {
		query += name_filter
		filter = append(filter, "%"+params["name"]+"%")
	}
	if len(params["mobile"]) > 0 {
		query += mobile_filter
		filter = append(filter, params["mobile"])
	}
	return query, filter
}
