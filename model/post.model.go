package model

// user schema for user table
type User struct {
	First_Name   string `json:"first_name"`
	Last_Name    string `json:"last_name"`
	Password     string `json:"password"`
	Email        string `json:"email"`
	Phone_Number int    `json:"phone_number"`
}

//table schema for user address
type UserAddress struct {
	AddressType string `json:"address_type"`
	HouseName   string `json:"house_name"`
	StreetName  string `json:"stree_name"`
	Landmark    string `json:"landmark"`
	District    string `json:"district"`
	State       string `json:"state"`
	Country     string `json:"country"`
	PinCode     int    `json:"pincode"`
}

type Admin struct {
	First_Name   string `json:"first_name"`
	Last_Name    string `json:"last_name"`
	Password     string `json:"password"`
	Email        string `json:"email"`
	Phone_Number int    `json:"phone_number"`
}
