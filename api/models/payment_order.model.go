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
	ContentID  *int    `json:"content_id,omitempty"`
	CustomerID *int    `json:"customer_id,omitempty"`
	OrderID    *string `json:"order_id,omitempty"`
	PaymentID  *int    `json:"payment_id,omitempty"`
	Amount     *int    `json:"amount,omitempty"`
	Status     *string `json:"status,omitempty"`
	CreatedOn  *string `json:"created_on,omitempty"`
}

type VerifyPaymentBody struct {
	RazorpayOrderID   *string `json:"razorpay_order_id,omitempty"`
	RazorpayPaymentID *string `json:"razorpay_payment_id,omitempty"`
	RazorpaySignature *string `json:"razorpay_signature,omitempty"`
}

func (obj *PaymentOrder) Create(token CustomerToken) map[string]any {
	var amount = 6900
	var response = make(map[string]any)
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
	return flag
}

func (obj *PaymentOrder) Update() bool {
	var fields = []any{obj.Status, obj.PaymentID, obj.OrderID}
	return utils.RunUpdate(update_payment_order_query, fields)
}
