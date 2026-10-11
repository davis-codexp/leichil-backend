package models

import (
	"encoding/hex"
	"hash"
	"leichil-backend/api/utils"
	"os"

	"crypto/hmac"
	"crypto/sha256"

	"github.com/gofiber/fiber/v2/log"
	razorpay "github.com/razorpay/razorpay-go"
)

var rzp_client *razorpay.Client
var rzp_secret string
var rzp_key string

type PaymentOrder struct {
	ContentID      *int    `json:"content_id,omitempty"`
	CustomerID     *int    `json:"customer_id,omitempty"`
	PlanID         *int    `json:"plan_id,omitempty"`
	OrderID        *string `json:"order_id,omitempty"`
	PaymentID      *string `json:"payment_id,omitempty"`
	Amount         *int    `json:"amount,omitempty"`
	Status         *string `json:"status,omitempty"`
	CreatedOn      *string `json:"created_on,omitempty"`
	CustomerName   *string `json:"customer_name,omitempty"`
	CustomerMobile *string `json:"customer_mobile,omitempty"`
	ContentName    *string `json:"content_name,omitempty"`
}

type VerifyPaymentBody struct {
	RazorpayOrderID   *string `json:"razorpay_order_id,omitempty"`
	RazorpayPaymentID *string `json:"razorpay_payment_id,omitempty"`
	RazorpaySignature *string `json:"razorpay_signature,omitempty"`
	PlanID            *int    `json:"plan_id,omitempty"`
}

func (obj PaymentOrder) Get(offset, limit int, params map[string]string) (CountResult, []any, error) {
	var fields = get_payment_order_fields(&obj)
	var query string = get_payment_orders_query
	filter_query, filter := obj.QueryBuilder(params)
	query += filter_query + " ORDER BY po.created_on DESC"
	return GetListData(filter, fields, &obj, "", query, offset, limit)
}

func (obj PaymentOrder) GetCustomerPurchase(token CustomerToken) (CountResult, []any, error) {
	var fields = get_payment_order_fields(&obj)
	var query string = get_payment_orders_query + payment_customer_filter + payment_order_status_filter
	var filter = []any{token.CustomerID, "SUCCESS"}
	return GetListData(filter, fields, &obj, "", query, 0, 100000)
}

func (obj PaymentOrder) VerifyCustomerPurchase(token CustomerToken, content_id string) bool {
	var fields = get_payment_order_fields(&obj)
	var query string = get_payment_orders_query + payment_customer_filter + payment_content_filter + payment_order_status_filter
	var filter = []any{token.CustomerID, content_id, "SUCCESS"}
	_, result, err := GetListData(filter, fields, &obj, "", query, 0, 100000)
	if err != nil {
		log.Error("Error in Verifying customer purchase: ", err.Error())
		return false
	}
	return len(result) > 0
}

func (obj *PaymentOrder) Create(token CustomerToken, order_type string) map[string]any {
	var response = make(map[string]any)
	var amount int

	if order_type == "payment" && obj.ContentID != nil {
		var content Content
		var content_fields = []any{&content.ID, &content.Price}
		var filter = []any{obj.ContentID}
		result, err := GetSingleData(filter, content_fields, &content, get_content_price_query)
		if err != nil {
			log.Error("Error in Getting Content By ID: ", err.Error())
			return response
		}
		content = utils.FormatSingleResult[Content](result)
		amount = SanitizeData(content.Price) * 100
	}

	if order_type == "subscription" && obj.PlanID != nil {
		var plan SubscriptionPlan
		var plan_fields []any = get_subscription_plan_fields(&plan)
		var filter = []any{obj.PlanID}
		var query string = get_subscription_plans_query + subscription_plan_id_filter
		result, err := GetSingleData(filter, plan_fields, &plan, query)
		if err != nil {
			log.Error("Error in Getting Subscription Plan By ID: ", err.Error())
			return response
		}
		plan = utils.FormatSingleResult[SubscriptionPlan](result)
		amount = SanitizeData(plan.Price) * 100
	}

	rzp_key = os.Getenv("RZP_KEY_ID")
	rzp_secret = os.Getenv("RZP_KEY_SECRET")
	rzp_client = razorpay.NewClient(rzp_key, rzp_secret)
	var data = map[string]interface{}{
		"amount":   amount,
		"currency": "INR",
		"receipt":  "receipt_srv_01",
	}
	body, err := rzp_client.Order.Create(data, nil)
	if err != nil {
		log.Error("Unable to create razorpay order: ", err.Error())
		return response
	}
	var order_id = body["id"].(string)
	obj.OrderID = &order_id
	obj.Amount = &amount
	obj.CustomerID = token.CustomerID
	var fields = get_payment_order_add_fields(obj)
	if utils.RunInsert(insert_payment_order_query, fields) {
		response = map[string]any{
			"order_id": order_id,
			"amount":   amount,
			"currency": "INR",
		}
	} else {
		log.Info("Error while inserting payment order into DB")
	}
	return response
}

func (obj *VerifyPaymentBody) VerifyPayment(token CustomerToken, order_type string) bool {
	var flag = true
	var status = "SUCCESS"
	var signature_payload = SanitizeData(obj.RazorpayOrderID) + "|" + SanitizeData(obj.RazorpayPaymentID)
	var h hash.Hash = hmac.New(sha256.New, []byte(rzp_secret))
	h.Write([]byte(signature_payload))
	var signature = hex.EncodeToString(h.Sum(nil))
	if signature != SanitizeData(obj.RazorpaySignature) {
		log.Info("Invalid transaction signature")
		status = "FAILED"
		flag = false
	}
	var fields = []any{status, obj.RazorpayPaymentID, obj.RazorpayOrderID}
	utils.RunUpdate(update_payment_order_query, fields)
	if order_type == "subscription" && obj.PlanID != nil {
		var sub = new(Subscription)
		sub.CustomerID = token.CustomerID
		sub.PlanID = obj.PlanID
		sub.Add(token)
	} else if order_type == "payment" {
		var sub Subscription
		var sub_query = get_subscriptions_query + subscription_customer_filter + subscription_not_expired_filter + " LIMIT 1"
		var sub_fields = get_subscription_fields(&sub)
		var sub_filter = []any{token.CustomerID}
		result, err := GetSingleData(sub_filter, sub_fields, &sub, sub_query)
		if err != nil {
			log.Error("Error in Getting Active Customer Subscriptions: ", err.Error())
			return flag
		}
		if len(result) > 0 {
			//Extend Subscription by 1 month
			sub = utils.FormatSingleResult[Subscription](result)
			utils.RunUpdate(extend_subscription_expiry_query, []any{30, sub.ID})
		} else {
			//Create New Subscription
			var plan SubscriptionPlan
			var plan_query string = get_subscription_plans_query + subscription_plan_name_filter
			var plan_fields []any = get_subscription_plan_fields(&plan)
			result, err := GetSingleData([]any{"BASIC"}, plan_fields, &plan, plan_query)
			if err != nil {
				log.Error("Error in getting BASIC plan ID: ", err.Error())
				return flag
			}
			plan = utils.FormatSingleResult[SubscriptionPlan](result)
			var sub = new(Subscription)
			sub.CustomerID = token.CustomerID
			sub.PlanID = plan.ID
			sub.Add(token)
		}
	}
	return flag
}

func (obj *PaymentOrder) Update() bool {
	var fields = []any{obj.Status, obj.PaymentID, obj.OrderID}
	return utils.RunUpdate(update_payment_order_query, fields)
}

func (obj PaymentOrder) QueryBuilder(params map[string]string) (string, []any) {
	var query string
	var filter []any
	if len(params["order_id"]) > 0 {
		query += payment_order_id_filter
		filter = append(filter, params["order_id"])
	}
	if len(params["payment_id"]) > 0 {
		query += payment_order_payment_filter
		filter = append(filter, params["payment_id"])
	}
	if len(params["amount"]) > 0 {
		query += payment_order_amount_filter
		filter = append(filter, params["amount"])
	}
	if len(params["status"]) > 0 {
		query += payment_order_status_filter
		filter = append(filter, params["status"])
	}
	if len(params["content_id"]) > 0 {
		query += payment_content_filter
		filter = append(filter, params["content_id"])
	}
	if len(params["customer_id"]) > 0 {
		query += payment_customer_filter
		filter = append(filter, params["customer_id"])
	}
	return query, filter
}
