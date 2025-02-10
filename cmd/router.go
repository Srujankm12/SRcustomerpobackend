package main

import (
	"database/sql"

	middlewares "github.com/Srujankm12/CustomerPoBackend/internal/middleware"
	"github.com/gorilla/mux"
)

func registerRouter(db *sql.DB) *mux.Router {
	router := mux.NewRouter()
	router.Use(middlewares.CorsMiddleware)
	return router
}
