package model

import (
	"encoding/json"
)

type RequestCreateOrder struct {
	UserID        string          `json:"user_id" validate:"required"`
	PaymentTypeID string          `json:"payment_type_id" validate:"required"`
	OrderNumber   string          `json:"order_number" validate:"required"`
	TotalPrice    float64         `json:"total_price" validate:"required"`
	ProductOrder  json.RawMessage `json:"product_order"`
	Status        string          `json:"status" validate:"required"`
	IsPaid        bool            `json:"is_paid"`
}
