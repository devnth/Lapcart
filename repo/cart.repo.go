package repo

import (
	"context"
	"database/sql"
	"errors"
	"lapcart/model"
	"log"
)

type CartRepository interface {
	AddToCart(cart model.Cart) (string, error)
	GetCartByUserId(user_id int) ([]model.GetCart, float64, error)
	DeleteCart(cart model.Cart) (model.Cart, error)
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
		).Err()

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

func (c *cartRepo) GetCartByUserId(user_id int) ([]model.GetCart, float64, error) {

	var totalPrice float64

	query := `
				WITH cart AS 
				(
				   SELECT
					  c.id,
					  c.user_id,
					  c.product_id,
					  c.count 
				   FROM
					  cart c 
				   WHERE
					  user_id = $1
				)
				,
				discount AS 
				(
				   SELECT
					  id,
					  name,
					  percentage 
				   FROM
					  discount d
				)
				,
				product AS 
				(
				   SELECT
					  c.id AS cart_id,
					  p.id AS product_id,
					  P.name,
					  p.color,
					  p.brand_id,
					  p.category_id,
					  p.processor_id,
					  c.count AS count,
					  p.price,
					  p.image,
					  COALESCE(d.name, '') AS discount_name,
					  COALESCE(d.percentage, 0)  AS percentage
				   FROM
					  product P 
					  JOIN
						 cart c 
						 ON p.id = c.product_id 
					  LEFT JOIN
						 discount d 
						 ON p.discount_id = d.id
				)
				SELECT
				   p.cart_id,
				   p.product_id,
				   p.image,
				   p.name,
				   p.color,
				   b.name,
				   cat.name,
				   proc.name,
				   p.discount_name,
				   p.count,
				   p.price AS unit_price,
				   p.price * p.count AS sub_total_price,
				   COALESCE(cast((p.price * (1 - p.percentage / 100)) AS NUMERIC(10,2)),0) AS discount_price 
				FROM
				   product p 
				   JOIN
					  category cat 
					  ON p.category_id = cat.id 
				   JOIN
					  brand b 
					  ON p.brand_id = b.id 
				   JOIN
					  processor proc 
					  ON p.processor_id = proc.id;`

	var carts []model.GetCart

	rows, err := c.db.Query(query, user_id)

	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	for rows.Next() {

		var cart model.GetCart

		err := rows.Scan(
			&cart.CartID,
			&cart.ProductID,
			&cart.Image,
			&cart.Name,
			&cart.Color,
			&cart.Brand,
			&cart.Category,
			&cart.Processor,
			&cart.DiscountName,
			&cart.Count,
			&cart.UnitPrice,
			&cart.SubTotalPrice,
			&cart.DiscountPrice,
		)

		if err != nil {
			return carts, totalPrice, err
		}
		totalPrice += cart.DiscountPrice
		cart.TotalPrice += cart.DiscountPrice
		carts = append(carts, cart)
	}

	if err := rows.Err(); err != nil {
		return carts, totalPrice, err
	}

	return carts, totalPrice, nil
}

func (c *cartRepo) DeleteCart(cart model.Cart) (model.Cart, error) {

	checkQuery := `
					SELECT
						product_id,
						count 
					 FROM
						cart 
					 WHERE
						id = $1 
						AND user_id = $2;`

	deletequery := `
						DELETE
					FROM
					   cart 
					WHERE
					   id = $1;`

	err := c.db.QueryRow(
		checkQuery,
		cart.ID,
		cart.User_Id).Scan(
		&cart.Product_Id,
		&cart.Count,
	)

	if err != nil {
		log.Println(err)
		return model.Cart{}, errors.New("could not find product in cart")
	}

	err = c.db.QueryRow(
		deletequery,
		cart.ID,
	).Err()

	if err != nil {
		log.Println(err)
		return model.Cart{}, errors.New("error deleting product from cart")
	}

	return cart, nil
}
