package service

import (
	"crypto/md5"
	"errors"
	"fmt"
	"lapcart/repo"
)

type AuthService interface {
	VerifyAdminCredential(email, password string) error
	VerifyUserCredential(email, password string) error
}

type authService struct {
	adminRepo repo.AdminRepository
	userRepo  repo.UserRepository
}

func NewAuthService(
	adminRepo repo.AdminRepository,
	userRepo repo.UserRepository,
) AuthService {
	return &authService{
		adminRepo: adminRepo,
		userRepo:  userRepo,
	}
}

func (c *authService) VerifyAdminCredential(email, password string) error {

	admin, err := c.adminRepo.FindAdminByEmail(email)

	if err != nil {
		return errors.New("failed to login. check your email")
	}

	isValidPassword := VerifyPassword(password, admin.Password)
	if !isValidPassword {
		return errors.New("failed to login. check your credential")
	}

	return nil
}

func (c *authService) VerifyUserCredential(email, password string) error {

	user, err := c.userRepo.FindUserByEmail(email)

	if err != nil {
		return errors.New("failed to login. check your email")
	}

	if !user.IsActive {
		return errors.New("failed to login. your account is blocked")
	}

	isValidPassword := VerifyPassword(password, user.Password)
	if !isValidPassword {
		return errors.New("failed to login. check your credential")
	}

	return nil
}

func VerifyPassword(requestPassword, dbPassword string) bool {

	requestPassword = fmt.Sprintf("%x", md5.Sum([]byte(requestPassword)))
	return requestPassword == dbPassword
}
