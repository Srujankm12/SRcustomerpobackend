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

	customerPoRepo := repository.NewCustomerPoRepository(db)
	customerPoHandler := handlers.NewCustomerPoHandler(customerPoRepo)
	router.HandleFunc("/dropdown", customerPoHandler.FetchDropDown).Methods("GET")
	router.HandleFunc("/submit", customerPoHandler.SubmitFormCustomerPoData).Methods("POST")
	router.HandleFunc("/fetch", customerPoHandler.FetchCustomerPoData).Methods("GET")
	router.HandleFunc("/update", customerPoHandler.UpdateCustomerPoData).Methods("PUT")
	router.HandleFunc("/delete/{id}", customerPoHandler.DeleteCustomerPoHandler).Methods("POST")

	excelDownloadCustomerPoHandler := handlers.NewExcelDownloadCustomerPoHandler(customerPoRepo)
	router.HandleFunc("/download", excelDownloadCustomerPoHandler.DownloadCustomerPo).Methods("GET")
	return router
}
