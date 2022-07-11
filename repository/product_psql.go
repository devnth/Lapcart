package repository

import "devnth/models"

type Prod struct {
	Id          int     `json:"product_id"`
	Name        string  `json:"product_name"`
	Description string  `json:"product_description"`
	Category    string  `josn:"product_category"`
	Brand       string  `json:"product_brand"`
	Processor   string  `json:"product_processor"`
	Price       float64 `json:"product_price"`
}

func (r Repository) Addproduct(product models.Product) (models.Product, error) {
	// Writing query
	query := `INSERT INTO product (
		product_name, 
		product_description,
		product_price,
		product_brand,
		product_processor,
		product_category)
		VALUES(
		$1, $2, $3, $4, $5, $6) 
		RETURNING 
		product_id, 
		product_name, 
		product_description,
		product_price,
		product_brand,
		product_processor,
		product_category;`
	//executing the query
	err := r.DB.QueryRow(query,
		product.Product_Name,
		product.Product_Description,
		product.Product_Price,
		product.Brand.Product_Brand_Id,
		product.Processor.Product_Processor_ID,
		product.Category.Category_ID).Scan(
		&product.Product_ID,
		&product.Product_Name,
		&product.Product_Description,
		&product.Product_Price,
		&product.Brand.Product_Brand_Id,
		&product.Processor.Product_Processor_ID,
		&product.Category.Category_ID,
	)

	return product, err

}

func (r Repository) ViewProduct() ([]Prod, error) {
	var products []Prod

	query := `SELECT 
	product.product_id, 
	product.product_name,
	product.product_description,
	product_category.category_name,
	product_branding.brand_name,
	product_processor.processor_name,
	product.product_price
	FROM product
	INNER JOIN product_category ON product_category.category_id = product.product_category
    INNER  JOIN product_branding ON product_branding.brand_id = product.product_brand
	INNER JOIN product_processor ON product_processor.processor_id = product.product_processor;`
	rows, err := r.DB.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Loop through rows, using scan to assign column data to struct fields
	for rows.Next() {
		var product Prod
		err := rows.Scan(
			&product.Id,
			&product.Name,
			&product.Description,
			&product.Category,
			&product.Brand,
			&product.Processor,
			&product.Price,
		)
		if err != nil {
			return products, err
		}
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return products, err
	}
	return products, nil

}

//function to add category
func (r Repository) AddCategory(category models.Product_Category) (models.Product_Category, error) {

	query := `INSERT INTO product_category(
		category_name,
		category_desc)
		VALUES ($1, $2) 
		RETURNING 
		category_id, 
		category_name, 
		category_desc, 
		category_created_at;`

	err := r.DB.QueryRow(query,
		category.Category_Name,
		category.Category_Description).Scan(
		&category.Category_ID,
		&category.Category_Name,
		&category.Category_Description,
		&category.Category_Created_At)

	return category, err
}

//function to add branding
func (r Repository) AddBranding(brand models.Product_Branding) (models.Product_Branding, error) {
	query := `INSERT INTO product_branding(brand_name) 
		VALUES($1) 
		RETURNING 
		brand_id, 
		brand_name, 
		brand_created_at;`
	err := r.DB.QueryRow(query,
		brand.Product_Brand_Name).Scan(
		&brand.Product_Brand_Id,
		&brand.Product_Brand_Name,
		&brand.Product_Brand_Created_At)
	return brand, err
}

//function to add processor
func (r Repository) AddProcessor(processor models.Product_Processor) (models.Product_Processor, error) {
	query := `INSERT INTO product_processor(
		processor_name, 
		processor_desc) 
		VALUES($1, $2) 
		RETURNING processor_id, 
		processor_name, 
		processor_desc,
		processor_created_at;`
	err := r.DB.QueryRow(query,
		processor.Product_Processor_Name,
		processor.Product_Processor_Description).Scan(
		&processor.Product_Processor_ID,
		&processor.Product_Processor_Name,
		&processor.Product_Processor_Description,
		&processor.Product_Processor_Created_At)
	return processor, err
}

//function to add color

func (r Repository) AddColor(color models.Product_Color) (models.Product_Color, error) {

	// writing query
	query := `INSERT INTO product_color(
			 product_id,
			 color_name)
			 VALUES
			 ($1, $2)
			 RETURNING 
			 color_id,
			 product_id,
			 color_name;`

	err := r.DB.QueryRow(query,
		color.Product_ID,
		color.Product_Color_Name).Scan(
		&color.Product_Color_ID,
		&color.Product_ID,
		&color.Product_Color_Name)

	return color, err
}

func (r Repository) ViewCategory() ([]models.Product_Category, error) {
	var categories []models.Product_Category

	//Writing and executing query
	query := `SELECT category_id, 
	category_name,
	category_desc, 
	category_created_at
	FROM product_category;`
	rows, err := r.DB.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Loop through rows, using scan to assign column data to struct fields
	for rows.Next() {
		var category models.Product_Category
		if err := rows.Scan(
			&category.Category_ID,
			&category.Category_Name,
			&category.Category_Description,
			&category.Category_Created_At); err != nil {
			return categories, err
		}
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return categories, err
	}
	return categories, nil

}

func (r Repository) ViewBranding() ([]models.Product_Branding, error) {

	//writing query
	query := `SELECT * FROM product_branding;`

	rows, err := r.DB.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var brands []models.Product_Branding

	//loop through each row
	for rows.Next() {
		var brand models.Product_Branding
		err := rows.Scan(
			&brand.Product_Brand_Id,
			&brand.Product_Brand_Name,
			&brand.Product_Brand_Created_At)
		if err != nil {
			return nil, err
		}
		brands = append(brands, brand)
	}

	if err := rows.Err(); err != nil {
		return brands, err
	}

	return brands, nil
}

func (r Repository) ViewProcessor() ([]models.Product_Processor, error) {

	//making query
	query := `SELECT * FROM product_processor;`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var processors []models.Product_Processor

	//fetching every record
	for rows.Next() {
		var processor models.Product_Processor
		err := rows.Scan(
			&processor.Product_Processor_ID,
			&processor.Product_Processor_Name,
			&processor.Product_Processor_Description,
			&processor.Product_Processor_Created_At)

		if err != nil {
			return processors, nil
		}

		processors = append(processors, processor)
	}

	if err := rows.Err(); err != nil {
		return processors, err
	}

	return processors, nil

}

func (r Repository) ViewEachProduct(product_id int) (models.Product, error) {
	var p models.Product
	query := `SELECT
	product.product_id,
	product.product_name,
	product.product_description,
	product_category.category_name,
	product_branding.brand_name,
	product_processor.processor_name
	FROM product
	INNER JOIN product_category ON product_category.category_id = product.product_category
    INNER  JOIN product_branding ON product_branding.brand_id = product.product_brand
	INNER JOIN product_processor ON product_processor.processor_id = product.product_processor
	WHERE product_id = $1;`

	err := r.DB.QueryRow(query,
		product_id).Scan(
		&p.Product_ID,
		&p.Product_Name,
		&p.Product_Description,
		&p.Category.Category_Name,
		&p.Brand.Product_Brand_Name,
		&p.Processor.Product_Processor_Name,
	)

	return p, err
}
