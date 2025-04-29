package utils

import (
	"fmt"
	"slices"
	"sync"

	"github.com/gorilla/websocket"
)

const (
	WS_ARR_START_CAP uint = 2
)

// TODO: Maintain map of 'session IDs' rather than 'user IDs' as keys
// A session ID can safely just be a UUID (probably)

type WsMap struct {
	conns map[uint][]*websocket.Conn
	mu    sync.Mutex
}

var (
	m WsMap = WsMap{
		conns: make(map[uint][]*websocket.Conn),
	}
)

func OnClientConnect(uid uint, ws *websocket.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()

	conns, ok := m.conns[uid]

	if !ok {
		conns = make([]*websocket.Conn, 0, WS_ARR_START_CAP)
	}

	conns = append(conns, ws)
}

func OnClientDisconnect(uid uint, ws *websocket.Conn) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	conns, ok := m.conns[uid]

	if !ok {
		return fmt.Errorf("connection being removed doesn't exist")
	}

	for i, c := range conns {
		if c == ws {
			conns = slices.Delete(conns, i, i+1)
		}
	}

	if len(conns) == 0 {
		delete(m.conns, uid)
	}

	return nil
}
