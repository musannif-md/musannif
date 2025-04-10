package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/masroof-maindak/musannif/internal/db"
	"github.com/masroof-maindak/musannif/internal/logger"
	"github.com/masroof-maindak/musannif/internal/utils"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignupRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type AuthResponse struct {
	Message string `json:"message"`
	Role    string `json:"role,omitempty"`
	Token   string `json:"token,omitempty"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	role, err := db.LoginUser(req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateToken(req.Username)
	if err != nil {
		logger.Log.Err(err).Msg("Failed to generate token")
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response := AuthResponse{
		Message: "Login successful",
		Role:    role,
		Token:   token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {

}
