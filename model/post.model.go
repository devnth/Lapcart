package model

import "time"

// user schema for user table
type User struct {
	First_Name   string `json:"first_name"`
	Last_Name    string `json:"last_name"`
	Password     string `json:"password"`
	Email        string `json:"email"`
	Phone_Number int    `json:"phone_number"`
	IsActive     bool   `json:"is_active,omitempty"`
}

//table schema for user address
type Address struct {
	User_id     int    `json:"user_id,omitempty"`
	AddressType string `json:"address_type"`
	HouseName   string `json:"house_name"`
	StreetName  string `json:"stree_name"`
	Landmark    string `json:"landmark"`
	District    string `json:"district"`
	State       string `json:"state"`
	Country     string `json:"country"`
	PinCode     int    `json:"pincode"`
	Created_At  time.Time
}

type Admin struct {
	First_Name   string `json:"first_name"`
	Last_Name    string `json:"last_name"`
	Password     string `json:"password"`
	Email        string `json:"email"`
	Phone_Number int    `json:"phone_number"`
}

type Product struct {
	ID          []uint    `json:"id,omitempty"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Brand       Brand     `json:"brand"`
	Processor   Processor `json:"processor"`
	Category    Category  `json:"category"`
	Colors      []Color   `json:"colors"`
	Price       float64   `json:"price"`
	// Discount    Discount  `json:"discount,omitempty"`
	Rating    float32 `json:"rating"`
	Image     string  `json:"image"`
	IsDeleted bool    `json:"is_deleted,omitempty"`
}

// table schema for product_category
type Category struct {
	ID          uint   `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

//table scehma for product_branding
type Brand struct {
	ID          uint   `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

//table schema for product color
type Color struct {
	Name  string `json:"name"`
	Stock int    `json:"quantity"`
}

//tabel schema for product_processor
type Processor struct {
	ID          uint   `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

//table schema for product_discount
type Discount struct {
	ArgId        int
	ID           int
	AllProducts  bool      `json:"all_products"`
	ProductCode  string    `json:"product_code"`
	Brand        string    `json:"brand"`
	Category     string    ` json:"category"`
	DiscountName string    `json:"discount_name"`
	Percentage   float32   `json:"percentage"`
	Status       bool      `json:"status"`
	Created_At   time.Time `json:"created_at"`
	Expiry_Date  time.Time `json:"expiry_date"`
	Updated_At   time.Time `json:"updated_at"`
	Deleted_At   time.Time `json:"deleted_at"`
}

type WishList struct {
	User_Id     int
	ProductCode string `json:"product_code"`
}

type Name struct {
	Name string `json:"name"`
}

type ProductCode struct {
	ProductCode string
}

type Filter struct {
	Category    []Category    `json:"category"`
	Brand       []Brand       `json:"brand"`
	Color       []Color       `json:"color"`
	Processor   []Processor   `json:"processor"`
	Name        []Name        `json:"name"`
	ProductCode []ProductCode `json:"product_code"`
}
