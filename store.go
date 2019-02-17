package gomatrix

import (
	"errors"
	"github.com/leighmacdonald/golm"
)

// Storer is an interface which must be satisfied to store client data.
//
// You can either write a struct which persists this data to disk, or you can use the
// provided "InMemoryStore" which just keeps data around in-memory which is lost on
// restarts.
type Storer interface {
	SaveFilterID(userID, filterID string)
	LoadFilterID(userID string) string
	SaveNextBatch(userID, nextBatchToken string)
	LoadNextBatch(userID string) string
	SaveRoom(room *Room)
	LoadRoom(roomID string) *Room
	SaveAccount(deviceId string, account golm.Account) error
	DeleteAccount(deviceId string) error
	LoadAccount(deviceId string) (golm.Account, error)
}

// InMemoryStore implements the Storer interface.
//
// Everything is persisted in-memory as maps. It is not safe to load/save filter IDs
// or next batch tokens on any goroutine other than the syncing goroutine: the one
// which called Client.Sync().
type InMemoryStore struct {
	Filters   map[string]string
	NextBatch map[string]string
	Rooms     map[string]*Room
	Published bool
	Accounts  map[string]string
}

func (s *InMemoryStore) LoadAccount(deviceId string) (golm.Account, error) {
	actStr, ok := s.Accounts[deviceId]
	if !ok {
		return golm.Account{}, errors.New("invalid device id")
	}
	act := golm.AccountFromPickle(deviceId, actStr)
	return act, nil
}

func (s *InMemoryStore) SaveAccount(deviceId string, account golm.Account) error {
	s.Accounts[deviceId] = account.Pickle(deviceId)
	return nil
}

func (s *InMemoryStore) DeleteAccount(deviceId string) error {
	delete(s.Accounts, deviceId)
	return nil
}

// SaveFilterID to memory.
func (s *InMemoryStore) SaveFilterID(userID, filterID string) {
	s.Filters[userID] = filterID
}

// LoadFilterID from memory.
func (s *InMemoryStore) LoadFilterID(userID string) string {
	return s.Filters[userID]
}

// SaveNextBatch to memory.
func (s *InMemoryStore) SaveNextBatch(userID, nextBatchToken string) {
	s.NextBatch[userID] = nextBatchToken
}

// LoadNextBatch from memory.
func (s *InMemoryStore) LoadNextBatch(userID string) string {
	return s.NextBatch[userID]
}

// SaveRoom to memory.
func (s *InMemoryStore) SaveRoom(room *Room) {
	s.Rooms[room.ID] = room
}

// LoadRoom from memory.
func (s *InMemoryStore) LoadRoom(roomID string) *Room {
	return s.Rooms[roomID]
}

// NewInMemoryStore constructs a new InMemoryStore.
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		Filters:   make(map[string]string),
		NextBatch: make(map[string]string),
		Rooms:     make(map[string]*Room),
		Accounts:  make(map[string]string),
		Published: false,
	}
}
