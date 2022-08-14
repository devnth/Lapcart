package service

import (
	"database/sql"
	"errors"
	"lapcart/model"
	"lapcart/repo"
	"lapcart/utils"
	"log"
)

type AdminService interface {
	FindAdminByEmail(email string) (*model.AdminResponse, error)
	AllUsers() (*[]model.UserResponse, error)
	ManageUsers(email string, isActive bool) error
	AddDiscount(discount model.Discount) error
	AddCoupon(coupon model.Coupon) error
	GetAllOrders(pagenation utils.Filter) (*[]model.GetOrders, *utils.Metadata, error)
	ManageOrders(data model.ManageOrder) error
}

type adminService struct {
	adminRepo   repo.AdminRepository
	userRepo    repo.UserRepository
	productRepo repo.ProductRepository
}

func NewAdminService(
	adminRepo repo.AdminRepository,
	userRepo repo.UserRepository,
	productRepo repo.ProductRepository) AdminService {
	return &adminService{
		adminRepo:   adminRepo,
		userRepo:    userRepo,
		productRepo: productRepo,
	}
}

func (c *adminService) FindAdminByEmail(email string) (*model.AdminResponse, error) {
	admin, err := c.adminRepo.FindAdminByEmail(email)

	if err != nil {
		return nil, err
	}

	return &admin, nil
}

func (c *adminService) AllUsers() (*[]model.UserResponse, error) {

	users, err := c.userRepo.AllUsers()
	if err != nil {
		return nil, err
	}

	return &users, nil
}

func (c *adminService) ManageUsers(email string, isActive bool) error {

	// check if user exist
	_, err := c.userRepo.FindUserByEmail(email)
	if err != nil {
		return errors.New("entered user does not exists")
	}

	err = c.userRepo.ManageUsers(email, isActive)

	return err

}

func (c *adminService) AddDiscount(discount model.Discount) error {

	var err error
	var ok bool

	discount.ID, err = c.adminRepo.AddDiscount(discount)

	if err != nil {
		return err
	}

	if discount.Category != "" {
		discount.ArgId, ok = c.productRepo.FindCategory(discount.Category)
		if !ok {
			return errors.New("category not found")
		}

	}

	if discount.Brand != "" {
		discount.ArgId, ok = c.productRepo.FindBrand(discount.Brand)
		if !ok {
			return errors.New("brand not found")
		}
	}

	if discount.ProductCode != "" {
		err = c.productRepo.FindProductCode(discount.ProductCode)

		if err == sql.ErrNoRows {
			return errors.New("product not found")
		}
	}

	err = c.adminRepo.AddDiscountToProduct(discount)

	if err != nil {
		return err
	}

	return nil

}

func (c *adminService) AddCoupon(coupon model.Coupon) error {

	err := c.adminRepo.AddCoupon(coupon)

	if err != nil {
		log.Println(err)
		return errors.New("error in adding coupon")
	}

	return nil
}

func (c *adminService) ManageOrders(data model.ManageOrder) error {

	err := c.adminRepo.ManageOrders(data)

	if err != nil {
		return err
	}

	return err
}

func (c *adminService) GetAllOrders(pagenation utils.Filter) (*[]model.GetOrders, *utils.Metadata, error) {

	orders, metadata, err := c.adminRepo.GetAllOrders(pagenation)

	if err != nil {
		return nil, nil, err
	}

	return &orders, &metadata, nil
}
