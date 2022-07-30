package service

import (
	"lapcart/model"
	"lapcart/repo"
)

type WishListService interface {
	AddOrDeleteWishList(wishList model.WishList) (string, error)
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
