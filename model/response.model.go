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
type UserAddressResponse struct {
	ID         uint      `json:"id"`
	HouseName  string    `json:"house_name"`
	StreetName string    `json:"stree_name"`
	Landmark   string    `json:"landmark"`
	District   string    `json:"district"`
	State      string    `json:"state"`
	Country    string    `json:"country"`
	PinCode    int       `json:"pincode"`
	Created_At time.Time `json:"created_at,omitempty"`
}

type ProductResponse struct {
	Product `json:"product"`
}
