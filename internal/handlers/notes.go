package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/musannif-md/musannif/internal/config"
	"github.com/musannif-md/musannif/internal/db"
	"github.com/musannif-md/musannif/internal/logger"
)

type noteCreateReq struct {
	Username string `json:"username"`
	NoteName string `json:"note_name"`
	Content  string `json:"content"`
}

type noteCreationResp struct {
	NoteId string `json:"note_id"`
}

type noteContent struct {
	Content string `json:"content"`
}

type noteListReq struct {
	Username string `json:"username"`
}

func CreateNote(cfg *config.AppConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req noteCreateReq
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.NoteName == "" || req.Username == "" {
			http.Error(w, "Username/note name not provided", http.StatusBadRequest)
			return
		}

		// Make directories
		notesDirPath := filepath.Join(cfg.App.NoteDirectory, req.Username)
		err = os.MkdirAll(notesDirPath, os.ModePerm)
		if err != nil {
			http.Error(w, "error initializing note directory: %v\n", http.StatusInternalServerError)
			logger.Log.Error().Err(err).Msgf("error initializing note directory: %s", notesDirPath)
			return
		}

		req.NoteName += ".md"
		path := filepath.Join(notesDirPath, req.NoteName)

		// Create file
		f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			http.Error(w, "failed to create note file", http.StatusInternalServerError)
			logger.Log.Error().Err(err).Msg("failed to create note file")
			return
		}

		// Write contents
		if req.Content != "" {
			_, err = f.WriteString(req.Content)
			if err != nil {
				http.Error(w, "failed to copy over contents", http.StatusInternalServerError)
				logger.Log.Error().Err(err).Msg("failed to copy over contents")
				return
			}
		}

		// Insert file info in DB
		id, err := db.CreateNote(req.Username, req.NoteName)
		if err != nil {
			http.Error(w, "failed to create note in DB", http.StatusInternalServerError)
			logger.Log.Error().Err(err).Msg("failed to create note in DB")
			return
		}

		// Construct and send response
		data := noteCreationResp{
			NoteId: strconv.FormatInt(id, 10),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}
}

func DeleteNote(cfg *config.AppConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req noteCreateReq
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.NoteName == "" || req.Username == "" {
			http.Error(w, "Username/note name not provided", http.StatusBadRequest)
			return
		}

		/*
			CHECK: how could we ensure consistency b/w the filesystem and DB?

			One horrible idea is to spin up two threads with a conditional variable in
			the database's function to rollback if the filesystem delete fails
		*/

		// Delete database entry
		req.NoteName += ".md"
		err = db.DeleteNote(req.Username, req.NoteName)
		if err != nil {
			http.Error(w, "failed to delete note from DB", http.StatusInternalServerError)
			logger.Log.Error().Err(err).Msg("failed to delete note from DB")
			return
		}

		// Delete file
		path := filepath.Join(cfg.App.NoteDirectory, req.Username, req.NoteName)
		err = os.Remove(path)
		if err != nil {
			http.Error(w, "failed to delete note file", http.StatusInternalServerError)
			logger.Log.Error().Err(err).Msg("failed to delete note file")
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func FetchNoteData(cfg *config.AppConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req noteCreateReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.NoteName == "" || req.Username == "" {
			http.Error(w, "Username/note name not provided", http.StatusBadRequest)
			return
		}

		req.NoteName += ".md"
		path := filepath.Join(cfg.App.NoteDirectory, req.Username, req.NoteName)

		content, err := os.ReadFile(path)
		if err != nil {
			http.Error(w, "failed to read note", http.StatusInternalServerError)
			logger.Log.Error().Err(err).Msg("failed to read note")
			return
		}
	
		// TODO: base64 encode contents over first...

		data := noteContent{
			Content: string(content),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}
}

func FetchNoteList(cfg *config.AppConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req noteListReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.Username == "" {
			http.Error(w, "Username not provided", http.StatusBadRequest)
			return
		}

		noteListMd, err := db.GetUserNoteMetadata(req.Username)
		if err != nil {
			http.Error(w, "failed to get metadata of user's notes", http.StatusInternalServerError)
			logger.Log.Error().Err(err).Msg("failed to get metadata of user's notes")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(noteListMd)
	}
}
