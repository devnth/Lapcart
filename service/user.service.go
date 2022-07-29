package service

import (
	"crypto/md5"
	"database/sql"
	"errors"
	"fmt"
	"lapcart/model"
	"lapcart/repo"
	"time"
)

type UserService interface {
	FindUserByEmail(email string) (*model.UserResponse, error)
	CreateUser(registerRequest model.User) error
	AddAddress(address model.Address) error
	GetAddressByUserID(user_id int) (*[]model.AddressResponse, error)
	DeleteAddress(user_id, address_id int) error
}

type userService struct {
	userRepo repo.UserRepository
}

func NewUserService(userRepo repo.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
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

func HashPassword(password string) string {
	data := []byte(password)
	password = fmt.Sprintf("%x", md5.Sum(data))
	return password

}
