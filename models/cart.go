package models

import "time"

type Session_Cart struct {
	Session_ID         uint      `json:"session_id"`
	User_ID            User      `json:"user"`
	Session_Created_At time.Time `json:"session_created_at"`
	Session_Updated_At time.Time `json:"session_updated_at"`
}

type Cart struct {
	Cart_ID uint `json:"cart_id"`
	// Session_ID      Session_Cart `json:"session_cart"`
	User_ID         int       `json:"user_id"`
	Product_ID      Product   `json:"product"`
	Product_Count   int       `json:"product_count"`
	Cart_Created_At time.Time `json:"cart_created_at"`
	Cart_Updated_At time.Time `json:"cart_updated_at"`
}
