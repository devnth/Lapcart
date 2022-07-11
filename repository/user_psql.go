package repository

import (
	"database/sql"

	"devnth/models"
	"log"
)

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

//function to add user details to user database
func (r Repository) UserSignup(user models.User) models.User {

	query := `INSERT INTO users (
		first_name, 
		last_name, 
		password, 
		email, 
		phone_number) 
		VALUES($1, $2, $3, $4, $5) 
		RETURNING 
		user_id, 
		is_active;`

	//Makes query
	err := r.DB.QueryRow(query,
		user.First_Name,
		user.Last_Name,
		user.Password,
		user.Email,
		user.Phone_Number).Scan(
		&user.User_ID,
		&user.IsActive)
	logFatal(err)
	return user
}

// To check whether r.User exists
func (r Repository) DoesUserExists(user models.User) bool {

	query := "SELECT email FROM users WHERE email = $1"

	err := r.DB.QueryRow(query,
		user.Email).Scan(
		&user.User_ID)
	// returns false if the user email does not exists
	return (err != sql.ErrNoRows)
}

func (r Repository) UserLogin(user models.User) (models.User, error) {

	query := "SELECT * FROM users WHERE email = $1"

	err := r.DB.QueryRow(query,
		user.Email).Scan(
		&user.User_ID,
		&user.First_Name,
		&user.Last_Name,
		&user.Password,
		&user.Email,
		&user.Phone_Number,
		&user.IsActive,
		&user.CreatedAt)

	return user, err
}

// // To check whether r.User exists
// func (r Repository) DoesAdminExists(admin models.Admin) bool {

// 	query := "SELECT email FROM admin WHERE email = $1"

// 	err := r.DB.QueryRow(query, admin.Email).Scan(&admin.Admin_ID)
// 	// returns false if the user email does not exists
// 	return (err != sql.ErrNoRows)
// }

func (r Repository) AdminLogin(admin models.Admin) (models.Admin, error) {

	query := "SELECT * FROM admin WHERE email = $1"

	err := r.DB.QueryRow(query,
		admin.Email).Scan(
		&admin.Admin_ID,
		&admin.First_Name,
		&admin.Last_Name,
		&admin.Password,
		&admin.Email,
		&admin.Phone_Number,
		&admin.CreatedAt)

	return admin, err
}

func (r Repository) AdminViewUsers() ([]models.User, error) {

	var users []models.User

	query := `SELECT user_id, 
			 first_name, 
			 last_name, 
			 email, 
			 phone_number, 
			 is_active,
			 created_at FROM users;`

	rows, err := r.DB.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	//Loop through each rows
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.User_ID,
			&user.First_Name,
			&user.Last_Name,
			&user.Email,
			&user.Phone_Number,
			&user.IsActive,
			&user.CreatedAt); err != nil {
			return users, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return users, err
	}

	return users, nil
}

//function to block or unblock users
func (r Repository) BlockUsers(user models.User) (models.User, error) {
	//Query
	query := `UPDATE users 
			SET is_active = $1 
			WHERE email = $2 
			RETURNING user_id, 
			first_name, 
			last_name, 
			email, 
			phone_number, 
			is_active, 
			created_at;`

	err := r.DB.QueryRow(query,
		user.IsActive,
		user.Email).Scan(
		&user.User_ID,
		&user.First_Name,
		&user.Last_Name,
		&user.Email,
		&user.Phone_Number,
		&user.IsActive,
		&user.CreatedAt)

	return user, err

}

func (r Repository) AddToCart(cart models.Cart) (models.Cart, error) {
	//writing query
	query := `INSERT INTO cart(
		user_id,
		product_id,
		product_count)
		VALUES ($1, $2, $3)
		RETURNING
		cart_id,
		user_id,
		product_id,
		product_count,
		cart_created_at;`

	err := r.DB.QueryRow(
		query,
		cart.User_ID,
		cart.Product_ID.Product_ID,
		cart.Product_Count,
	).Scan(
		&cart.Cart_ID,
		&cart.User_ID,
		&cart.Product_ID.Product_ID,
		&cart.Product_Count,
		&cart.Cart_Created_At,
	)

	return cart, err
}

func (r Repository) ViewCart(cart models.Cart) ([]models.Cart, error) {

	query := `SELECT 
		cart.cart_id, 
		cart.user_id, 
		product.product_id,
		product.product_name,
		cart.product_count
		FROM product
		INNER JOIN cart ON product.product_id=  cart.product_id
		WHERE user_id = $1;`

	rows, err := r.DB.Query(query,
		cart.User_ID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var carts []models.Cart

	for rows.Next() {
		err := rows.Scan(
			&cart.Cart_ID,
			&cart.User_ID,
			&cart.Product_ID.Product_ID,
			&cart.Product_ID.Product_Name,
			&cart.Product_Count,
		)
		if err != nil {
			return carts, err
		}

		// cart.Product_ID, err = r.ViewEachProduct(int(cart.Product_ID.Product_ID))

		// if err != nil {
		// 	log.Println(err)
		// 	return carts, err
		// }

		carts = append(carts, cart)
	}
	if err = rows.Err(); err != nil {
		return carts, err
	}

	return carts, nil

}
