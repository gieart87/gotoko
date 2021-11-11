package controllers

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) initializeRoutes() {
	s.Router = mux.NewRouter()
	s.Router.HandleFunc("/", s.Home).Methods("GET")
	s.Router.HandleFunc("/products", s.Products).Methods("GET")

	staticFileDirectory := http.Dir("./assets/")
	staticFileHandler := http.StripPrefix("/public/", http.FileServer(staticFileDirectory))
	s.Router.PathPrefix("/public/").Handler(staticFileHandler).Methods("GET")
}
