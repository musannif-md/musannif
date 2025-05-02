package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/musannif-md/musannif/internal/config"
	"github.com/musannif-md/musannif/internal/logger"
	"github.com/musannif-md/musannif/internal/resolver"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	closeDeadline = 10 * time.Second
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func validateSid(r *http.Request) (uuid.UUID, error) {
	uuidStr := r.URL.Query().Get("sid")
	if uuidStr == "" {
		return uuid.Nil, fmt.Errorf("expected session ID (uuid) via query parameter `/sid`")
	}

	sidUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("couldn't parse session ID (uuid): %w", err)
	}

	return sidUUID, nil
}

func CreateWsConn(cfg *config.AppConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			_, ok := err.(websocket.HandshakeError)
			if !ok {
				logger.Log.Error().Err(err).Msg("error establishing websocket connection")
			}
			return
		}
		defer ws.Close()

		sidUUID, err := validateSid(r)
		if err != nil {
			logger.Log.Err(err).Msgf("session ID (uuid) error")

			err := ws.WriteControl(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseUnsupportedData, "couldn't parse session ID (uuid)"),
				time.Now().Add(closeDeadline),
			)

			if err != nil {
				logger.Log.Err(err).Msgf("couldn't send close message to connection")
			}

			return
		}

		err = diffReader(sidUUID, ws)
		if err != nil {
			logger.Log.Err(err).Msg("reading ws json msg failed")
			err := ws.WriteControl(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseUnsupportedData, "expected"),
				time.Now().Add(closeDeadline),
			)

			if err != nil {
				logger.Log.Err(err).Msgf("couldn't send close message to connection [%s]", sidUUID)
			}
		}
	}
}

func diffReader(sid uuid.UUID, ws *websocket.Conn) error {
	resolver.OnClientConnect(sid, ws)
	defer resolver.OnClientDisconnect(sid, ws)

	// TODO: add ping-ponging via ticker in for loop w/ select or some shit, I don't know

	for {
		var msg map[string]any
		err := ws.ReadJSON(&msg)
		if err != nil {
			return fmt.Errorf("reading from connection with session id [%s] failed: %w", sid.String(), err)
		}

		textMsg, ok := msg["text"].(string)
		if !ok {
			return fmt.Errorf("json object didn't contain key 'text'")
		}

		err = resolver.OnClientWrite(sid, ws, textMsg)
		if err != nil {
			return fmt.Errorf("resolver failed to write in session id [%s] w/ err: %w", sid.String(), err)
		}
	}
}
