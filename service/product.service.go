package service

import (
	"database/sql"
	"errors"
	"lapcart/model"
	"lapcart/repo"
	"lapcart/utils"
	"log"
)

type ProductService interface {
	AddProduct(product model.Product) (*string, error)
	GetProductByCode(productCode string) (*model.ProductResponse, error)
	ViewProducts(pagenation utils.Filter) (*[]model.ProductResponse, *utils.Metadata, error)
	UpdateProduct(product model.UpdateProduct) error
	DeleteProduct(data model.DeleteProduct) error
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

func (c *productService) UpdateProduct(data model.UpdateProduct) error {

	var err error

	data.OldCode, err = c.productRepo.GetProductCodeById(data.ProductID)

	if err != nil {
		log.Println(err)
		return errors.New("the product you are looking for is not in the database")
	}

	err = c.productRepo.UpdateProductByCode(data)

	if err != nil {
		log.Println("error in updating products")
		return err
	}

	if data.ChangeColor != "" {
		err = c.productRepo.ChangeColor(data)

		if err != nil {
			return err
		}
	}

	if data.ChangeQuantity != 0 {
		err = c.productRepo.ChangeStock(data)

		if err != nil {
			return err
		}
	}

	if data.NewColor != "" {

		if data.NewQuantity != 0 {

			err = c.productRepo.InsertNewColor(data)

			if err != nil {
				return err
			}

			return nil
		}

		return errors.New("quantity for new color not given")

	}

	return nil
}

func (c *productService) DeleteProduct(data model.DeleteProduct) error {

	err := c.productRepo.DeleteProduct(data)

	if err != nil {
		return err
	}

	return err

}
