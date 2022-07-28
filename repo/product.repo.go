package repo

import (
	"context"
	"database/sql"
	"lapcart/model"
	"lapcart/utils"
)

type ProductRepository interface {
	AddProduct(product model.Product) (string, error)
	FindProductByCode(productCode string) (model.ProductResponse, error)
	UpdateProductColor(newStock int, inStockColor string) error
	UpdateProduct(color model.Color, product model.ProductResponse) error
	GetAllProductCode(pagenation utils.Filter) ([]string, utils.Metadata, error)
	// GetAllProducts(pagenation utils.Filter) ([]model.Product, utils.Metadata, error)
}

type productRepo struct {
	db *sql.DB
}

func NewProductRepo(db *sql.DB) ProductRepository {
	return &productRepo{
		db: db,
	}
}

func (c *productRepo) FindProductByCode(productCode string) (model.ProductResponse, error) {

	ctx := context.Background()
	var product model.ProductResponse

	query1 := `SELECT 
				product.code,
				product.name,
				product.description,
				brand.id,
				brand.name,
				category.id,
				category.name,
				processor.id,
				processor.name,
				product.price,
				product.rating,
				product.image,
				product.is_deleted
				FROM product
				INNER JOIN category ON category.id = product.category_id
				INNER  JOIN brand ON brand.id = product.brand_id
				INNER JOIN processor ON processor.id = product.processor_id
				WHERE product.code = $1;
				`
	query2 := `
				SELECT id, color, stock
				FROM product WHERE code = $1;
	`
	tx, err := c.db.BeginTx(ctx, nil)

	if err != nil {
		return model.ProductResponse{}, err
	}

	err = tx.QueryRow(query1,
		productCode).Scan(
		&product.Code,
		&product.Name,
		&product.Description,
		&product.Brand.ID,
		&product.Brand.Name,
		&product.Category.ID,
		&product.Category.Name,
		&product.Processor.ID,
		&product.Processor.Name,
		&product.Price,
		&product.Rating,
		&product.Image,
		&product.IsDeleted,
	)

	if err != nil {
		return model.ProductResponse{}, err
	}

	rows, err := tx.Query(query2, productCode)

	if err != nil {
		return model.ProductResponse{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var id uint
		var color model.Color

		if err := rows.Scan(&id, &color.Name, &color.Stock); err != nil {
			return product, err
		}

		product.Colors = append(product.Colors, color)
		product.ID = append(product.ID, id)

	}

	if err = rows.Err(); err != nil {
		return product, err
	}

	err = tx.Commit()

	if err != nil {
		return model.ProductResponse{}, err
	}

	return product, nil
}

func (c *productRepo) AddProduct(product model.Product) (string, error) {

	var color string
	var stock int
	ctx := context.Background()

	query3 := `INSERT INTO processor(name)
				VALUES($1)
				RETURNING id;`

	query2 := `INSERT INTO brand(name)
				VALUES($1)
				RETURNING id;`

	query1 := `INSERT INTO category(name)
				VALUES($1)
				RETURNING id;`

	query4 := `
				INSERT INTO product(
				code,
				name,
				color,
				brand_id,
				processor_id,
				category_id,
				price,
				stock)
				VALUES($1, $2, $3, $4, $5, $6, $7, $8)
				RETURNING code;`

	// First You begin a transaction with a call to db.Begin()
	tx, err := c.db.BeginTx(ctx, nil)

	if err != nil {
		return "", err
	}

	id, ok := c.FindCategory(product.Category.Name)

	if !ok {
		err = tx.QueryRow(query1,
			product.Category.Name).Scan(&product.Category.ID)
		if err != nil {
			// Incase we find any error in the query execution, rollback the transaction
			tx.Rollback()
			return "", err
		}
	}
	if ok {
		product.Category.ID = uint(id)
	}

	id, ok = c.FindBrand(product.Brand.Name)
	if !ok {
		err = tx.QueryRow(query2,
			product.Brand.Name).Scan(&product.Brand.ID)
		if err != nil {
			// Incase we find any error in the query execution, rollback the transaction
			tx.Rollback()
			return "", err
		}
	}
	if ok {
		product.Brand.ID = uint(id)
	}

	id, ok = c.FindProcessor(product.Processor.Name)
	if !ok {
		err = tx.QueryRow(query3,
			product.Processor.Name).Scan(&product.Processor.ID)
		if err != nil {
			// Incase we find any error in the query execution, rollback the transaction
			tx.Rollback()
			return "", err
		}
	}
	if ok {
		product.Processor.ID = uint(id)
	}

	for _, v := range product.Colors {
		color = v.Name
		stock = v.Stock

		err = tx.QueryRow(query4,
			product.Code,
			product.Name,
			color,
			product.Brand.ID,
			product.Processor.ID,
			product.Category.ID,
			product.Price,
			stock).Scan(
			&product.Code,
		)

		if err != nil {
			// Incase we find any error in the query execution, rollback the transaction
			tx.Rollback()
			return "", err
		}
	}

	err = tx.Commit()
	if err != nil {
		return "", err
	}

	return product.Code, nil

}

func (c *productRepo) UpdateProductColor(newStock int, inStockColor string) error {

	query := `UPDATE product 
				SET 
				stock = (stock + $1) 
				WHERE color = $2;`

	err := c.db.QueryRow(query, newStock, inStockColor).Err()

	if err != nil {
		return err
	}

	return nil
}

func (c *productRepo) UpdateProduct(color model.Color, product model.ProductResponse) error {

	query := `INSERT INTO product
				(code, 
				name,
				color,
				brand_id,
				processor_id,
				category_id,
				price,
				stock)
				VALUES($1, $2, $3, $4, $5, $6, $7, $8);
				`

	err := c.db.QueryRow(
		query,
		product.Code,
		product.Name,
		color.Name,
		product.Brand.ID,
		product.Processor.ID,
		product.Category.ID,
		product.Price,
		color.Stock,
	).Err()

	return err
}

func (c *productRepo) FindCategory(category string) (int, bool) {

	var id int

	query := `SELECT id FROM category WHERE name = $1;`

	err := c.db.QueryRow(query, category).Scan(&id)

	if err == sql.ErrNoRows {
		return 0, false
	}
	return id, true
}

func (c *productRepo) FindBrand(brand string) (int, bool) {

	var id int

	query := `SELECT id FROM brand WHERE name = $1;`

	err := c.db.QueryRow(query, brand).Scan(&id)

	if err == sql.ErrNoRows {
		return 0, false
	}
	return id, true
}

func (c *productRepo) FindProcessor(processor string) (int, bool) {

	var id int

	query := `SELECT id FROM processor WHERE name = $1;`

	err := c.db.QueryRow(query, processor).Scan(&id)

	if err == sql.ErrNoRows {
		return 0, false
	}
	return id, true
}

func (c *productRepo) GetAllProductCode(pagenation utils.Filter) ([]string, utils.Metadata, error) {

	var codes []string
	var totalRecords int

	query := `WITH cte AS (
		SELECT DISTINCT code FROM product)
		SELECT DISTINCT code, count(*) over() FROM cte
		LIMIT $1 OFFSET $2;`

	rows, err := c.db.Query(
		query,
		pagenation.Limit(),
		pagenation.Offset())

	if err != nil {
		return nil, utils.Metadata{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var code string

		if err = rows.Scan(&code, &totalRecords); err != nil {
			return codes, utils.ComputeMetaData(totalRecords, pagenation.Page, pagenation.PageSize), err
		}
		codes = append(codes, code)
	}

	if err := rows.Err(); err != nil {
		return codes, utils.ComputeMetaData(totalRecords, pagenation.Page, pagenation.PageSize), err
	}

	return codes, utils.ComputeMetaData(totalRecords, pagenation.Page, pagenation.PageSize), nil

}
