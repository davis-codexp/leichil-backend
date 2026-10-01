package models

import (
	"leichil-backend/api/utils"
	"time"

	"github.com/gofiber/fiber/v2/log"
)

type SubscriptionPlan struct {
	ID             *int    `json:"id"`
	Name           *string `form:"name" json:"name"`
	Features       *string `form:"features" json:"features"`
	Price          *int    `form:"price" json:"price"`
	DurationMonths *int    `form:"duration_months" json:"duration_months"`
	IsActive       *int    `form:"is_active" json:"is_active"`
	CreatedOn      *string `json:"created_on"`
	CreatedBy      *int    `json:"created_by"`
	CreatedByName  *string `json:"created_by_name"`
}

func (obj SubscriptionPlan) Get(offset, limit int, params map[string]string) (CountResult, []any, error) {
	var fields = get_subscription_plan_fields(&obj)
	var query string = get_subscription_plans_query
	var count_query string = count_subscription_plans_query
	filter_query, filter := obj.QueryBuilder(params)
	query += filter_query
	count_query += filter_query
	return GetListData(filter, fields, &obj, count_query, query, offset, limit)
}

func (obj SubscriptionPlan) GetById(id string) ([]any, error) {
	var fields = get_subscription_plan_fields(&obj)
	var query string = get_subscription_plans_query + subscription_plan_id_filter
	var filter = []any{id}
	return GetSingleData(filter, fields, &obj, query)
}

func (obj *SubscriptionPlan) Add(token UserToken) bool {
	var fields = get_subscription_plan_add_fields(obj)
	fields = append(fields, token.UserID)
	return utils.RunInsert(insert_subscription_plan_query, fields)
}

func (obj *SubscriptionPlan) Update(id string) bool {
	var fields = get_subscription_plan_add_fields(obj)
	fields = append(fields, id)
	return utils.RunUpdate(update_subscription_plan_query, fields)
}

func (obj SubscriptionPlan) Delete(id string) bool {
	return utils.RunUpdate(delete_subscription_plan_query, []any{id})
}

func (obj SubscriptionPlan) QueryBuilder(params map[string]string) (string, []any) {
	var query string
	var filter []any
	if len(params["name"]) > 0 {
		query += subscription_plan_name_filter
		filter = append(filter, "%"+params["name"]+"%")
	}
	if len(params["is_active"]) > 0 {
		query += subscription_plan_status_filter
		filter = append(filter, params["is_active"])
	}
	if len(params["price"]) > 0 {
		query += subscription_plan_price_filter
		filter = append(filter, params["price"])
	}
	if len(params["duration"]) > 0 {
		query += subscription_plan_duration_filter
		filter = append(filter, params["duration"])
	}
	return query, filter
}

/*--------------------- SUBSCRIPTIONS ---------------------*/

type Subscription struct {
	ID             *int    `json:"id"`
	CustomerID     *int    `form:"customer_id" json:"customer_id"`
	PlanID         *int    `form:"plan_id" json:"plan_id"`
	Price          *int    `form:"price" json:"price"`
	ExpiresOn      *string `json:"expires_on"`
	CreatedOn      *string `json:"created_on"`
	CustomerMobile *string `json:"customer_mobile"`
	CustomerName   *string `json:"customer_name"`
	CustomerAvatar *string `json:"customer_avatar,omitempty"`
	PlanName       *string `json:"plan_name"`
}

func (obj Subscription) Get(offset, limit int, params map[string]string) (CountResult, []any, error) {
	var fields = get_subscription_fields(&obj)
	var query string = get_subscriptions_query
	var count_query string = count_subscriptions_query
	filter_query, filter := obj.QueryBuilder(params)
	query += filter_query
	count_query += filter_query
	return GetListData(filter, fields, &obj, count_query, query, offset, limit)
}

func (obj Subscription) GetById(id string) ([]any, error) {
	var fields = get_subscription_fields(&obj)
	var query string = get_subscriptions_query + subscription_id_filter
	var filter = []any{id}
	return GetSingleData(filter, fields, &obj, query)
}

func (obj Subscription) GetCustomerSubscription(token CustomerToken) ([]any, error) {
	var query string = get_subscriptions_query + subscription_customer_filter + subscription_not_expired_filter
	var fields = get_subscription_fields(&obj)
	var filter = []any{token.CustomerID}
	return GetSingleData(filter, fields, &obj, query)
}

func (obj *Subscription) Add(token CustomerToken) bool {
	var sub_plan SubscriptionPlan
	var fields = get_subscription_add_fields(obj)
	var sub_plan_fields = get_subscription_plan_fields(&sub_plan)
	var sub_plan_query string = get_subscription_plans_query + subscription_plan_id_filter
	var expires_on string
	var days int

	result, err := GetSingleData([]any{obj.PlanID}, sub_plan_fields, &sub_plan, sub_plan_query)
	if err != nil {
		log.Error("Error in Getting Subscription Plan By ID: ", err.Error())
		return false
	}

	sub_plan = utils.FormatSingleResult[SubscriptionPlan](result)
	var duration int = SanitizeData(sub_plan.DurationMonths)
	days = duration * 30
	if duration == 12 {
		days += 5
	}
	expires_on = time.Now().AddDate(0, 0, days).Format("2006-01-02")
	fields = append(fields, sub_plan.Price, token.CustomerID, expires_on)
	return utils.RunInsert(insert_subscription_query, fields)
}

func (obj *Subscription) Update(id string) bool {
	var sub_plan SubscriptionPlan
	var fields = get_subscription_add_fields(obj)
	var sub_plan_fields = get_subscription_plan_fields(&sub_plan)
	var sub_plan_query string = get_subscription_plans_query + subscription_plan_id_filter
	result, err := GetSingleData([]any{obj.PlanID}, sub_plan_fields, &sub_plan, sub_plan_query)
	if err != nil {
		log.Error("Error in Getting Subscription Plan By ID: ", err.Error())
		return false
	}

	sub_plan = utils.FormatSingleResult[SubscriptionPlan](result)
	fields = append(fields, sub_plan.Price, obj.ExpiresOn, id)
	return utils.RunUpdate(update_subscription_query, fields)
}

func (obj Subscription) Delete(id string) bool {
	return utils.RunUpdate(delete_subscription_query, []any{id})
}

func (obj Subscription) QueryBuilder(params map[string]string) (string, []any) {
	var query string
	var filter []any
	if len(params["customer"]) > 0 {
		query += subscription_customer_filter
		filter = append(filter, params["customer"])
	}
	if len(params["plan"]) > 0 {
		query += subscription_plan_filter
		filter = append(filter, params["plan"])
	}
	if len(params["date"]) > 0 {
		query += subscription_date_filter
		filter = append(filter, params["date"])
	}
	if len(params["expired"]) > 0 {
		if params["expired"] == "true" {
			query += subscription_expired_filter
		} else if params["expired"] == "false" {
			query += subscription_not_expired_filter
		}
	}
	return query, filter
}
