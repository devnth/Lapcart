package models

import "time"

// product schema for product table
type Product struct {
	Product_ID          uint              `json:"product_id"`
	Product_Name        string            `json:"product_name"`
	Product_Description string            `json:"product_description"`
	Brand               Product_Branding  `json:"product_brand"`
	Processor           Product_Processor `json:"product_processor"`
	Category            Product_Category  `json:"product_category"`
	Product_Price       float64           `json:"product_price"`
	Discount            Product_Discount  `json:"product_discount"`
	Product_Rating      float32           `json:"product_rating"`
	Product_Image       string            `json:"product_image"`
	Product_Stock       Product_Inventory `json:"product_stock"`
	Product_Created_At  time.Time         `json:"product_created_at"`
	// Product_Modified_At time.Time         `json:"product_modified_at"`
	// Product_Deleted_At  time.Time         `json:"product_deleted_at"`
}

// type for usercart
type ProductUser struct {
	ProductUser_ID          uint              `json:"product_id"`
	ProductUser_Name        string            `json:"product_name"`
	ProductUser_Description string            `json:"product_description"`
	Brand                   Product_Branding  `json:"product_brand"`
	Processor               Product_Processor `json:"product_processor"`
	Category                Product_Category  `json:"product_category"`
	Product_Price           float64           `json:"product_price"`
	ProductUser_Count       uint              `json:"product_count"`
	Discount                Product_Discount  `json:"product_discount"`
	ProductUser_Image       string            `json:"product_image"`
}

// table schema for product_category
type Product_Category struct {
	Category_ID          uint      `json:"category_id"`
	Category_Name        string    `json:"category_name"`
	Category_Description string    `json:"category_description"`
	Category_Created_At  time.Time `json:"category_created_at"`
	// Category_Updated_At  time.Time `json:"category_updated-at"`
	// Category_Deleted_At  time.Time `json:"category_deleted_at"`
}

//table scehma for product_branding
type Product_Branding struct {
	Product_Brand_Id         uint      `json:"product_brand_id"`
	Product_Brand_Name       string    `json:"product_brand_name"`
	Product_Brand_Created_At time.Time `json:"product_brnad_created_at"`
	// Product_Brand_Updated_At time.Time `json:"product_brand_updated_at"`
	// Product_Brand_Deleted_At time.Time `json:"product_brand_deleted_at"`
}

//table schema for product color
type Product_Color struct {
	Product_Color_ID   uint    `json:"product_color_id"`
	Product_Color_Name string  `json:"product_color_name"`
	Product_ID         Product `json:"product"`
}

//tabel schema for product_processor
type Product_Processor struct {
	Product_Processor_ID          uint      `json:"product_processor_id"`
	Product_Processor_Name        string    `json:"product_processor_name"`
	Product_Processor_Description string    `json:"product_processor_description"`
	Product_Processor_Created_At  time.Time `json:"product_processor_created_at"`
	Product_Processor_Updated_At  time.Time `json:"product_processor_updated_at"`
	Product_Processor_Deleted_At  time.Time `json:"product_processor_deleted_at"`
}

//table schema for product inventory
type Product_Inventory struct {
	Product_Inventory_ID     uint      `json:"product_inventory-id"`
	Product_Quantity         int       `json:"product_stock_count"`
	Product_Stock_Created_At time.Time `json:"product_stocked_at"`
	// Product_Stock_Updated_At time.Time `json:"product_stock_updated_at"`
	// Product_Stock_Deleted_At time.Time `json:"product_stock_deleted_at"`
}

//table schema for product_discount
type Product_Discount struct {
	Product_Discount_ID          uint      `json:"product_discount_id"`
	Product_Discount_Name        string    `json:"product_discount_name"`
	Product_Discount_Description string    `json:"product_discount_description"`
	Product_Discount_Percentage  float32   `json:"product_discount_percentage"`
	Product_Discount_Status      bool      `json:"product_discount_status"`
	Product_Discount_Created_At  time.Time `json:"product_discount_created_at"`
	// Product_Discount_Updated_At  time.Time `json:"product_discount_updated_at"`
	// Product_Discount_Deleted_At  time.Time `json:"product_discount_deleted_at"`
}
