package controllers

import (
	"net/http"

	"github.com/gieart87/gotoko/app/core/session/auth"
	"github.com/gieart87/gotoko/app/core/session/flash"
	"github.com/gieart87/gotoko/app/models"
	"github.com/google/uuid"

	"github.com/unrolled/render"
)

func (server *Server) Login(w http.ResponseWriter, r *http.Request) {
	render := render.New(render.Options{
		Layout:     "layout",
		Extensions: []string{".html", ".tmpl"},
	})

	_ = render.HTML(w, http.StatusOK, "login", map[string]interface{}{
		"error": flash.GetFlash(w, r, "error"),
	})
}

func (server *Server) DoLogin(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	userModel := models.User{}
	user, err := userModel.FindByEmail(server.DB, email)
	if err != nil {
		flash.SetFlash(w, r, "error", "email or password invalid")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if !auth.ComparePassword(password, user.Password) {
		flash.SetFlash(w, r, "error", "email or password invalid")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	session, _ := auth.GetSessionUser(r)
	session.Values["id"] = user.ID
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (server *Server) Register(w http.ResponseWriter, r *http.Request) {
	render := render.New(render.Options{
		Layout:     "layout",
		Extensions: []string{".html", ".tmpl"},
	})

	_ = render.HTML(w, http.StatusOK, "register", map[string]interface{}{
		"error": flash.GetFlash(w, r, "error"),
	})
}

func (server *Server) DoRegister(w http.ResponseWriter, r *http.Request) {
	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	email := r.FormValue("email")
	password := r.FormValue("password")

	if firstName == "" || lastName == "" || email == "" || password == "" {
		flash.SetFlash(w, r, "error", "First name, last name, email and password are required!")
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	userModel := models.User{}
	existUser, _ := userModel.FindByEmail(server.DB, email)
	if existUser != nil {
		flash.SetFlash(w, r, "error", "Sorry, email already registered")
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	hashedPassword, _ := auth.MakePassword(password)
	params := &models.User{
		ID:        uuid.New().String(),
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  hashedPassword,
	}

	user, err := userModel.CreateUser(server.DB, params)
	if err != nil {
		flash.SetFlash(w, r, "error", "Sorry, registration failed")
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	session, _ := auth.GetSessionUser(r)
	session.Values["id"] = user.ID
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (server *Server) Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := auth.GetSessionUser(r)

	session.Values["id"] = nil
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
