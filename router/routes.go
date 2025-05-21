package router

import (
	"database/sql"
	"net/http"

	"github.com/themrgeek/airline-management-backend/controller"

	"github.com/gorilla/mux"
)

func SetupRoutes(db *sql.DB) *mux.Router {
	r := mux.NewRouter()

	// User Auth Routes
	r.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		controller.RegisterUser(w, r, db)
	}).Methods("POST")

	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		controller.LoginUser(w, r, db)
	}).Methods("POST")

	return r
}
