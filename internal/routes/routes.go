package routes

import (
	"net/http"

	"github.com/musannif-md/musannif/internal/config"
	"github.com/musannif-md/musannif/internal/handlers"
	"github.com/musannif-md/musannif/internal/middlewares"
)

func AddRoutes(mux *http.ServeMux, cfg *config.AppConfig) {
	// JWT protection
	auth := func(handler http.HandlerFunc) http.HandlerFunc {
		return middlewares.AuthMiddleware(handler)
	}

	mux.HandleFunc("POST /login", handlers.LoginHandler)
	mux.HandleFunc("POST /signup", handlers.SignupHandler)

	mux.HandleFunc("POST /notes", auth(handlers.CreateNote(cfg)))
	mux.HandleFunc("DELETE /notes", auth(handlers.DeleteNote(cfg)))
}
