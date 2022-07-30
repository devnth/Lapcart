package repo

import (
	"context"
	"database/sql"
	"lapcart/model"
)

type WishListRepository interface {
	AddOrDeleteWishList(wishList model.WishList) (string, error)
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
