package repo

import (
	"database/sql"
	"fmt"
	"lapcart/model"
	"lapcart/utils"
	"log"
	"time"
)

type AdminRepository interface {
	FindAdminByEmail(email string) (model.AdminResponse, error)
	AddDiscount(discount model.Discount) (int, error)
	AddDiscountToProduct(discount model.Discount) error
	AddCoupon(coupon model.Coupon) error
	FindDiscountByName(name string) (uint, error)
	GetAllOrders(pagenation utils.Filter) ([]model.GetOrders, utils.Metadata, error)
	ManageOrders(data model.ManageOrder) error
	AddCategory(data model.Category) error
	GetAllCategory() ([]model.Category, error)
	UpdateCategory(data model.Category) error
	DeleteCategory(id uint) error
	DeleteProductByCategory(id uint) error
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

func (c *adminRepo) AddCategory(data model.Category) error {

	query := `INSERT INTO 
				category 
					   (name, description)
   			VALUES 
						($1, $2);`

	err := c.db.QueryRow(query, data.Name, data.Description).Err()

	return err
}

func (c *adminRepo) GetAllCategory() ([]model.Category, error) {

	var categories []model.Category

	query := ` SELECT * FROM category;`

	rows, err := c.db.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	//Loop through each rows
	for rows.Next() {
		var category model.Category
		err := rows.Scan(&category.ID,
			&category.Name,
			&category.Description)

		if err != nil {
			return categories, err
		}
		categories = append(categories, category)
	}
	if err = rows.Err(); err != nil {
		return categories, err
	}

	return categories, nil

}

func (c *adminRepo) UpdateCategory(data model.Category) error {

	query := `UPDATE 
				category 
				SET
				`

	var arg []interface{}
	i := 1

	if data.Name != "" {

		query = query + fmt.Sprintf(`name = $%d`, i)
		arg = append(arg, data.Name)
		i++

	}

	if data.Description != "" {

		if i > 1 {
			query = query + `, `
		}

		query = query + fmt.Sprintf(`description = $%d`, i)
		arg = append(arg, data.Description)
		i++
	}

	statement, err := c.db.Prepare(query)

	if err != nil {
		log.Println("Error", "query exec failed", err)
		return err
	}

	err = statement.QueryRow(arg...).Err()

	if err != nil {
		log.Println("Error", "query exec failed: ", err)
	}

	return err
}

func (c *adminRepo) DeleteCategory(id uint) error {

	query := `Update 
				 category 
				SET is_deleted = true
				WHERE id = $1;`

	err := c.db.QueryRow(query, id).Err()

	return err

}

func (c *adminRepo) DeleteProductByCategory(id uint) error {

	deleted_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	query := `Update 
				 product 
				SET is_deleted = true, deleted_at = $1
				WHERE category_id = $2;
				`

	err := c.db.QueryRow(query, deleted_at, id).Err()

	return err

}
