package handlers

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/musannif-md/musannif/internal/config"
	"github.com/musannif-md/musannif/internal/connection"
	"github.com/musannif-md/musannif/internal/logger"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

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

		uuidStr := r.URL.Query().Get("sid")
		if uuidStr == "" {
			http.Error(w, "expected sid (uuid) via query parameter", http.StatusBadRequest)
			return
		}

		uuid, err := uuid.Parse(uuidStr)
		if err != nil {
			http.Error(w, "couldn't parse UUID", http.StatusBadRequest)
			return
		}

		connection.OnClientConnect(uuid, ws)
		defer connection.OnClientDisconnect(uuid, ws)

		var wg sync.WaitGroup
		wg.Add(1)

		go diffReader(uuid, ws, &wg)

		wg.Wait()
	}
}

func diffReader(uuid uuid.UUID, ws *websocket.Conn, wg *sync.WaitGroup) error {
	defer func() {
		ws.Close()
		connection.OnClientDisconnect(uuid, ws)
		wg.Done()
	}()

	for {
		var msg map[string]any
		err := ws.ReadJSON(&msg)
		if err != nil {
			logger.Log.Err(err).Msgf("reading from connection with id [%s] failed", uuid.String())
			break
		}

		// TODO: ping/pong and auto-disconnect after inactivity?

		textMsg, ok := msg["text"].(string)
		if !ok {
			return fmt.Errorf("json object didn't contain map")
		}

		connection.OnClientWrite(uuid, ws, textMsg)
	}

	return nil
}
