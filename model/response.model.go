package model

import "time"

type AdminResponse struct {
	ID           int    `json:"id"`
	First_Name   string `json:"first_name"`
	Last_Name    string `json:"last_name"`
	Email        string `json:"email"`
	Password     string `json:"password,omitempty"`
	Phone_Number int    `json:"phone_number"`
	Token        string `json:"token,omitempty"`
}

// user schema for user table
type UserResponse struct {
	ID           int       `json:"id"`
	First_Name   string    `json:"first_name"`
	Last_Name    string    `json:"last_name"`
	Email        string    `json:"email"`
	Phone_Number int       `json:"phone_number"`
	Password     string    `json:"password,omitempty"`
	IsActive     bool      `json:"is_active,omitempty"`
	Created_At   time.Time `json:"created_at,omitempty"`
	// UserAddress  []UserAddressResponse `json:"address,omitempty"`
	Token string `json:"token,omitempty"`
}

//table schema for user address
type AddressResponse struct {
	Id          int    `json:"address_id"`
	AddressType string `json:"address_type"`
	HouseName   string `json:"house_name"`
	StreetName  string `json:"stree_name"`
	Landmark    string `json:"landmark"`
	District    string `json:"district"`
	State       string `json:"state"`
	Country     string `json:"country"`
	PinCode     int    `json:"pincode"`
}

type ProductResponse struct {
	Product `json:"product"`
}

type GetProduct struct {
	ID           []int64 `json:"product_id"`
	Code         string  `json:"product_code"`
	Name         string  `json:"name"`
	Description  string  `json:"description,omitempty"`
	GetBrand     `json:"brand"`
	GetProcessor `json:"processor"`
	GetCategory  `json:"category"`
	GetColor     `json:"colors"`
	Price        float64 `json:"price"`
	Image        string  `json:"image"`
	WishList     bool    `json:"wishlist,omitempty"`
}

type GetBrand struct {
	ID   uint   `json:"id,omitempty"`
	Name string `json:"name"`
}

// table schema for product_category
type GetCategory struct {
	ID          uint   `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

//table schema for product color
type GetColor struct {
	Name []string `json:"name"`
}

type GetProcessor struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name"`
}

//table schema for product_discount
type GetDiscount struct {
	ID          uint    `json:"id,omitempty"`
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	Percentage  float32 `json:"percentage,omitempty"`
	Status      bool    `json:"status,omitempty"`
}
