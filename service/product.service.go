package service

import (
	"database/sql"
	"lapcart/model"
	"lapcart/repo"
	"lapcart/utils"
)

type ProductService interface {
	AddProduct(product model.Product) (*string, error)
	GetProductByCode(productCode string) (*model.ProductResponse, error)
	ViewProducts(pagenation utils.Filter) (*[]model.ProductResponse, *utils.Metadata, error)
}

type productService struct {
	productRepo repo.ProductRepository
}

func NewProductService(
	productRepo repo.ProductRepository) ProductService {
	return &productService{
		productRepo: productRepo,
	}
}

func (c *productService) AddProduct(product model.Product) (*string, error) {

	productInStock, err := c.productRepo.FindProductByCode(product.Code)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if err != sql.ErrNoRows {
		flag := 0

		for _, New := range product.Colors {
			flag = 0
			for _, InStock := range productInStock.Colors {

				if New.Name == InStock.Name {
					err := c.productRepo.UpdateProductColor(New.Stock, InStock.Name)
					if err != nil {
						return nil, err
					}
					flag = 1
					break
				}
			}
			if flag == 0 {
				err := c.productRepo.UpdateProduct(New, productInStock)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	if err == sql.ErrNoRows {

		_, err := c.productRepo.AddProduct(product)

		if err != nil {
			return nil, err
		}
	}

	return &productInStock.Code, nil

}

func (c *productService) GetProductByCode(productCode string) (*model.ProductResponse, error) {

	product, err := c.productRepo.FindProductByCode(productCode)

	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (c *productService) ViewProducts(pagenation utils.Filter) (*[]model.ProductResponse, *utils.Metadata, error) {

	var products []model.ProductResponse

	productCodes, metadata, err := c.productRepo.GetAllProductCode(pagenation)

	if err != nil {
		return nil, &metadata, err
	}

	for _, productCode := range productCodes {

		product, err := c.productRepo.FindProductByCode(productCode)

		if err != nil {
			return &products, &metadata, err
		}

		products = append(products, product)
	}

	return &products, &metadata, nil

}
