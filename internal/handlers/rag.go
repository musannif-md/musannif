package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"bytes"
	"github.com/musannif-md/musannif/internal/config"
)

type RAGHandler struct {
	ragServiceURL string
}

func newRAGHandler(ragServiceURL string) *RAGHandler {
	return &RAGHandler{
		ragServiceURL: ragServiceURL,
	}
}

type RAGQueryRequest struct {
	Question string `json:"question"`
}

type RAGIndexRequest struct {
	NoteName string `json:"note_name"`
}

func QueryRAG(cfg *config.AppConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RAGQueryRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		ragHandler := newRAGHandler(cfg.RAG.ServiceURL)
		response, err := ragHandler.queryRAG(req.Question)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{
			"response": response,
		})
	}
}

func IndexNote(cfg *config.AppConfig) http.HandlerFunc {
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
		ragHandler := newRAGHandler(cfg.RAG.ServiceURL)
		if err := ragHandler.addNoteToRAG(noteContent); err != nil {
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


// AddNoteToRAG adds a note's content to the RAG system
func (h *RAGHandler) addNoteToRAG(noteContent string) error {
	documents := []string{noteContent}
	payload := map[string][]string{
		"documents": documents,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling request: %v", err)
	}

	resp, err := http.Post(h.ragServiceURL+"/documents", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error making request to RAG service: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("RAG service returned non-200 status code: %d", resp.StatusCode)
	}

	return nil
}

// QueryRAG queries the RAG system with a question
func (h *RAGHandler) queryRAG(question string) (string, error) {
	payload := map[string]string{
		"question": question,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %v", err)
	}

	resp, err := http.Post(h.ragServiceURL+"/query", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error making request to RAG service: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("RAG service returned non-200 status code: %d", resp.StatusCode)
	}

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("error decoding response: %v", err)
	}

	return result["response"], nil
} 