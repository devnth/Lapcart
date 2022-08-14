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
	Full_Name    string    `json:"full_name,omitempty"`
	First_Name   string    `json:"first_name"`
	Last_Name    string    `json:"last_name"`
	Email        string    `json:"email"`
	Phone_Number int       `json:"phone_number"`
	Password     string    `json:"password,omitempty"`
	IsActive     bool      `json:"is_active,omitempty"`
	IsVerified   bool      `json:"is_verified,omitempty"`
	Created_At   time.Time `json:"created_at,omitempty"`

	// UserAddress  []UserAddressResponse `json:"address,omitempty"`
	Token string `json:"token,omitempty"`
}

//table schema for user address
type AddressResponse struct {
	Id          int    `json:"address_id,omitempty"`
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
	ID            []int64 `json:"product_id"`
	Code          string  `json:"product_code"`
	Name          string  `json:"name"`
	Description   string  `json:"description,omitempty"`
	GetBrand      `json:"brand"`
	GetProcessor  `json:"processor"`
	GetCategory   `json:"category"`
	GetColor      `json:"colors"`
	DiscountName  string  `json:"discount_name,omitempty"`
	Price         float64 `json:"price"`
	DiscountPrice float64 `json:"discount_price,omitempty"`
	Image         string  `json:"image"`
	WishList      bool    `json:"wishlist,omitempty"`
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

type GetCart struct {
	CartID        uint    `json:"cart_id,omitempty"`
	ProductID     uint    `json:"product_id,omitempty"`
	Name          string  `json:"name"`
	Brand         string  `json:"brand"`
	Category      string  `json:"category"`
	Processor     string  `json:"processor"`
	Color         string  `json:"color"`
	Count         int     `json:"count"`
	DiscountName  string  `json:"discount_name,omitempty"`
	UnitPrice     float64 `json:"unit_price,omitempty"`
	SubTotalPrice float64 `json:"sub_total_price,omitempty"`
	DiscountPrice float64 `json:"discount_price,omitempty"`
	TotalPrice    float64 `json:"total_price,omitempty"`
	Image         string  `json:"image"`
}

type GetOrders struct {
	OrderID     uint      `json:"order_id"`
	User_ID     uint      `json:"user_id"`
	User_Name   string    `json:"user_name"`
	Email       string    `json:"email"`
	TotalAmount float64   `json:"total_amount"`
	Is_Paid     bool      `json:"is_paid"`
	Status      string    `json:"order_status"`
	Ordered_At  time.Time `json:"ordered_at"`
	Updated_At  time.Time `json:"updated_at"`
}

type Orders struct {
	OrderID    int       `json:"order_id,omitempty"`
	ProductID  int       `json:"product_id"`
	Name       string    `json:"product_name"`
	Brand      string    `json:"brand"`
	Processor  string    `json:"processor"`
	Category   string    `json:"category"`
	Color      string    `json:"color"`
	Image      string    `json:"image"`
	Price      float64   `json:"price,omitempty"`
	Status     string    `json:"status"`
	Ordered_At time.Time `json:"ordere_at"`
}
