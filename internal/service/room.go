package service

import (
	"time"

	"github.com/BogdanBratsky/proverb/internal/model"
	"github.com/BogdanBratsky/proverb/internal/ws"
	"github.com/BogdanBratsky/proverb/pkg/id"
)

const (
	WritingDuration = 2 * time.Minute
	VotingDuration  = 2 * time.Minute
)

type RoomStore interface {
	Create(room *model.Room)
	Get(id string) (*model.Room, bool)
	Save(room *model.Room)
	Delete(id string)
}

type ProverbProvider interface {
	GetProverb()
}

type RoomService struct {
	store    RoomStore
	proverbs ProverbProvider
	hub      *ws.Hub
}

func NewRoomService(store RoomStore, proverbs ProverbProvider, hub *ws.Hub) *RoomService {
	return &RoomService{
		store:    store,
		proverbs: proverbs,
		hub:      hub,
	}
}

// =========================
// helpers
// =========================

func (s *RoomService) broadcast(room *model.Room) {
	s.hub.Broadcast <- ws.Message{
		Type:   ws.RoomUpdated,
		RoomID: room.ID,
		Data:   room,
	}
}

func (s *RoomService) GetRoom(id string) (*model.Room, error) {
	room, ok := s.store.Get(id)
	if !ok {
		return nil, ErrRoomNotFound
	}
	return room, nil
}

func (s *RoomService) startTimer(roomID string, state string) {
	duration := WritingDuration

	if state == model.RoomVoting {
		duration = VotingDuration
	}

	time.AfterFunc(duration, func() {
		s.autoAdvance(roomID, state)
	})
}
func (s *RoomService) autoAdvance(roomID string, state string) {
	room, ok := s.store.Get(roomID)
	if !ok {
		return
	}

	if room.State != state {
		return
	}

	switch state {

	case model.RoomWriting:
		room.State = model.RoomVoting
		s.store.Save(room)
		s.broadcast(room)

		// ⏱ старт voting таймера
		s.startTimer(roomID, model.RoomVoting)
		return

	case model.RoomVoting:
		s.calculateScores(room)
		s.finishRound(room)

		room.Index++

		if s.isGameOver(room) {
			room.State = model.RoomResults
			s.store.Save(room)
			s.broadcast(room)
			return
		}

		round := model.Round{
			Proverb:  room.Proverbs[room.Index],
			Answers:  make(map[string]model.Answer),
			Votes:    make(map[string]string),
			Answered: make(map[string]bool),
			Voted:    make(map[string]bool),
		}

		room.Rounds = append(room.Rounds, round)
		room.State = model.RoomWriting

		s.store.Save(room)
		s.broadcast(room)

		// ⏱ старт writing таймера
		s.startTimer(roomID, model.RoomWriting)
	}
}

// =========================
// room lifecycle
// =========================

func (s *RoomService) CreateRoom(hostName string) (*model.Room, error) {
	host := &model.Player{
		ID:   id.GenerateID(),
		Name: hostName,
	}

	room := &model.Room{
		ID:      id.GenerateID(),
		HostID:  host.ID,
		State:   model.RoomWaiting,
		Players: make(map[string]*model.Player),
		Proverbs: []model.Proverb{
			{ID: "1", Text: "Без труда не вытащишь ___"},
			{ID: "2", Text: "Любишь кататься — ___"},
			{ID: "3", Text: "Семь раз отмерь ___"},
		},
		Index:  0,
		Rounds: []model.Round{},
	}

	room.Players[host.ID] = host

	s.store.Create(room)

	s.broadcast(room)

	return room, nil
}

func (s *RoomService) JoinRoom(roomID string, name string) (*model.Player, *model.Room, error) {
	room, ok := s.store.Get(roomID)
	if !ok {
		return nil, nil, ErrRoomNotFound
	}

	if room.State != model.RoomWaiting {
		return nil, nil, ErrAlreadyStarted
	}

	player := &model.Player{
		ID:   id.GenerateID(),
		Name: name,
	}

	room.Players[player.ID] = player

	s.store.Save(room)
	s.broadcast(room)

	return player, room, nil
}

func (s *RoomService) StartGame(roomID, hostID string) (*model.Room, error) {
	room, ok := s.store.Get(roomID)
	if !ok {
		return nil, ErrRoomNotFound
	}

	if room.HostID != hostID {
		return nil, ErrNotHost
	}

	if room.State != model.RoomWaiting {
		return nil, ErrAlreadyStarted
	}

	room.Index = 0
	room.State = model.RoomWriting

	round := model.Round{
		Proverb:  room.Proverbs[room.Index],
		Answers:  make(map[string]model.Answer),
		Votes:    make(map[string]string),
		Answered: make(map[string]bool),
		Voted:    make(map[string]bool),
	}

	room.Rounds = append(room.Rounds, round)

	s.store.Save(room)
	s.broadcast(room)

	// ⏱ СТАРТ ТАЙМЕРА НА WRITING
	s.startTimer(room.ID, model.RoomWriting)

	return room, nil
}

// =========================
// gameplay
// =========================

func (s *RoomService) SubmitAnswer(roomID, playerID, text string) error {
	room, ok := s.store.Get(roomID)
	if !ok {
		return ErrRoomNotFound
	}

	if room.State != model.RoomWriting {
		return ErrInvalidState
	}

	round := &room.Rounds[room.Index]

	if round.Answered[playerID] {
		return ErrAlreadyAnswered
	}

	answer := model.Answer{
		ID:       id.GenerateID(),
		PlayerId: playerID,
		Text:     text,
	}

	round.Answers[answer.ID] = answer
	round.Answered[playerID] = true

	if len(round.Answered) == len(room.Players) {
		room.State = model.RoomVoting
	}

	if len(round.Answered) == len(room.Players) {
		room.State = model.RoomVoting

		s.store.Save(room)
		s.broadcast(room)

		// ⏱ старт таймера voting
		s.startTimer(room.ID, model.RoomVoting)

		return nil
	}

	return nil
}

func (s *RoomService) Vote(roomID, playerID, answerID string) error {
	room, ok := s.store.Get(roomID)
	if !ok {
		return ErrRoomNotFound
	}

	if room.State != model.RoomVoting {
		return ErrInvalidState
	}

	round := &room.Rounds[room.Index]

	if round.Voted[playerID] {
		return ErrAlreadyVoted
	}

	answer, ok := round.Answers[answerID]
	if !ok {
		return ErrInvalidAnswer
	}

	if answer.PlayerId == playerID {
		return ErrCantVoteSelf
	}

	round.Votes[playerID] = answerID
	round.Voted[playerID] = true

	s.store.Save(room)
	s.broadcast(room)

	return nil
}

func (s *RoomService) NextRound(roomID, hostID string) (*model.Room, error) {
	room, ok := s.store.Get(roomID)
	if !ok {
		return nil, ErrRoomNotFound
	}

	if room.HostID != hostID {
		return nil, ErrNotHost
	}

	s.calculateScores(room)
	s.finishRound(room)

	room.Index++

	if s.isGameOver(room) {
		s.calculateScores(room)
		s.finishRound(room)

		room.State = model.RoomResults

		s.store.Save(room)
		s.broadcast(room)

		return room, nil
	}

	round := model.Round{
		Proverb:  room.Proverbs[room.Index],
		Answers:  make(map[string]model.Answer),
		Votes:    make(map[string]string),
		Answered: make(map[string]bool),
		Voted:    make(map[string]bool),
	}

	room.Rounds = append(room.Rounds, round)
	room.State = model.RoomWriting

	s.store.Save(room)
	s.broadcast(room)

	return room, nil
}

func (s *RoomService) calculateScores(room *model.Room) {
	round := &room.Rounds[room.Index]

	for _, answerID := range round.Votes {
		answer := round.Answers[answerID]
		room.Players[answer.PlayerId].Score++
	}
}

func (s *RoomService) finishRound(room *model.Room) {
	room.Rounds[room.Index].Done = true
}

func (s *RoomService) isGameOver(room *model.Room) bool {
	return room.Index >= len(room.Proverbs)-1
}
