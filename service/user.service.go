package service

import (
	"crypto/md5"
	"database/sql"
	"errors"
	"fmt"
	"lapcart/model"
	"lapcart/repo"
)

type UserService interface {
	FindUserByEmail(email string) (*model.UserResponse, error)
	CreateUser(registerRequest model.User) error
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

func HashPassword(password string) string {
	data := []byte(password)
	password = fmt.Sprintf("%x", md5.Sum(data))
	return password

}
