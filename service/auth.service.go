package service

import (
	"crypto/md5"
	"errors"
	"fmt"
	"lapcart/repo"
)

type AuthService interface {
	VerifyAdminCredential(email, password string) error
}

type authService struct {
	adminRepo repo.AdminRepository
}

func NewAuthService(adminRepo repo.AdminRepository) AuthService {
	return &authService{
		adminRepo: adminRepo,
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

func HashPassword(password string) string {
	data := []byte(password)
	password = fmt.Sprintf("%x", md5.Sum(data))
	return password

}

func VerifyPassword(requestPassword, dbPassword string) bool {

	requestPassword = fmt.Sprintf("%x", md5.Sum([]byte(requestPassword)))
	return requestPassword == dbPassword
}
