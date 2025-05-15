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

	// Auth
	mux.HandleFunc("POST /login", handlers.LoginHandler)
	mux.HandleFunc("POST /signup", handlers.SignupHandler)

	// Single note
	mux.HandleFunc("POST /note", auth(handlers.CreateNote(cfg)))        // Upload a note to the user's directory
	mux.HandleFunc("POST /get-note", auth(handlers.FetchNoteData(cfg))) // Get the contents of one note in the user's directory
	mux.HandleFunc("POST /del-note", auth(handlers.DeleteNote(cfg)))    // Delete a note from the user's directory

	// User & note metadata
	mux.HandleFunc("POST /notes", auth(handlers.FetchNoteList(cfg))) // Return a list of notes in user's directory

	// RAG endpoints
	mux.HandleFunc("POST /rag/query", auth(handlers.QueryRAG(cfg))) // Query the RAG system
	mux.HandleFunc("POST /rag/index", auth(handlers.IndexNote(cfg))) // Index a note in the RAG system

	// Connection
	mux.HandleFunc("/connect", auth(handlers.CreateWsConn(cfg))) // Establish connection and start sending/receiving diffs
}
