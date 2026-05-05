package dto

import (
	"github.com/BogdanBratsky/proverb/internal/model"
)

type RoomResponse struct {
	ID      string                   `json:"id"`
	State   string                   `json:"state"`
	Players map[string]*model.Player `json:"players"`

	Round *RoundResponse `json:"round,omitempty"`
}

type RoundResponse struct {
	Proverb model.Proverb `json:"proverb"`

	Answers []AnswerResponse  `json:"answers"`
	Votes   map[string]string `json:"votes,omitempty"`
}

type AnswerResponse struct {
	ID       string `json:"id"`
	PlayerID string `json:"playerId"`
	Text     string `json:"text"`
}

func BuildRoomResponse(room *model.Room) RoomResponse {
	resp := RoomResponse{
		ID:      room.ID,
		State:   room.State,
		Players: room.Players,
	}

	if len(room.Rounds) == 0 {
		return resp
	}

	round := room.Rounds[room.Index]

	r := RoundResponse{
		Proverb: round.Proverb,
		Votes:   round.Votes,
	}

	for _, a := range round.Answers {
		r.Answers = append(r.Answers, AnswerResponse{
			ID:       a.ID,
			PlayerID: a.PlayerId,
			Text:     a.Text,
		})
	}

	resp.Round = &r

	return resp
}

type Resp struct {
	Message  string       `json:"message"`
	Room     RoomResponse `json:"room"`
	PlayerID string       `json:"playerId,omitempty"`
}
