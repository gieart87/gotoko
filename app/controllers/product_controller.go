package controllers

import (
	"net/http"
	"strconv"

	"github.com/gieart87/gotoko/app/models"

	"github.com/unrolled/render"
)

func (server *Server) Products(w http.ResponseWriter, r *http.Request) {
	render := render.New(render.Options{
		Layout: "layout",
	})

	var page = 1
	if r.URL.Query().Get("page") != "" {
		page, _ = strconv.Atoi(r.URL.Query().Get("page"))
	}

	productModel := models.Product{}
	products, err := productModel.GetProducts(server.DB, page)
	if err != nil {
		return
	}

	_ = render.HTML(w, http.StatusOK, "products", map[string]interface{}{
		"products": products,
	})
}
