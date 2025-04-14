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

	mux.HandleFunc("GET /note", auth(handlers.FetchNoteData(cfg))) // Get the contents of one note in the user's directory
	mux.HandleFunc("POST /note", auth(handlers.CreateNote(cfg)))   // Upload a note to the user's directory
	mux.HandleFunc("DELETE /note", auth(handlers.DeleteNote(cfg))) // Delete a note from the user's directory

	mux.HandleFunc("GET /notes", auth(handlers.FetchNoteList(cfg))) // Return a list of notes in user's directory
}
