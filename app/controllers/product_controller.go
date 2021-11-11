package controllers

import (
	"net/http"

	"github.com/gieart87/gotoko/app/models"

	"github.com/unrolled/render"
)

func (s *Server) Products(w http.ResponseWriter, r *http.Request) {
	render := render.New(render.Options{
		Layout: "layout",
	})

	productModel := models.Product{}
	products, err := productModel.GetProducts(s.DB)
	if err != nil {
		return
	}

	_ = render.HTML(w, http.StatusOK, "products", map[string]interface{}{
		"products": products,
	})
}
