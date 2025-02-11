package main

import (
	"database/sql"

	"github.com/Srujankm12/CustomerPoBackend/internal/handlers"
	middlewares "github.com/Srujankm12/CustomerPoBackend/internal/middleware"
	"github.com/Srujankm12/CustomerPoBackend/repository"
	"github.com/gorilla/mux"
)

func registerRouter(db *sql.DB) *mux.Router {
	router := mux.NewRouter()
	router.Use(middlewares.CorsMiddleware)

	CustomerPoDropDown := handlers.NewCustomerPoController(repository.NewCustomerPoRepository(db))
	router.HandleFunc("/dropdown", CustomerPoDropDown.FetchDropDown).Methods("GET")
	router.HandleFunc("/submit", CustomerPoDropDown.SubmitFormCustomerPoData).Methods("POST")
	router.HandleFunc("/fetch", CustomerPoDropDown.FetchCustomerPoData).Methods("GET")
	return router
}
