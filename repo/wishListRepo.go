package repo

import (
	"context"
	"database/sql"
	"lapcart/model"

	"github.com/lib/pq"
)

type WishListRepository interface {
	AddOrDeleteWishList(wishList model.WishList) (string, error)
	GetWishList(user_id int) ([]model.GetProduct, error)
}

type wishListRepo struct {
	db *sql.DB
}

func NewWishListRepo(db *sql.DB) WishListRepository {
	return &wishListRepo{
		db: db,
	}
}

func (c *wishListRepo) AddOrDeleteWishList(wishList model.WishList) (string, error) {

	ctx := context.Background()
	var exists bool
	message := "product added to wishlist"

	checkquery := ` SELECT true 
					FROM wishlist 
					WHERE user_id = $1 AND product_code = $2;
					`
	deleteQuery := `DELETE FROM wishlist WHERE user_id = $1 AND product_code = $2;`

	insertQuery := `INSERT INTO wishlist(
					user_id,
					product_code)
					VALUES ($1, $2);`

	tx, err := c.db.BeginTx(ctx, nil)

	if err != nil {
		return "", err
	}

	err = tx.QueryRow(checkquery, wishList.User_Id, wishList.ProductCode).Scan(&exists)

	if err != nil && err != sql.ErrNoRows {
		tx.Rollback()
		return "", err

	}

	if err == sql.ErrNoRows {
		err = tx.QueryRow(insertQuery, wishList.User_Id, wishList.ProductCode).Err()

		if err != nil {
			tx.Rollback()
			return "", err
		}
	}

	if exists {
		err = tx.QueryRow(deleteQuery, wishList.User_Id, wishList.ProductCode).Err()
		if err != nil {
			tx.Rollback()
			return "", err
		}
		message = "product removed from wishlist "
	}

	err = tx.Commit()

	if err != nil {
		return "", err
	}

	return message, nil
}

func (c *wishListRepo) GetWishList(user_id int) ([]model.GetProduct, error) {

	query := `WITH wishlist as (	
					SELECT 
					w.product_code
					FROM wishlist w
					WHERE user_id = $1),
			  product as (	  
					  SELECT 
					  ARRAY_AGG(p.id) as id, 
					  p.name, 
					  p.code,
					  p.image,
					  p.price,
					  ARRAY_AGG(p.color) as color, 
					  p.brand_id,
					  p.category_id,
					  p.processor_id
					  FROM product p 
					  JOIN wishlist w ON p.code = w.product_code
					  GROUP BY 
					  p.name, p.brand_id, p.category_id, p.processor_id, p.code, p.image, p.price
					   )
			SELECT pr.id, pr.code, pr.name, pr.color, c.name, b.name, p.name ,pr.image, pr.price
			FROM product pr
			JOIN category c ON pr.category_id = c.id
			JOIN brand b ON pr.brand_id = b.id
			JOIN processor p ON pr.processor_id = p.id; `

	rows, err := c.db.Query(query, user_id)

	if err != nil {
		return []model.GetProduct{}, nil
	}

	defer rows.Close()

	var products []model.GetProduct

	for rows.Next() {
		var product model.GetProduct

		err := rows.Scan(
			pq.Array(&product.ID),
			&product.Code,
			&product.Name,
			pq.Array(&product.GetColor.Name),
			&product.GetCategory.Name,
			&product.GetBrand.Name,
			&product.GetProcessor.Name,
			&product.Image,
			&product.Price)

		if err != nil {
			return products, err
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return products, err
	}

	return products, nil
}
