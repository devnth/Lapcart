package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"lapcart/model"
	"log"
)

type UserRepository interface {
	FindUserByEmail(email string) (model.UserResponse, error)
	InsertUser(user model.User) (int, error)
	AllUsers() ([]model.UserResponse, error)
	ManageUsers(email string, isActive bool) error
	AddAddress(address model.Address) (int, error)
	GetAddressByUserID(user_id int) ([]model.AddressResponse, error)
	DeleteAddressById(user_id, address_id int) error
	FindAddressByUserID(user_id int) (uint, error)
	AddOrder(orderDetails model.OrderDetails) (uint, error)
	AddOrderItems(orderItems model.OrderItems) error
	FindOrderByUserID(user_id int) (uint, float64, error)
	VerifyCoupon(code string) (uint, float64, float64, error)
	Payment(data model.Payment) error
	UpdateUser(data model.User) error
	GetUserByID(id int) (model.UserResponse, error)
	CreateVerifyData(user_id, code int) error
	GetCodeByUserID(user_id int) (int, error)
	DeleteVerifyData(user_id int) error
	UpdateOfferPrice(data model.Payment) error
	GetAllOrders(user_id int) ([]model.Orders, error)
	CancelOrderById(order_id, user_id int) error
	FindProductAndCountIdFromOrder(order_id int) ([]int, []int, error)
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) UserRepository {
	return &userRepo{
		db: db,
	}
}

func (c *userRepo) FindUserByEmail(email string) (model.UserResponse, error) {

	var user model.UserResponse

	query := `SELECT 
				id,
				first_name,
				last_name,
				password,
				email,
				phone_number,
				is_active
				FROM users 
				WHERE email = $1;`

	err := c.db.QueryRow(query,
		email).Scan(
		&user.ID,
		&user.First_Name,
		&user.Last_Name,
		&user.Password,
		&user.Email,
		&user.Phone_Number,
		&user.IsActive,
	)

	return user, err
}

func (c *userRepo) InsertUser(user model.User) (int, error) {

	var user_id int

	query := `INSERT INTO users(
			first_name,
			last_name,
			email,
			password,
			phone_number)
			VALUES
			($1, $2, $3, $4, $5)
			RETURNING id;`

	err := c.db.QueryRow(query,
		user.First_Name,
		user.Last_Name,
		user.Email,
		user.Password,
		user.Phone_Number).Scan(
		&user_id,
	)
	return user_id, err
}

func (c *userRepo) AllUsers() ([]model.UserResponse, error) {

	var users []model.UserResponse

	query := `SELECT 
			 id, 
			 first_name, 
			 last_name, 
			 email, 
			 phone_number, 
			 is_active,
			 created_at
			 FROM users;`

	rows, err := c.db.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	//Loop through each rows
	for rows.Next() {
		var user model.UserResponse
		err := rows.Scan(&user.ID,
			&user.First_Name,
			&user.Last_Name,
			&user.Email,
			&user.Phone_Number,
			&user.IsActive,
			&user.Created_At)

		if err != nil {
			return users, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return users, err
	}

	return users, nil
}

func (c *userRepo) ManageUsers(email string, isActive bool) error {
	//Query
	query := `UPDATE users 
			SET is_active = $1 
			WHERE email = $2 ;`

	err := c.db.QueryRow(query,
		isActive,
		email).Err()

	return err
}

func (c *userRepo) AddAddress(address model.Address) (int, error) {

	var id int

	query := `INSERT INTO address(
				type,
				user_id,
				house_name,
				street_name,
				landmark,
				district,
				state,
				country,
				pincode,
				created_at)
				VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)	
				RETURNING id;`

	err := c.db.QueryRow(query,
		address.AddressType,
		address.User_id,
		address.HouseName,
		address.StreetName,
		address.Landmark,
		address.District,
		address.State,
		address.Country,
		address.PinCode,
		address.Created_At).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, err
}

func (c *userRepo) GetAddressByUserID(user_id int) ([]model.AddressResponse, error) {

	//writing query
	query := `SELECT 
				id,
				type,
				house_name,
				street_name,
				landmark,
				district,
				state,
				country,
				pincode
				FROM address
				WHERE user_id = $1;`

	rows, err := c.db.Query(query, user_id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var addresses []model.AddressResponse

	for rows.Next() {
		var address model.AddressResponse
		err := rows.Scan(
			&address.Id,
			&address.AddressType,
			&address.HouseName,
			&address.StreetName,
			&address.Landmark,
			&address.District,
			&address.State,
			&address.Country,
			&address.PinCode)
		if err != nil {
			return addresses, err
		}

		addresses = append(addresses, address)
	}

	if err := rows.Err(); err != nil {
		return addresses, err
	}

	return addresses, nil
}

func (c *userRepo) DeleteAddressById(user_id, address_id int) error {

	deletequery := `
				DELETE FROM address
				WHERE id = $1; `

	checkquery := `SELECT id FROM address WHERE user_id = $1 AND id = $2;`

	ctx := context.Background()

	tx, err := c.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	err = tx.QueryRow(checkquery, user_id, address_id).Scan(&address_id)

	if err != nil {
		return err
	}

	err = tx.QueryRow(deletequery, address_id).Err()
	if err != nil {
		fmt.Println("from delete:", err.Error())
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (c *userRepo) FindAddressByUserID(user_id int) (uint, error) {

	var address_id uint

	query := `
				SELECT
					id 
				 FROM
					address 
				 WHERE
					user_id = $1;`

	err := c.db.QueryRow(query, user_id).Scan(&address_id)

	return address_id, err

}

func (c *userRepo) AddOrder(orderDetails model.OrderDetails) (uint, error) {

	query := `
				INSERT INTO
				   order_details( user_id, shipping_address_id, total, created_at, updated_at) 
				VALUES
				   (
				      $1, $2, $3, $4, $5 
				   )
				RETURNING id
				   ;`

	err := c.db.QueryRow(query,
		orderDetails.User_ID,
		orderDetails.Address_ID,
		orderDetails.TotalPrice,
		orderDetails.Created_At,
		orderDetails.Updated_At).
		Scan(&orderDetails.ID)

	return orderDetails.ID, err
}

func (c *userRepo) AddOrderItems(orderItems model.OrderItems) error {

	var query string

	if orderItems.DiscountID != 0 {

		query = `
				INSERT INTO
				   order_items ( order_id, product_id, quantity, created_at,discount_id) 
				VALUES
				   (
				      $1, $2, $3, $4, $5 
				   )
				;`

		err := c.db.QueryRow(query,
			orderItems.OrderID,
			orderItems.ProductID,
			orderItems.Quantity,
			orderItems.Created_At,
			orderItems.DiscountID).Err()

		return err
	}

	query = `
				INSERT INTO
				   order_items ( order_id, product_id, quantity, created_at) 
				VALUES
				   (
				      $1, $2, $3, $4
				   )
				;`

	err := c.db.QueryRow(query,
		orderItems.OrderID,
		orderItems.ProductID,
		orderItems.Quantity,
		orderItems.Created_At).Err()

	return err

}

func (c *userRepo) FindOrderByUserID(user_id int) (uint, float64, error) {

	var order_id uint
	var total float64

	query := `
				SELECT
				   id,
				   total
				FROM
				   order_details 
				WHERE
				   user_id = $1 
				   AND is_paid = false AND status <> 'cancelled';`

	err := c.db.QueryRow(query, user_id).Scan(&order_id, &total)

	return order_id, total, err
}

func (c *userRepo) VerifyCoupon(code string) (uint, float64, float64, error) {

	var couponValue, minAmount float64
	var id uint

	query := `
				SELECT
					id,
					minimum_amount,
					value 
				 FROM
					coupons 
				 WHERE
					code = $1 
					AND valid_till > NOW();`

	err := c.db.QueryRow(
		query,
		code).Scan(
		&id,
		&minAmount,
		&couponValue,
	)

	return id, minAmount, couponValue, err
}

func (c *userRepo) UpdateOfferPrice(data model.Payment) error {

	query := `
				UPDATE
					order_details
				 SET
					coupon_id = $1 ,total = $2 ,updated_at = $3`

	err := c.db.QueryRow(query, data.Coupon_Id, data.Amount, data.Updated_At).Err()

	return err
}

func (c *userRepo) Payment(data model.Payment) error {

	var arg []interface{}

	insertQuery := `
				INSERT INTO
				   payment ( payment_type, created_at, updated_at, payment_id, razor_order_id, payment_signature )
				VALUES
				   (
				      $1, $2, $3, $4, $5, $6
				   )
				RETURNING id;`

	updateQuery := `
					UPDATE
					   order_details
					SET
					   is_paid = true, payment_id = $1 ,updated_at = $2
					
	`

	err := c.db.QueryRow(
		insertQuery,
		data.PaymentType,
		data.Created_At,
		data.Updated_At,
		data.Razorpay_payment_id,
		data.Razorpay_order_id,
		data.Razorpay_signature).Scan(&data.ID)

	if err != nil {
		return err
	}

	arg = append(arg, data.ID)
	arg = append(arg, data.Updated_At)

	i := 3

	if data.Coupon_Id != 0 {
		updateQuery = updateQuery + `,total = $3, coupon_id = $4`
		i = 5
		arg = append(arg, data.Amount)
		arg = append(arg, data.Coupon_Id)
	}

	updateQuery = updateQuery + `
								WHERE
								id = $` + fmt.Sprintf(`%d;`, i)

	arg = append(arg, data.Order_ID)

	stmt, err := c.db.Prepare(updateQuery)

	if err != nil {
		log.Println("Error", err)
		log.Println("Error", "Query prepare failed")
		return err
	}

	_, err = stmt.Query(arg...)

	if err != nil {
		log.Println("Error", err)
		log.Println("Error", "Query Exec failed")
		return err
	}

	if err != nil {
		return err
	}

	return nil
}

func (c *userRepo) UpdateUser(data model.User) error {

	query := `
				UPDATE
					users 
				 SET `
	i := 1
	var arg []interface{}

	if data.First_Name != "" {
		query = query + `first_name = $` + fmt.Sprintf(`%d`, i)
		i++
		arg = append(arg, data.First_Name)
	}

	if data.Last_Name != "" {
		if i > 1 {
			query = query + `, `
		}
		query = query + `last_name = $` + fmt.Sprintf(`%d`, i)
		i++
		arg = append(arg, data.Last_Name)
	}

	if data.Email != "" {
		if i > 1 {
			query = query + `, `
		}
		query = query + `email = $` + fmt.Sprintf(`%d`, i)
		arg = append(arg, data.Email)
		i++
	}

	if data.Password != "" {
		if i > 1 {
			query = query + `, `
		}
		query = query + `password = $` + fmt.Sprintf(`%d`, i)
		arg = append(arg, data.Password)
		i++
	}

	if data.Phone_Number != 0 {
		if i > 1 {
			query = query + `, `
		}
		query = query + `phone_number = $` + fmt.Sprintf(`%d`, i)
		arg = append(arg, data.Phone_Number)
		i++
	}

	if data.IsVerified {
		if i > 1 {
			query = query + `, `
		}
		query = query + `is_verified = $` + fmt.Sprintf(`%d`, i)
		arg = append(arg, data.IsVerified)
		i++
	}

	if i > 1 {
		query = query + `, `
	}
	query = query + `updated_at = $` + fmt.Sprintf(`%d 
											WHERE id = $%d;`, i, i+1)
	arg = append(arg, data.Updated_At)
	arg = append(arg, data.ID)

	statement, err := c.db.Prepare(query)

	if err != nil {
		log.Println("Error ", "error in preparing query: ", err)
		return err
	}

	_, err = statement.Query(arg...)

	if err != nil {
		log.Println("Error ", "error in query execution: ", err)
		return err
	}

	return nil
}

func (c *userRepo) GetUserByID(id int) (model.UserResponse, error) {

	var user model.UserResponse

	query := `SELECT 
				CONCAT (first_name,' ',last_name) AS fullname,
				email,
				phone_number,
				is_verified
			  FROM users
			  WHERE id = $1;`

	err := c.db.QueryRow(query, id).Scan(
		&user.Full_Name,
		&user.Email,
		&user.Phone_Number,
		&user.IsVerified)

	return user, err
}

func (c *userRepo) CreateVerifyData(user_id, code int) error {

	query := `
				INSERT INTO
					verify_email ( user_id, code) 
				 VALUES
					(
					   $1, $2
					)
				 ;`

	err := c.db.QueryRow(query, user_id, code).Err()

	return err
}

func (c *userRepo) GetCodeByUserID(user_id int) (int, error) {

	var code int

	query := `
			SELECT 
				code 
			FROM verify_email
			WHERE user_id = $1;	`

	err := c.db.QueryRow(query, user_id).Scan(&code)

	return code, err
}

func (c *userRepo) DeleteVerifyData(user_id int) error {

	query := `
					DELETE FROM verify_email
					WHERE user_id = $1;`

	err := c.db.QueryRow(query, user_id).Err()

	return err
}

func (c *userRepo) GetAllOrders(user_id int) ([]model.Orders, error) {

	query := `
						WITH od AS 
						(
						   SELECT
						      *,
						      created_at AS ordered_at 
						   FROM
						      order_details 
						   WHERE
						      user_id = 8
						)
						,
						items AS 
						(
						   SELECT
						      * 
						   FROM
						      order_items oi 
						      JOIN
						         od od 
						         ON oi.order_id = od.id 
						   WHERE
						      oi.order_id = od.id
						)
						,
						products AS 
						(
						   SELECT
						      * 
						   FROM
						      PRODUCT p 
						      JOIN
						         items it 
						         ON p.id = it.product_id 
						)
						SELECT
						   p.order_id,
						   p.product_id,
						   p.image,
						   p.name,
						   p.color,
						   b.name,
						   c.name,
						   proc.name,
						   p.status,
						   p.ordered_at
						FROM
						   products p 
						   JOIN
						      category c 
						      ON p.category_id = c.id 
						   JOIN
						      brand b 
						      ON p.brand_id = b.id 
						   JOIN
						      processor proc 
						      ON p.processor_id = proc.id;`

	rows, err := c.db.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var orders []model.Orders

	for rows.Next() {

		var order model.Orders

		err = rows.Scan(
			&order.OrderID,
			&order.ProductID,
			&order.Image,
			&order.Name,
			&order.Color,
			&order.Brand,
			&order.Category,
			&order.Processor,
			&order.Status,
			&order.Ordered_At,
		)

		if err != nil {
			return orders, err
		}

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return orders, err
	}

	return orders, nil

}

func (c *userRepo) CancelOrderById(order_id, user_id int) error {

	querycheck :=

		`SELECT 
		TRUE 
		FROM order_details
		WHERE id = $1 AND user_id = $2 AND status = 'waiting';`

	var ok bool

	c.db.QueryRow(querycheck, order_id, user_id).Scan(&ok)

	if !ok {
		return errors.New("you have no erros to cancel")
	}

	query := `UPDATE
				 order_details 
			     SET status = 'cancelled'
			  WHERE
			     id = $1 AND user_id = $2 AND status <> 'cancelled';`

	_, err := c.db.Exec(query, order_id, user_id)

	return err
}

func (c *userRepo) FindProductAndCountIdFromOrder(order_id int) ([]int, []int, error) {

	query := `SELECT 
				product_id, quantity
			 FROM order_items
			 WHERE order_id = $1;`

	rows, err := c.db.Query(query, order_id)

	if err != nil {
		return nil, nil, err
	}

	defer rows.Close()

	var product_ids, quantities []int

	for rows.Next() {

		var product_id, quantity int

		if err := rows.Scan(&product_id, &quantity); err != nil {
			return product_ids, quantities, err
		}

		product_ids = append(product_ids, product_id)
		quantities = append(quantities, quantity)

	}

	if err := rows.Err(); err != nil {
		return product_ids, quantities, err
	}

	return product_ids, quantities, nil
}
