package ws

import "encoding/json"

type MessageType string

const (
	RoomUpdated  MessageType = "room_updated"
	PlayerJoined MessageType = "player_joined"
	StateChanged MessageType = "state_changed"
)

type Message struct {
	Type   MessageType `json:"type"`
	RoomID string      `json:"roomId"`
	Data   any         `json:"data"`
}

func (m Message) Marshal() ([]byte, error) {
	return json.Marshal(m)
}
