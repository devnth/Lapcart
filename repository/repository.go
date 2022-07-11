package repository

import (
	"database/sql"
	"devnth/models"
)

type Repository struct {
	DB *sql.DB
}

type ProductRepository interface {
	Addproduct(models.Product) (models.Product, error)
	ViewProduct() ([]Prod, error)
	AddCategory(models.Product_Category) (models.Product_Category, error)
	AddBranding(models.Product_Branding) (models.Product_Branding, error)
	AddProcessor(models.Product_Processor) (models.Product_Processor, error)
	ViewBranding() ([]models.Product_Branding, error)
	ViewCategory() ([]models.Product_Category, error)
	ViewProcessor() ([]models.Product_Processor, error)
	AddColor(models.Product_Color) (models.Product_Color, error)
	ViewEachProduct(int) (models.Product, error)
}

type UserRepository interface {
	UserSignup(models.User) models.User
	DoesUserExists(models.User) bool
	// DoesAdminExists(models.Admin) bool
	UserLogin(models.User) (models.User, error)
	AdminLogin(models.Admin) (models.Admin, error)
	AdminViewUsers() ([]models.User, error)
	BlockUsers(models.User) (models.User, error)
	AddToCart(models.Cart) (models.Cart, error)
	ViewCart(models.Cart) ([]models.Cart, error)
}
