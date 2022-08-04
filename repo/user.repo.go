package repo

import (
	"context"
	"database/sql"
	"fmt"
	"lapcart/model"
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

	query := `
				INSERT INTO
				   order_items ( order_id, product_id, discount_id, quantity, created_at) 
				VALUES
				   (
				      $1, $2, $3, $4, $5
				   )
				;`

	err := c.db.QueryRow(query,
		orderItems.OrderID,
		orderItems.ProductID,
		orderItems.DiscountID,
		orderItems.Quantity,
		orderItems.Created_At).Err()

	return err
}
