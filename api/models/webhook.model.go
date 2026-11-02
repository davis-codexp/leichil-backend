package models

import "encoding/json"

type RazorpayEvent struct {
	Entity    string          `json:"entity"`
	AccountID string          `json:"account_id"`
	Event     string          `json:"event"`
	Contains  []string        `json:"contains"`
	Payload   json.RawMessage `json:"payload"`
	CreatedAt int64           `json:"created_at"`
}

type PaymentEntity struct {
	ID       string `json:"id"`
	Entity   string `json:"entity"`
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
	Status   string `json:"status"`
	OrderID  string `json:"order_id"`
	Method   string `json:"method"`
	Email    string `json:"email"`
	Contact  string `json:"contact"`
}
