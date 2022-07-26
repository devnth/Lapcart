package repo

import (
	"database/sql"
	"lapcart/model"
)

type UserRepository interface {
	FindUserByEmail(email string) (model.UserResponse, error)
	InsertUser(user model.User) (int, error)
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
				phone_number
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
