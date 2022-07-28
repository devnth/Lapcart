package service

import (
	"errors"
	"lapcart/model"
	"lapcart/repo"
)

type AdminService interface {
	FindAdminByEmail(email string) (*model.AdminResponse, error)
	AllUsers() (*[]model.UserResponse, error)
	ManageUsers(email string, isActive bool) error
}

type adminService struct {
	adminRepo repo.AdminRepository
	userRepo  repo.UserRepository
}

func NewAdminService(
	adminRepo repo.AdminRepository,
	userRepo repo.UserRepository) AdminService {
	return &adminService{
		adminRepo: adminRepo,
		userRepo:  userRepo,
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
