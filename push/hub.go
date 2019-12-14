package push

import "github.com/sirupsen/logrus"

type Hub struct {
	// Registered clients.
	connections map[string]*Connection
}

var hub *Hub

func Init() {
	hub = &Hub{
		connections:    make(map[string]*Connection),
	}
}

func (h *Hub) register(connection *Connection) {
	h.connections[connection.key] = connection
	logrus.Infof("register: %v", h.connections)
}

func (h *Hub) unregister(connection *Connection) {
	delete(h.connections, connection.key)
	close(connection.send)
}


