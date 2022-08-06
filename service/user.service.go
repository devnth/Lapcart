package service

import (
	"crypto/md5"
	"database/sql"
	"errors"
	"fmt"
	"lapcart/model"
	"lapcart/repo"
	"lapcart/utils"
	"log"
	"time"
)

type UserService interface {
	FindUserByEmail(email string) (*model.UserResponse, error)
	CreateUser(registerRequest model.User) error
	AddAddress(address model.Address) error
	GetAddressByUserID(user_id int) (*[]model.AddressResponse, error)
	DeleteAddress(user_id, address_id int) error
	GetAllProducts(filter model.Filter, user_id int, pagenation utils.Filter) (*[]model.GetProduct, *utils.Metadata, error)
	ProceedToCheckout(user_id int) error
	Payment(data model.Payment) error
}

type userService struct {
	userRepo    repo.UserRepository
	productRepo repo.ProductRepository
	cartRepo    repo.CartRepository
	adminRepo   repo.AdminRepository
}

func NewUserService(
	userRepo repo.UserRepository,
	productRepo repo.ProductRepository,
	cartRepo repo.CartRepository,
	adminRepo repo.AdminRepository) UserService {
	return &userService{
		userRepo:    userRepo,
		productRepo: productRepo,
		cartRepo:    cartRepo,
		adminRepo:   adminRepo,
	}
}

func (c *userService) FindUserByEmail(email string) (*model.UserResponse, error) {
	user, err := c.userRepo.FindUserByEmail(email)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (c *userService) CreateUser(registerRequest model.User) error {

	_, err := c.userRepo.FindUserByEmail(registerRequest.Email)

	if err == nil {
		return errors.New("user already Exists")
	}

	if err != nil && err != sql.ErrNoRows {
		return err
	}

	//hashing password
	registerRequest.Password = HashPassword(registerRequest.Password)

	_, err = c.userRepo.InsertUser(registerRequest)
	if err != nil {
		return errors.New("error inserting user in the database")
	}
	return nil

}

func (c *userService) AddAddress(address model.Address) error {

	address.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	_, err := c.userRepo.AddAddress(address)

	if err != nil {
		return err
	}

	return nil
}

func (c *userService) GetAddressByUserID(user_id int) (*[]model.AddressResponse, error) {

	address, err := c.userRepo.GetAddressByUserID(user_id)

	if err != nil {
		return nil, err
	}

	if address == nil {
		return nil, errors.New("no address found")
	}

	return &address, nil
}

func (c *userService) DeleteAddress(user_id, address_id int) error {

	err := c.userRepo.DeleteAddressById(user_id, address_id)

	if err != nil {
		return errors.New("unable to delete the requested address")
	}
	return nil
}

func (c *userService) GetAllProducts(filter model.Filter, user_id int, pagenation utils.Filter) (*[]model.GetProduct, *utils.Metadata, error) {

	products, metadata, err := c.productRepo.GetAllProducts(filter, user_id, pagenation)

	if err != nil {
		return nil, &metadata, err

	}

	return &products, &metadata, nil

}

func (c *userService) ProceedToCheckout(user_id int) error {

	var orderDetails model.OrderDetails
	var err error
	var carts []model.GetCart
	var orderItems model.OrderItems
	var inCart model.Cart

	orderDetails.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	orderDetails.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	orderItems.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	orderDetails.User_ID = uint(user_id)

	orderDetails.Address_ID, err = c.userRepo.FindAddressByUserID(user_id)

	if err == sql.ErrNoRows {
		return errors.New("please add shippping address")
	}

	carts, orderDetails.TotalPrice, _ = c.cartRepo.GetCartByUserId(user_id)

	if len(carts) == 0 {
		return errors.New("cart empty")
	}

	orderItems.OrderID, err = c.userRepo.AddOrder(orderDetails)

	if err != nil {
		log.Println("unable to add to order table")
		return err
	}

	for _, cart := range carts {

		orderItems.ProductID = cart.ProductID
		orderItems.DiscountID, _ = c.adminRepo.FindDiscountByName(cart.DiscountName)
		orderItems.Quantity = cart.Count

		err := c.userRepo.AddOrderItems(orderItems)

		if err != nil {
			return err
		}

		inCart.ID = int(cart.CartID)
		inCart.User_Id = user_id

		_, err = c.cartRepo.DeleteCart(inCart)

		if err != nil {
			log.Println(err.Error())
			return errors.New("unable to delete from cart")
		}

	}
	return nil
}

func (c *userService) Payment(data model.Payment) error {

	var err error
	var minAmount, couponValue float64

	data.Order_ID, data.Amount, err = c.userRepo.FindOrderByUserID(data.User_ID)

	if err != nil {
		log.Println("error in finding order: ", err)
		return errors.New("unable to find order")
	}

	if data.Coupon_Code != "" {
		data.Coupon_Id, minAmount, couponValue, err = c.userRepo.VerifyCoupon(data.Coupon_Code)

		if err != nil {
			log.Println("error in verifying coupon: ", err)
			return errors.New("invalid coupon")
		}

		if data.Amount < minAmount {
			return errors.New("coupon not applicable for this purchase")
		}

		data.Amount = data.Amount - couponValue
	}

	err = c.userRepo.Payment(data)

	if err != nil {
		log.Println("error in payment: ", err)
		return err
	}

	return nil

}

func HashPassword(password string) string {
	data := []byte(password)
	password = fmt.Sprintf("%x", md5.Sum(data))
	return password

}
