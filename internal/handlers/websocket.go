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
	pongWait      = 20 * time.Second
	pingPeriod    = pongWait * 9 / 10
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
		return uuid.Nil, fmt.Errorf("expected session ID via query parameter `/sid`")
	}

	sidUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("couldn't parse session ID: %w", err)
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

		err = wsTransmission(cfg, sidUUID, ws, r)
		if err != nil {
			logger.Log.Err(err).Msg("websocket died: ")
		}
	}
}

func wsTransmission(cfg *config.AppConfig, sid uuid.UUID, ws *websocket.Conn, r *http.Request) error {
	readerFinished := make(chan error)
	pingTicker := time.NewTicker(pingPeriod)
	err := resolver.OnClientConnect(cfg, sid, ws, r)
	if err != nil {
		err2 := ws.WriteControl(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseUnsupportedData, err.Error()),
			time.Now().Add(closeDeadline),
		)

		if err2 != nil {
			logger.Log.Err(err2).Msgf("couldn't send close message to connection")
		}

		return err
	}

	defer func() {
		resolver.OnClientDisconnect(sid, ws)
		pingTicker.Stop()
		ws.Close()
	}()

	go readWs(sid, ws, readerFinished)

	for {
		select {
		case reason := <-readerFinished: // `readWs` loop has been broken
			return reason
		case <-pingTicker.C: // Send sporadic pings
			err := ws.WriteControl(
				websocket.PingMessage,
				[]byte{},
				time.Now().Add(closeDeadline),
			)

			if err != nil {
				return fmt.Errorf("failed to send ping: %w", err)
			}
		}
	}
}

func readWs(sid uuid.UUID, ws *websocket.Conn, readerFinished chan error) {
	defer close(readerFinished)

	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(appdata string) error {
		ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		var msg map[string]any

		err := ws.ReadJSON(&msg)
		if err != nil {
			if _, ok := err.(*websocket.CloseError); ok {
				readerFinished <- nil
			} else {
				readerFinished <- fmt.Errorf("reading from connection with session id [%s] failed: %w", sid.String(), err)
			}
			break
		}

		textMsg, ok := msg["text"].(string)
		if !ok {
			readerFinished <- fmt.Errorf("json object didn't contain key 'text'")
			break
		}

		err = resolver.OnClientWrite(sid, ws, textMsg)
		if err != nil {
			readerFinished <- fmt.Errorf("resolver failed to write in session id [%s] w/ err: %w", sid.String(), err)
			break
		}
	}
}
