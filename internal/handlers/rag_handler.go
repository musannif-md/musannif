package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type RAGHandler struct {
	ragServiceURL string
}

func NewRAGHandler(ragServiceURL string) *RAGHandler {
	return &RAGHandler{
		ragServiceURL: ragServiceURL,
	}
}

// AddNoteToRAG adds a note's content to the RAG system
func (h *RAGHandler) AddNoteToRAG(noteContent string) error {
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
func (h *RAGHandler) QueryRAG(question string) (string, error) {
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