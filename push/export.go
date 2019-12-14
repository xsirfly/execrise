package push

import (
	"net/http"
	"github.com/sirupsen/logrus"
)

type Message struct {
	Command string `json:"command"`
	Content interface{} `json:"content"`
}

const (
	CmdRunLog = "run_log"
)

// ServeWs handles websocket requests from the peer.
func ServeWs(connId string, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.WithError(err).Error("websocket upgrade failed")
		return
	}
	connection := &Connection{key: connId, conn: conn, send: make(chan []byte, 256)}
	hub.register(connection)
	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go connection.writePump()
	go connection.readPump()
}

func GetPushClient(key string) (*Connection, bool) {
	c, ok := hub.connections[key]
	return c, ok
}
