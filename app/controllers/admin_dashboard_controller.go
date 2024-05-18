package controllers

import (
	"fmt"
	"net/http"

	"github.com/gieart87/gotoko/app/utils"

	"github.com/gieart87/gotoko/app/core/session/auth"

	"github.com/unrolled/render"
)

func (server *Server) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	render := render.New(render.Options{
		Layout:     "admin_layout",
		Extensions: []string{".html", ".tmpl"},
	})

	user := auth.CurrentUser(server.DB, w, r)

	fmt.Println("user ===>", utils.PrintJSON(user))

	_ = render.HTML(w, http.StatusOK, "admin_dashboard", map[string]interface{}{
		"user": utils.PrintJSON(user),
	})
}
