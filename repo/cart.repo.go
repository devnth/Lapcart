package repo

import (
	"context"
	"database/sql"
	"lapcart/model"
	"log"
)

type CartRepository interface {
	AddToCart(cart model.Cart) (string, error)
}

type cartRepo struct {
	db *sql.DB
}

func NewCartRepository(db *sql.DB) CartRepository {
	return &cartRepo{
		db: db,
	}

}

func (c *cartRepo) AddToCart(cart model.Cart) (string, error) {

	ctx := context.Background()
	var message string
	var id int

	checkQuery := `SELECT
				  	id
				   FROM
				  	cart 
				   WHERE
				  	user_id = $1 
				   AND product_id = $2;`

	insertQuery := `INSERT INTO
					   cart ( user_id, product_id, count, created_at, updated_at) 
					VALUES
					   (
					      $1, $2, $3, $4, $5
					   );`

	updateQuery := `UPDATE
				  	cart 
				   SET
				  	count = 
				  	(
				  	   count + $1
				  	)
				   WHERE
				  	id = $2 ;`

	tx, err := c.db.BeginTx(ctx, nil)

	if err != nil {
		return "", err
	}

	err = tx.QueryRow(
		checkQuery,
		cart.User_Id,
		cart.Product_Id).Scan(
		&id)

	if err != sql.ErrNoRows && err != nil {
		return "", err
	}

	if err == sql.ErrNoRows {

		err = tx.QueryRow(
			insertQuery,
			cart.User_Id,
			cart.Product_Id,
			cart.Count,
			cart.Created_At,
			cart.Updated_At,
		).Err()

		if err != nil {
			log.Println("error in inserting cart")
			tx.Rollback()
			return "", err
		}
		message = "added to cart"
	} else {

		err = tx.QueryRow(
			updateQuery,
			cart.Count,
			id,
		).Scan(&cart.Count)

		if err != nil {
			log.Println("error in updating cart")
			tx.Rollback()
			return "", err
		}

		message = "updated cart"
	}

	err = tx.Commit()
	if err != nil {
		log.Println("error in commiting in cart")
		return "", err
	}

	return message, nil

}
