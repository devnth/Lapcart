package models

import (
	"time"
)

// user schema for user table
type User struct {
	User_ID        int           `json:"user_id"`
	First_Name     string        `json:"first_name"`
	Last_Name      string        `json:"last_name"`
	Password       string        `json:"password"`
	Email          string        `json:"email"`
	Phone_Number   int           `json:"phone_number"`
	IsActive       bool          `json:"is_active"`
	UserCart       []ProductUser `json:"user_cart"`
	AddressDetails []UserAddress `json:"address_details"`
	OrderStatus    []Order       `json:"order"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
	DeletedAt      time.Time     `json:"deleted_at"`
}

//table schema for user address
type UserAddress struct {
	UserAddress_ID uint   `json:"user_address_id"`
	HouseName      string `json:"house_name"`
	StreetName     string `json:"stree_name"`
	Landmark       string `json:"landmark"`
	District       string `json:"district"`
	State          string `json:"state"`
	Country        string `json:"country"`
	PinCode        int    `json:"pincode"`
}
