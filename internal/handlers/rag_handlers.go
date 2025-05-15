package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/musannif-md/musannif/internal/config"
)

type RAGQueryRequest struct {
	Question string `json:"question"`
}

type RAGIndexRequest struct {
	NoteName string `json:"note_name"`
}

func QueryRAGHandler(cfg *config.AppConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RAGQueryRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		ragHandler := NewRAGHandler(cfg.RAG.ServiceURL)
		response, err := ragHandler.QueryRAG(req.Question)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{
			"response": response,
		})
	}
}

func IndexNoteHandler(cfg *config.AppConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RAGIndexRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Extract username from context (JWT middleware must set this)
		username, ok := r.Context().Value("username").(string)
		if !ok || username == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Construct full note path: /base/dir/username/noteName.md
		notePath := filepath.Join(cfg.App.NoteDirectory, username, req.NoteName+".md")

		// Read content from note file
		noteContent, err := readNoteContent(notePath)
		if err != nil {
			http.Error(w, "Failed to fetch note content: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Send content to RAG handler
		ragHandler := NewRAGHandler(cfg.RAG.ServiceURL)
		if err := ragHandler.AddNoteToRAG(noteContent); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Respond success
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Note indexed successfully",
		})
	}
}

// readNoteContent reads the full content of a note file given its absolute path
func readNoteContent(fullPath string) (string, error) {
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to read note content: %v", err)
	}
	return string(content), nil
}
