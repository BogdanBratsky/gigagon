package ws

import (
	"encoding/json"
)

type Hub struct {
	rooms map[string]map[*Client]bool

	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan Message
}

func NewHub() *Hub {
	return &Hub{
		rooms:      make(map[string]map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan Message),
	}
}

func (h *Hub) Run() {
	for {
		select {

		case client := <-h.Register:
			if h.rooms[client.RoomID] == nil {
				h.rooms[client.RoomID] = make(map[*Client]bool)
			}
			h.rooms[client.RoomID][client] = true

		case client := <-h.Unregister:
			if clients, ok := h.rooms[client.RoomID]; ok {
				delete(clients, client)
			}

		case msg := <-h.Broadcast:
			data, _ := json.Marshal(msg)

			clients := h.rooms[msg.RoomID]
			for c := range clients {
				select {
				case c.Send <- data:
				default:
					close(c.Send)
					delete(clients, c)
				}
			}
		}
	}
}
