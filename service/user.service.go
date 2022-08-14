package service

import (
	"crypto/md5"
	"database/sql"
	"errors"
	"fmt"
	"lapcart/config"
	"lapcart/model"
	"lapcart/repo"
	"lapcart/utils"
	"log"
	"math/rand"
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
	ProcessingPayment(data model.Payment) (*model.Payment, error)
	AddPayment(data model.Payment) error
	SendVerificationEmail(user_id int) (*string, error)
	VerifyEmail(user model.User) error
	GetAllOrders(UserID int) (*[]model.Orders, error)
	CancelOrder(order_id, user_id int) error
}

type userService struct {
	userRepo    repo.UserRepository
	productRepo repo.ProductRepository
	cartRepo    repo.CartRepository
	adminRepo   repo.AdminRepository
	mailConfig  config.MailConfig
}

func NewUserService(
	userRepo repo.UserRepository,
	productRepo repo.ProductRepository,
	cartRepo repo.CartRepository,
	adminRepo repo.AdminRepository,
	mailConfig config.MailConfig) UserService {
	return &userService{
		userRepo:    userRepo,
		productRepo: productRepo,
		cartRepo:    cartRepo,
		adminRepo:   adminRepo,
		mailConfig:  mailConfig,
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

	user, _ := c.userRepo.GetUserByID(user_id)

	if !user.IsVerified {
		return errors.New("please verify your email first to proceed to checkout")
	}

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

func (c *userService) CancelOrder(order_id, user_id int) error {

	err := c.userRepo.CancelOrderById(order_id, user_id)

	if err != nil {
		return err
	}

	product_ids, quantities, err := c.userRepo.FindProductAndCountIdFromOrder(order_id)

	if err != nil {
		log.Println("error in find product id and count :", err)
		return err
	}

	for i, product_id := range product_ids {

		err := c.productRepo.ReUpdateStockById(product_id, quantities[i])

		if err != nil {
			log.Println("errors in reupdating stock: ", err)
			return err
		}

	}

	return nil

}

func (c *userService) ProcessingPayment(data model.Payment) (*model.Payment, error) {

	var err error
	var minAmount, couponValue float64

	data.Order_ID, data.Amount, err = c.userRepo.FindOrderByUserID(data.User_ID)

	if err != nil {
		log.Println("error in finding order: ", err)
		return nil, errors.New("unable to find order")
	}

	if data.Coupon_Code != "" {
		data.Coupon_Id, minAmount, couponValue, err = c.userRepo.VerifyCoupon(data.Coupon_Code)

		if err != nil {
			log.Println("error in verifying coupon: ", err)
			return nil, errors.New("invalid coupon")
		}

		if data.Amount < minAmount {
			return nil, errors.New("coupon not applicable for this purchase")
		}

		data.Amount = data.Amount - couponValue
		err = c.userRepo.UpdateOfferPrice(data)

		if err != nil {
			return nil, err
		}
	}

	user, _ := c.userRepo.GetUserByID(data.User_ID)

	data.Email = user.Email
	data.Full_Name = user.Full_Name
	data.Phone_Number = user.Phone_Number

	return &data, nil

}

func (c *userService) AddPayment(data model.Payment) error {

	err := c.userRepo.Payment(data)

	if err != nil {
		log.Println("error in payment: ", err)
		return err
	}

	user, _ := c.userRepo.GetUserByID(data.User_ID)

	message := fmt.Sprintf(
		"Hello, %s ..\nYour order (ORDER No: %d) has been placed.\n\nTotal Amount: %.2f has been paid.\n\nVisit Again!\nThanks and regards,\n\n Lapcart team.",
		user.Full_Name,
		data.Order_ID,
		data.Amount,
	)

	err = c.mailConfig.SendMail(user.Email, message)

	if err != nil {
		return err
	}

	return nil
}

func (c *userService) SendVerificationEmail(user_id int) (*string, error) {

	user, err := c.userRepo.GetUserByID(user_id)

	if err != nil {
		return nil, err
	}

	//to generate random code
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(100000)

	message := fmt.Sprintf(
		"Hello, %s ..\nThis is the verification code as you requested:\n\n%d\nUse this code to verify your email.\nThanks and regards,\n\n Lapcart team.",
		user.Full_Name,
		code,
	)

	// send code to email
	if err = c.mailConfig.SendMail(user.Email, message); err != nil {
		return nil, err
	}

	err = c.userRepo.CreateVerifyData(user_id, code)

	if err != nil {
		return nil, err
	}

	return &user.Email, nil
}

func (c *userService) VerifyEmail(data model.User) error {

	user, _ := c.userRepo.GetUserByID(data.ID)

	if user.IsVerified {
		return errors.New("user already verified")
	}

	code, err := c.userRepo.GetCodeByUserID(data.ID)

	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if err == sql.ErrNoRows {
		return errors.New("please send verification code request first")
	}

	if code != data.Verification_Code {
		return errors.New("invalid code, please try again later")
	}

	data.IsVerified = true

	err = c.userRepo.UpdateUser(data)

	if err != nil {
		return errors.New("error in update user data")
	}

	err = c.userRepo.DeleteVerifyData(data.ID)

	if err != nil {
		return err
	}

	return nil
}

func (c *userService) GetAllOrders(UserID int) (*[]model.Orders, error) {

	orders, err := c.userRepo.GetAllOrders(UserID)

	return &orders, err

}

func HashPassword(password string) string {
	data := []byte(password)
	password = fmt.Sprintf("%x", md5.Sum(data))
	return password

}
