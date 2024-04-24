package controllers

import (
	"net/http"

	"github.com/gieart87/gotoko/app/core/session/auth"

	"github.com/unrolled/render"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	render := render.New(render.Options{
		Layout:     "layout",
		Extensions: []string{".html", ".tmpl"},
	})

	user := auth.CurrentUser(server.DB, w, r)

	_ = render.HTML(w, http.StatusOK, "home", map[string]interface{}{
		"user": user,
	})
}
