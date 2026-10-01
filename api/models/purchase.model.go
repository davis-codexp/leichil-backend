package models

type Purchase struct {
	ContentID      *int    `json:"content_id,omitempty"`
	CustomerID     *int    `json:"customer_id,omitempty"`
	Price          *int    `json:"price,omitempty"`
	Discount       *int    `json:"discount,omitempty"`
	NetAmount      *int    `json:"net_amount,omitempty"`
	PaymentMode    *string `json:"payment_mode,omitempty"`
	CreatedOn      *string `json:"created_on,omitempty"`
	CustomerName   *string `json:"customer_name,omitempty"`
	CustomerMobile *string `json:"customer_mobile,omitempty"`
	ContentName    *string `json:"content_name,omitempty"`
	ContentType    *string `json:"content_type,omitempty"`
}
