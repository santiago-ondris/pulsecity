package state

import (
	"sync"

	"github.com/pulsecity/services/gateway/internal/domain"
)

type MapSnapshots struct {
	mu    sync.RWMutex
	store map[string]domain.MapClientState
}

func NewMapSnapshots() *MapSnapshots {
	return &MapSnapshots{
		store: make(map[string]domain.MapClientState),
	}
}

func (s *MapSnapshots) ApplyProgress(progress domain.MapGenerationProgress) (domain.MapClientState, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	current, existed := s.store[progress.GameID]
	if !existed {
		current = domain.MapClientState{
			GameID: progress.GameID,
		}
	}

	current.Stage = progress.Stage
	current.Progress = progress.Progress
	current.Message = progress.Message
	if progress.MapData != nil {
		current.MapData = progress.MapData
	}
	if progress.Stadium != nil {
		current.Stadium = progress.Stadium
	}

	s.store[progress.GameID] = current

	return current, existed
}

func (s *MapSnapshots) Get(gameID string) (domain.MapClientState, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state, ok := s.store[gameID]
	return state, ok
}
