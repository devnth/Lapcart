package repo

import (
	"database/sql"
	"lapcart/model"
)

type UserRepository interface {
	FindUserByEmail(email string) (model.UserResponse, error)
	InsertUser(user model.User) (int, error)
	AllUsers() ([]model.UserResponse, error)
	ManageUsers(email string, isActive bool) error
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
