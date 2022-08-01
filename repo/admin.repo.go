package repo

import (
	"database/sql"
	"lapcart/model"
)

type AdminRepository interface {
	FindAdminByEmail(email string) (model.AdminResponse, error)
	AddDiscount(discount model.Discount) (int, error)
	AddDiscountToProduct(discount model.Discount) error
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
			id,
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

func (c *adminRepo) AddDiscount(discount model.Discount) (int, error) {

	query := `INSERT INTO
				DISCOUNT ( name, percentage, created_at, updated_at, valid_till) 
			  VALUES
				(
				   $1, $2, $3, $4, $5
				)
				ON CONFLICT(name) DO NOTHING;`

	query2 := `SELECT
				id 
			 FROM
				discount 
			 WHERE
				name = $1;`

	err := c.db.QueryRow(
		query,
		discount.DiscountName,
		discount.Percentage,
		discount.Created_At,
		discount.Updated_At,
		discount.Expiry_Date).Err()

	if err != nil {
		return discount.ID, err
	}

	err = c.db.QueryRow(query2,
		discount.DiscountName).Scan(&discount.ID)

	return discount.ID, err
}

func (c *adminRepo) AddDiscountToProduct(discount model.Discount) error {

	var query string

	var arg int
	var err error

	if discount.Category != "" {

		query = `
				UPDATE
				   product 
				SET
				   discount_id = $1
				WHERE
				   category_id = $2 ;`

		arg = discount.ArgId

	}

	if discount.Brand != "" {

		query = `
				UPDATE
				   product 
				SET
				   discount_id = $1
				WHERE
				   brand_id = $2;`
		arg = discount.ArgId
	}

	if discount.ProductCode != "" {
		query = `
				UPDATE
				   product 
				SET
				   discount_id = $1
				WHERE
				   code = $2;`

		arg := discount.ProductCode
		err = c.db.QueryRow(query, discount.ID, arg).Err()

		return err
	}
	err = c.db.QueryRow(query, discount.ID, arg).Err()

	return err
}
