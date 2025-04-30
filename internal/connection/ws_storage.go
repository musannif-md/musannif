package connection

import (
	"fmt"
	"slices"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	WS_ARR_START_CAP uint = 2
)

type sessionInfo struct {
	sockets []*websocket.Conn
	solver  DiffSolver
}

type WsMap struct {
	conns map[uuid.UUID]sessionInfo
	mu    sync.Mutex
}

var (
	m WsMap = WsMap{
		conns: make(map[uuid.UUID]sessionInfo),
	}
)

func OnClientConnect(uuid uuid.UUID, ws *websocket.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()

	si, ok := m.conns[uuid]
	wsArr := si.sockets

	if !ok {
		wsArr = make([]*websocket.Conn, 0, WS_ARR_START_CAP)
	}

	wsArr = append(wsArr, ws)
}

// Write to all clients sharing the same session
func OnClientWrite(uuid uuid.UUID, ws *websocket.Conn, msg string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	si, ok := m.conns[uuid]
	if !ok {
		return fmt.Errorf("connection doesn't exist in map")
	}

	// TODO: run it through the diff solver against all other connections in this session first

	for _, c := range si.sockets {
		err := c.WriteMessage(websocket.TextMessage, []byte("oquw!eriqupwe}r"))
		if err != nil {
			// TODO: log error if not of type client unreachable
			OnClientDisconnect(uuid, c)
		}
	}

	return nil
}

func OnClientDisconnect(uuid uuid.UUID, ws *websocket.Conn) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	si, ok := m.conns[uuid]

	if !ok {
		return fmt.Errorf("connection being removed doesn't exist")
	}

	for i, s := range si.sockets {
		if s == ws {
			si.sockets = slices.Delete(si.sockets, i, i+1)
		}
	}

	if len(si.sockets) == 0 {
		delete(m.conns, uuid)
	}

	return nil
}
