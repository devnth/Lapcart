package repo

import (
	"database/sql"
	"lapcart/model"
)

type AdminRepository interface {
	FindAdminByEmail(email string) (model.AdminResponse, error)
	// ViewUsers([]model.User, error)
	// BlockUser(int) (int, error)
	// ViewUser(int) (model.User, error)
}

type adminRepo struct {
	db *sql.DB
}

func NewAdminRepo(db *sql.DB) AdminRepository {
	return &adminRepo{
		db: db,
	}
}

func (c *adminRepo) FindAdminByEmail(email string) (model.AdminResponse, error) {

	var admin model.AdminResponse

	query := `SELECT 
			user_id,
			first_name,
			last_name,
			password,
			email,
			phone_number
			FROM admin WHERE email = $1;`

	err := c.db.QueryRow(query,
		email).Scan(
		&admin.ID,
		&admin.First_Name,
		&admin.Last_Name,
		&admin.Password,
		&admin.Email,
		&admin.Phone_Number)
	return admin, err
}

// func (c *adminRepo) ViewUsers([]model.User, error)

// func (c *adminRepo) BlockUser(int) (int, error)

// func (c *adminRepo) ViewUser(int) (model.User, error)
