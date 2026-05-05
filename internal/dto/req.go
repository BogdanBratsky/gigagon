package dto

type CreateRoomReq struct {
	Name string `json:"name"`
}

type JoinRoomReq struct {
	RoomID string `json:"roomId"`
	Name   string `json:"name"`
}

type StartGameReq struct {
	RoomID string `json:"roomId"`
	HostID string `json:"hostId"`
}

type SubmitAnswerReq struct {
	RoomID   string `json:"roomId"`
	PlayerID string `json:"playerId"`
	Text     string `json:"text"`
}

type VoteReq struct {
	RoomID   string `json:"roomId"`
	PlayerID string `json:"playerId"`
	AnswerID string `json:"answerId"`
}

type NextRoundReq struct {
	RoomID string `json:"roomId"`
	HostID string `json:"hostId"`
}
