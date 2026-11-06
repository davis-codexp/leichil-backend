package models

import (
	"encoding/hex"
	"hash"
	"leichil-backend/api/utils"
	"os"
	"strings"
	"time"

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
	Thumbnail      *string `json:"thumbnail,omitempty"`
}

type VerifyPaymentBody struct {
	RazorpayOrderID   *string `json:"razorpay_order_id,omitempty"`
	RazorpayPaymentID *string `json:"razorpay_payment_id,omitempty"`
	RazorpaySignature *string `json:"razorpay_signature,omitempty"`
	PlanID            *int    `json:"plan_id,omitempty"`
}

type PaymentQRBody struct {
	ContentID   *int    `json:"content_id,omitempty"`
	PlanID      *int    `json:"plan_id,omitempty"`
	PaymentType *string `json:"payment_type,omitempty"`
	DeviceID    *string `json:"device_id,omitempty"`
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

	if len(rzp_key) == 0 {
		rzp_key = os.Getenv("RZP_KEY_ID")
	}
	if len(rzp_secret) == 0 {
		rzp_secret = os.Getenv("RZP_KEY_SECRET")
	}
	if rzp_client == nil {
		rzp_client = razorpay.NewClient(rzp_key, rzp_secret)
	}
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

func (obj *VerifyPaymentBody) VerifyPayment() bool {
	var flag = true
	var status = "SUCCESS"
	var signature_payload = SanitizeData(obj.RazorpayOrderID) + "|" + SanitizeData(obj.RazorpayPaymentID)
	log.Info("RAZORPAY ORDER ID: ", SanitizeData(obj.RazorpayOrderID))
	log.Info("RAZORPAY PAYMENT ID: ", SanitizeData(obj.RazorpayPaymentID))
	var h hash.Hash = hmac.New(sha256.New, []byte(rzp_secret))
	h.Write([]byte(signature_payload))
	var signature = hex.EncodeToString(h.Sum(nil))
	if signature != SanitizeData(obj.RazorpaySignature) {
		log.Info("Invalid transaction signature: ", signature, SanitizeData(obj.RazorpaySignature))
		status = "FAILED"
		flag = false
	}
	var fields = []any{status, obj.RazorpayPaymentID, obj.RazorpayOrderID}
	utils.RunUpdate(update_payment_order_query, fields)
	return flag
}

func (po PaymentOrder) UpdateOrder(status, payment_id, order_id string) bool {
	var fields = []any{status, payment_id, order_id}
	return utils.RunUpdate(update_payment_order_query, fields)
}

func (po PaymentOrder) UpdateSubscription(order_id string) {
	payment_type, po := get_payment_type_by_id(order_id)
	switch strings.ToUpper(payment_type) {
	case "SUBSCRIPTION":
		var sub = new(Subscription)
		var token = CustomerToken{CustomerID: po.CustomerID}
		sub.CustomerID = po.CustomerID
		sub.PlanID = po.PlanID
		sub.Add(token)
	case "PURCHASE":
		complimentary_subscription(po)
	default:
		log.Info("Invalid Payment Type: ", payment_type)
	}
}

func complimentary_subscription(po PaymentOrder) {
	var sub Subscription
	var sub_query = get_subscriptions_query + subscription_customer_filter + subscription_not_expired_filter + " LIMIT 1"
	var sub_fields = get_subscription_fields(&sub)
	var sub_filter = []any{po.CustomerID}
	result, err := GetSingleData(sub_filter, sub_fields, &sub, sub_query)
	if err != nil {
		log.Error("Error in Getting Active Customer Subscriptions: ", err.Error())
		return
	}
	if len(result) > 0 {
		sub = utils.FormatSingleResult[Subscription](result)
		utils.RunUpdate(extend_subscription_expiry_query, []any{30, sub.ID})
		var list = get_device_tokens(SanitizeData(po.CustomerID))
		if len(list) > 0 {
			for _, device_id := range list {
				utils.SendPushNotification(device_id, "Subscription Extened", "Your BASIC Subscription is extended by 1 month as a complimentary offer with your purchase.", map[string]string{})
			}
		}
	} else {
		var plan SubscriptionPlan
		var plan_query string = get_subscription_plans_query + subscription_plan_name_filter
		var plan_fields []any = get_subscription_plan_fields(&plan)
		result, err := GetSingleData([]any{"BASIC"}, plan_fields, &plan, plan_query)
		if err != nil {
			log.Error("Error in getting BASIC plan ID: ", err.Error())
			return
		}
		plan = utils.FormatSingleResult[SubscriptionPlan](result)
		var sub = new(Subscription)
		var token = CustomerToken{CustomerID: po.CustomerID}
		sub.CustomerID = po.CustomerID
		sub.PlanID = plan.ID
		sub.Add(token)
		var list = get_device_tokens(SanitizeData(po.CustomerID))
		if len(list) > 0 {
			for _, device_id := range list {
				utils.SendPushNotification(device_id, "Subscription Unlocked", "1 Month BASIC Subscription Activated as complimentary offer with your purchase.", map[string]string{})
			}
		}
	}
}

func get_payment_type_by_id(order_id string) (string, PaymentOrder) {
	var po PaymentOrder
	var fields = get_payment_order_fields(&po)
	var query string = get_payment_orders_query + payment_order_id_filter
	var filter = []any{order_id}
	result, err := GetSingleData(filter, fields, &po, query)
	if err != nil {
		log.Error("Unable to Get Payment Type: ", err.Error())
		return "", po
	}
	po = utils.FormatSingleResult[PaymentOrder](result)
	if po.PlanID != nil {
		return "SUBSCRIPTION", po
	} else if po.ContentID != nil {
		return "PURCHASE", po
	}
	return "", po
}

func get_device_tokens(customer_id int) []string {
	var list []string
	if customer_id > 0 {
		var cus Customer
		var cus_query = get_customers_query + id_filter
		var cus_fields = get_customer_fields(&cus)
		result, err := GetSingleData([]any{customer_id}, cus_fields, &cus, cus_query)
		if err != nil {
			log.Error("Unable to get customer data: ", err.Error())
			return list
		}
		cus = utils.FormatSingleResult[Customer](result)
		var mobile = SanitizeData(cus.Mobile)
		if len(mobile) > 0 {
			list, err = utils.GetDevicesByPhone(mobile)
			if err != nil {
				log.Error("Error in Getting device tokens by phone: ", err.Error())
				return list
			}
		}
	}
	return list
}

func (obj *PaymentOrder) Update() bool {
	var fields = []any{obj.Status, obj.PaymentID, obj.OrderID}
	return utils.RunUpdate(update_payment_order_query, fields)
}

func GetPaymentLink() string {
	var link string
	//var callback_url = os.Getenv("RZP_CALLBACK_URL")
	if len(rzp_key) == 0 {
		rzp_key = os.Getenv("RZP_KEY_ID")
	}
	if len(rzp_secret) == 0 {
		rzp_secret = os.Getenv("RZP_KEY_SECRET")
	}
	if rzp_client == nil {
		rzp_client = razorpay.NewClient(rzp_key, rzp_secret)
	}
	var data = map[string]interface{}{
		"amount":         9900,
		"currency":       "INR",
		"accept_partial": false,
		"description":    "Test Payment for TV",
		"reference_id":   "order_5565",
		"customer": map[string]interface{}{
			"name":   "Kakashi Hatake",
			"email":  "kakashi@konoha.com",
			"mobile": "9090909090",
		},
		//"callback_url":    callback_url,
		"callback_method": "get",
		"expire_by":       time.Now().Add(20 * time.Minute).Unix(),
	}
	payment_link, err := rzp_client.PaymentLink.Create(data, nil)
	if err != nil {
		log.Error("Error in Generating Payment Link: ", err.Error())
		return link
	}
	link, ok := payment_link["short_url"].(string)
	if !ok || link == "" {
		log.Error("Invalid response from Razorpay")
		return link
	}
	return link
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
	if len(params["start_date"]) > 0 && len(params["end_date"]) > 0 {
		query += payment_date_range_filter
		filter = append(filter, params["start_date"], params["end_date"])
	}
	return query, filter
}
