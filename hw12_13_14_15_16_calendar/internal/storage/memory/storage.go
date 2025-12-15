package memorystorage

import (
	"errors"
	"sync"
	"time"

	"github.com/DGKot/otus_hw/hw12_13_14_15_16_calendar/internal/storage"
)

type Storage struct {
	mu     sync.RWMutex
	events map[string]storage.Event
}

func New() *Storage {
	return &Storage{
		events: make(map[string]storage.Event),
	}
}

func (s *Storage) Create(event storage.Event) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if event.ID != "" {
		return event.ID, errors.New("event already created")
	}
	_, isBusy := s.Check(event)
	if isBusy {
		return "", storage.ErrDateBusy
	}
	event.ID = getID()
	s.events[event.ID] = event
	return event.ID, nil
}

func (s *Storage) Get(eventID string) (storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	event, inStorage := s.events[eventID]
	if !inStorage {
		return storage.Event{}, storage.ErrEventNotFound
	}
	return event, nil
}

func (s *Storage) Delete(eventID string) error {
	delete(s.events, eventID)
	return nil
}

func (s *Storage) Edit(event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if event.ID == "" {
		return errors.New("event not found")
	}
	s.events[event.ID] = event
	return nil
}

func (s *Storage) GetAll() ([]storage.Event, error) {
	events := make([]storage.Event, 0, len(s.events))
	for _, v := range s.events {
		events = append(events, v)
	}
	return events, nil
}

func (s *Storage) Check(event storage.Event) (storage.Event, bool) {
	if s.events == nil {
		return storage.Event{}, false
	}
	for _, v := range s.events {
		if overlaps(event.Datetime, event.Datetime.Add(event.Duration), v.Datetime, v.Datetime.Add(v.Duration)) {
			return v, true
		}
	}
	return storage.Event{}, false
}

func overlaps(start1, end1, start2, end2 time.Time) bool {
	return start1.Before(end2) && start2.Before(end1)
}

func getID() string {
	return time.Now().String()
}
