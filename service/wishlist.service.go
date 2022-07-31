package service

import (
	"lapcart/model"
	"lapcart/repo"
)

type WishListService interface {
	AddOrDeleteWishList(wishList model.WishList) (string, error)
	GetWishList(user_id int) (*[]model.GetProduct, error)
}

type wishListService struct {
	wishListRepo repo.WishListRepository
}

func NewWishListService(
	wishListRepo repo.WishListRepository,
) WishListService {
	return &wishListService{
		wishListRepo: wishListRepo,
	}
}

func (c *wishListService) AddOrDeleteWishList(wishList model.WishList) (string, error) {

	message, err := c.wishListRepo.AddOrDeleteWishList(wishList)

	if err != nil {
		return "", err
	}

	return message, nil

}

func (c *wishListService) GetWishList(user_id int) (*[]model.GetProduct, error) {

	wishList, err := c.wishListRepo.GetWishList(user_id)

	if err != nil {
		return nil, err
	}

	return &wishList, err

}
