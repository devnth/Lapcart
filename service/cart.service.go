package service

import (
	"errors"
	"lapcart/model"
	"lapcart/repo"
	"log"
)

type CartService interface {
	AddToCart(cart model.Cart) (*string, error)
}

type cartService struct {
	cartRepo    repo.CartRepository
	productRepo repo.ProductRepository
}

func NewCartService(cartRepo repo.CartRepository,
	productRepo repo.ProductRepository) CartService {
	return &cartService{
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

func (c *cartService) AddToCart(cart model.Cart) (*string, error) {

	stock, err := c.productRepo.FindStockById(int(cart.Product_Id))

	if err != nil {
		log.Println("error in finding stock")
		return nil, err
	}

	if stock < cart.Count {
		log.Println("error, product out of stock, product left: ", stock)
		return nil, errors.New("product out of stock")
	}

	message, err := c.cartRepo.AddToCart(cart)

	if err != nil {
		log.Println("error adding product to cart")
		return nil, err
	}
	c.productRepo.UpdateStockById(cart)

	return &message, nil
}
