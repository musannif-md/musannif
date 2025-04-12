package routes

import (
	"net/http"

	"github.com/musannif-md/musannif/internal/handlers"
	_ "github.com/musannif-md/musannif/internal/middlewares"
)

func AddRoutes(mux *http.ServeMux) {
	// auth := func(handler http.HandlerFunc) http.HandlerFunc {
	// 	return middlewares.AuthMiddleware(handler)
	// }

	mux.HandleFunc("POST /login", handlers.LoginHandler)
	mux.HandleFunc("POST /signup", handlers.SignupHandler)
}
