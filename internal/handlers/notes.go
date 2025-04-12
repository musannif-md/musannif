package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"path"

	"github.com/musannif-md/musannif/internal/config"
	"github.com/musannif-md/musannif/internal/logger"
)

type noteCreateReq struct {
	Username string `json:"username"`
	NoteName string `json:"name"`
	// optional content too...
}

func CreateNote(cfg *config.AppConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req noteCreateReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// TODO: check if a note w/ the same name exists already first

		// Create file
		path := path.Join(cfg.App.NoteDirectory, req.Username, "/", req.NoteName)
		_, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 666)
		if err != nil {
			http.Error(w, "failed to create note file", http.StatusInternalServerError)
			logger.Log.Error().Err(err).Msg("failed to create note file")
			return
		}

		// TODO: optionally copy over content, if provided

		// TODO: insert external blob/file pointer in DB

		w.WriteHeader(http.StatusOK)
	}
}

func DeleteNote(cfg *config.AppConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req noteCreateReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Delete file
		path := path.Join(cfg.App.NoteDirectory, req.Username, "/", req.NoteName)
		err := os.Remove(path)
		if err != nil {
			http.Error(w, "failed to delete note file", http.StatusInternalServerError)
			logger.Log.Error().Err(err).Msg("failed to delete note file")
			return
		}

		// TODO: remove entry from DB

		w.WriteHeader(http.StatusOK)
	}
}

// CHECK: websocket...?
func FetchNote(w http.ResponseWriter, r *http.Request) {
}
