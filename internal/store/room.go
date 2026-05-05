package store

import (
	"sync"

	"github.com/BogdanBratsky/proverb/internal/model"
)

type MemoryRoomStore struct {
	mu    sync.RWMutex
	rooms map[string]*model.Room
}

func NewMemoryRoomStore() *MemoryRoomStore {
	return &MemoryRoomStore{
		rooms: make(map[string]*model.Room),
	}
}

func (s *MemoryRoomStore) Create(room *model.Room) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.rooms[room.ID] = room
}

func (s *MemoryRoomStore) Get(id string) (*model.Room, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	room, ok := s.rooms[id]
	return room, ok
}

func (s *MemoryRoomStore) Save(room *model.Room) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.rooms[room.ID] = room
}

func (s *MemoryRoomStore) Delete(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.rooms, id)
}
