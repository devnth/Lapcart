package repo

import (
	"database/sql"
	"lapcart/model"
	"lapcart/utils"
)

type AdminRepository interface {
	FindAdminByEmail(email string) (model.AdminResponse, error)
	AddDiscount(discount model.Discount) (int, error)
	AddDiscountToProduct(discount model.Discount) error
	AddCoupon(coupon model.Coupon) error
	FindDiscountByName(name string) (uint, error)
	GetAllOrders(pagenation utils.Filter) ([]model.GetOrders, utils.Metadata, error)
	ManageOrders(data model.ManageOrder) error
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

func (c *adminRepo) AddCoupon(coupon model.Coupon) error {

	query := `
				INSERT INTO
				   coupons (name, code, description, minimum_amount, value, created_at, valid_till) 
				VALUES
				   (
				      $1, $2, $3, $4, $5, $6, $7
				   )
				   ON CONFLICT(name) DO NOTHING;`

	err := c.db.QueryRow(
		query,
		coupon.Name,
		coupon.Code,
		coupon.Description,
		coupon.Min_Amount,
		coupon.Value,
		coupon.Created_At,
		coupon.Valid_Till,
	).Err()
	return err
}

func (c *adminRepo) FindDiscountByName(name string) (uint, error) {

	var discount_id uint

	query := `
				SELECT
					id 
				 FROM
					discount
				 WHERE
					name = $1
				 ;`

	err := c.db.QueryRow(query, name).Scan(&discount_id)

	return discount_id, err
}

func (c *adminRepo) ManageOrders(data model.ManageOrder) error {

	query := `
				UPDATE
				   order_details 
				SET
				   status = $1 , updated_at = $2
				WHERE
				   id = $3;`

	err := c.db.QueryRow(query,
		data.Status,
		data.Updated_At,
		data.Order_ID).Err()

	return err
}

func (c *adminRepo) GetAllOrders(pagenation utils.Filter) ([]model.GetOrders, utils.Metadata, error) {

	query :=

		`WITH orders AS 
		(
		   SELECT
			  * 
		   FROM
			  order_details
		)
		SELECT
		   COUNT(*) OVER(),
		   u.id,
		   CONCAT(u.first_name, ' ', u.last_name) AS fullname,
		   u.email,
		   o.id,
		   o.total,
		   o.is_paid,
		   o.status,
		   o.created_at,
		   o.updated_at 
		FROM
		   users u 
		   JOIN
			  orders o 
			  ON u.id = o.user_id
		LIMIT $1 OFFSET $2;`

	rows, err := c.db.Query(query, pagenation.Limit(), pagenation.Offset())

	if err != nil {
		return nil, utils.Metadata{}, err
	}

	defer rows.Close()

	var orders []model.GetOrders
	var totalRecords int

	for rows.Next() {
		var order model.GetOrders

		err = rows.Scan(
			&totalRecords,
			&order.User_ID,
			&order.User_Name,
			&order.Email,
			&order.OrderID,
			&order.TotalAmount,
			&order.Is_Paid,
			&order.Status,
			&order.Ordered_At,
			&order.Updated_At)

		if err != nil {
			return orders, utils.ComputeMetaData(totalRecords, pagenation.Page, pagenation.PageSize), err
		}

		orders = append(orders, order)

	}

	if err := rows.Err(); err != nil {
		return orders, utils.ComputeMetaData(totalRecords, pagenation.Page, pagenation.PageSize), err
	}

	return orders, utils.ComputeMetaData(totalRecords, pagenation.Page, pagenation.PageSize), nil

}
