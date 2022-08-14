package model

import "time"

type UserRequest struct {
	First_Name   string `json:"first_name" validate:"required,min=2,max=50,alpha"`
	Last_Name    string `json:"last_name" validate:"required,alpha"`
	Email        string `json:"email" validate:"required,email"`
	Password     string `json:"password" validate:"required,min=6"`
	Phone_Number int    `json:"phone_number" validate:"required"`
}

// user schema for user table
type User struct {
	ID                int       `json:"user_id,omitempty"`
	First_Name        string    `json:"first_name"`
	Last_Name         string    `json:"last_name"`
	Password          string    `json:"password"`
	Email             string    `json:"email"`
	Phone_Number      int       `json:"phone_number"`
	IsActive          bool      `json:"is_active,omitempty"`
	IsVerified        bool      `json:"is_verified,omitempty"`
	Created_At        time.Time `json:"created_at"`
	Updated_At        time.Time `json:"updated_at"`
	Verification_Code int       `json:"verification_code,omitempty"`
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
	ProductID   int       `json:"product_id,omitempty"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Brand       Brand     `json:"brand"`
	Processor   Processor `json:"processor"`
	Category    Category  `json:"category"`
	Colors      []Color   `json:"colors"`
	Price       float64   `json:"price"`
	Rating      float32   `json:"rating"`
	Image       string    `json:"image"`
	IsDeleted   bool      `json:"is_deleted,omitempty"`
	Updated_At  time.Time `json:"updated_at,omitempty"`
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
	Category    []string `json:"category"`
	Brand       []string `json:"brand"`
	Color       []string `json:"color"`
	Processor   []string `json:"processor"`
	Name        []string `json:"name"`
	ProductCode []string `json:"product_code"`
	Sort
	PriceRange
}

type Sort struct {
	Name   string
	Price  string
	Latest string
}

type PriceRange struct {
	Min int
	Max int
}

type Cart struct {
	ID         int `json:"cart_id"`
	User_Id    int
	Product_Id uint `json:"product_id"`
	Count      int  `json:"count"`
	Created_At time.Time
	Updated_At time.Time
}

type Coupon struct {
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	Min_Amount  float64   `json:"min_amount"`
	Value       float64   `json:"coupon_value"`
	Valid_Till  time.Time `json:"valid_till"`
	Created_At  time.Time `json:"created_at"`
}

type OrderDetails struct {
	ID         uint      `json:"id"`
	User_ID    uint      `json:"user_id"`
	Address_ID uint      `json:"address_id"`
	TotalPrice float64   `json:"total_price"`
	Status     string    `json:"status"`
	Created_At time.Time `json:"created_at"`
	Updated_At time.Time `json:"updated_at"`
}

type OrderItems struct {
	ID         uint      `json:"id"`
	OrderID    uint      `json:"order_id"`
	ProductID  uint      `json:"product_id"`
	DiscountID uint      `json:"discount_id"`
	Quantity   int       `json:"quantity"`
	Created_At time.Time `json:"created_at"`
}

type Payment struct {
	ID                  uint `json:"payment_id"`
	Order_ID            uint `json:"order_id"`
	User_ID             int  `json:"user_id"`
	Full_Name           string
	Email               string
	Phone_Number        int
	Status              bool      `json:"status"`
	PaymentType         string    `json:"payment_type"`
	Amount              float64   `json:"amount"`
	Coupon_Code         string    `json:"coupon_code"`
	Coupon_Id           uint      `json:"coupon_id"`
	Created_At          time.Time `json:"created_at"`
	Updated_At          time.Time `json:"updated_at"`
	Razorpay_payment_id string    `json:"razor_payment_id"`
	Razorpay_order_id   string    `json:"razor_order_id"`
	Razorpay_signature  string    `json:"razor_signature"`
}

type ManageOrder struct {
	Order_ID   uint      `json:"order_id"`
	Status     string    `json:"status"`
	Updated_At time.Time `json:"updated_at"`
}

type UpdateProduct struct {
	ProductID      int       `json:"product_id"`
	Code           string    `json:"new_item_code"`
	OldCode        string    `json:"item_code"`
	Name           string    `json:"product_name"`
	Brand          string    `json:"brand_name"`
	Category       string    `json:"category"`
	Processor      string    `json:"processor"`
	ChangeColor    string    `json:"change_color"`
	ChangeQuantity int       `json:"change_quantity"`
	NewColor       string    `json:"new_color"`
	NewQuantity    int       `json:"new_quantity"`
	Price          float64   `json:"price"`
	Image          string    `json:"image"`
	Updated_At     time.Time `json:"updated_at"`
}

type DeleteProduct struct {
	ProductId    int       `json:"product_id"`
	Product_Code string    `json:"product_code"`
	Deleted_At   time.Time `json:"deleted_at"`
}
