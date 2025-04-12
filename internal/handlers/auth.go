package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/musannif-md/musannif/internal/db"
	"github.com/musannif-md/musannif/internal/logger"
	"github.com/musannif-md/musannif/internal/utils"
)

type loginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// type signupReq struct {
// 	Username string `json:"username"`
// 	Password string `json:"password"`
// 	Role     string `json:"role"`
// }

type authResp struct {
	Message string `json:"message"`
	Role    string `json:"role,omitempty"`
	Token   string `json:"token,omitempty"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req loginReq
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

	response := authResp{
		Message: "Login successful",
		Role:    role,
		Token:   token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	var req loginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := db.SignupUser(req.Username, req.Password, "user")
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

	response := authResp{
		Message: "Login successful",
		Role:    "user",
		Token:   token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
