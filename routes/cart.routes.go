package routes

import (
	v1 "lapcart/handler/v1"
	m "lapcart/middleware"

	"github.com/go-chi/chi"
)

type CartRoute interface {
	CartRouter(router chi.Router,
		middleware m.Middleware,
		cartHandler v1.CartHandler)
}

type cartRoute struct{}

func NewCartRoute() CartRoute {
	return &cartRoute{}
}

func (r *cartRoute) CartRouter(
	router chi.Router,
	middleware m.Middleware,
	cartHandler v1.CartHandler) {
	router.Group(func(r chi.Router) {
		r.Use(middleware.AuthorizeJwt)
		r.Post("/user/add/cart", cartHandler.AddToCart())
		// r.Get("/user/view/wishlist", wishListHandler.ViewWishList())
	})

}
