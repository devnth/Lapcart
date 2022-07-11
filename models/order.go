package models

import "time"

// table schema for Order
type Order struct {
	Order_ID        uint          `json:"order_id"`
	Ordered_Cart    []ProductUser `json:"ordered_cart"`
	Ordered_At      time.Time     `json:"ordered_at_time"`
	Total_Price     float64       `json:"total_price"`
	Discount_Amount float64       `json:"discount_amount"`
	Payment_Method  Payment       `json:"payment_method"`
}
