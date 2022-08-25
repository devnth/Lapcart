package routes

import (
	v1 "lapcart/handler"
	m "lapcart/middleware"

	"github.com/go-chi/chi"
)

type WishListRoute interface {
	WishListRouter(router chi.Router,
		middleware m.Middleware,
		wishListHandler v1.WishListHandler)
}

type wishListRoute struct{}

func NewWishListRoute() WishListRoute {
	return &wishListRoute{}
}

func (r *wishListRoute) WishListRouter(router chi.Router,
	middleware m.Middleware,
	wishListHandler v1.WishListHandler) {
	router.Group(func(r chi.Router) {
		r.Use(middleware.AuthorizeJwt)
		r.Post("/user/add/wishlist", wishListHandler.AddOrDeleteWishList())
		r.Get("/user/view/wishlist", wishListHandler.ViewWishList())
	})

}
