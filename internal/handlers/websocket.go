package handlers

import (
	"net/http"

	"github.com/musannif-md/musannif/internal/config"
	"github.com/musannif-md/musannif/internal/logger"
	"github.com/musannif-md/musannif/internal/utils"

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

		// TODO: get UID from body
		var uid uint = 1

		utils.OnClientConnect(uid, ws)
		defer utils.OnClientDisconnect(uid, ws)

		// TODO: get file name from body & check if the (newly-met) client against this
		// particular UID is attempting to open an already-opened file or a new file

		go diffReader(uid, ws)
		diffWriter(uid, ws)
	}
}

func diffReader(uid uint, ws *websocket.Conn) error {
	defer func() {
		ws.Close()
		utils.OnClientDisconnect(uid, ws)
	}()

	for {
		var msg map[string]any
		err := ws.ReadJSON(&msg)
		if err != nil {
			logger.Log.Err(err).Msgf("reading from connection with id [%d] failed", uid)
			break
		}

		// TODO: ping/pong and auto-disconnect after inactivity?

		// TODO: Propagate new diff to diff handler
	}

	return nil
}

func diffWriter(uid uint, ws *websocket.Conn) error {
	defer func() {
		ws.Close()
		utils.OnClientDisconnect(uid, ws)
	}()

	/* TODO: Propagate 'updated diffs (post conflict-resolution)' to all
	listeners i.e all websocket connections against a uid/sid */
	return nil
}
