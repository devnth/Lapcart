package service

import (
	"lapcart/model"
	"lapcart/repo"
)

type AdminService interface {
	FindAdminByEmail(email string) (*model.AdminResponse, error)
}

type adminService struct {
	adminRepo repo.AdminRepository
}

func NewAdminService(adminRepo repo.AdminRepository) AdminService {
	return &adminService{
		adminRepo: adminRepo,
	}
}

func (c *adminService) FindAdminByEmail(email string) (*model.AdminResponse, error) {
	admin, err := c.adminRepo.FindAdminByEmail(email)

	if err != nil {
		return nil, err
	}

	return &admin, nil
}
