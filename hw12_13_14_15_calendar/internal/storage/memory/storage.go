package memorystorage

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/evetaell13/hw-test/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	mu       sync.RWMutex
	filename string // TODO close file
	arrays   ArraysStorage
}

type ArraysStorage struct {
	Events []storage.Event
	Users  []storage.User
}

func New(filename string) (*Storage, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	arrays := ArraysStorage{
		make([]storage.Event, 0),
		make([]storage.User, 0),
	}
	err = json.Unmarshal(bytes, &arrays)
	if err != nil {
		return nil, err
	}
	fmt.Println("json arrays: ", arrays)
	// TODO open json file RW
	return &Storage{
		filename: filename,
		arrays:   arrays,
	}, nil
}

func (s *Storage) Connect(ctx context.Context) error {
	// TODO
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	// TODO
	return nil
}

func (s *Storage) GetEventsByPeriod(start time.Time, finish time.Time) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	events := make([]storage.Event, 0, 1)
	for _, event := range s.arrays.Events {
		if (event.Start.Before(start) && event.Start.After(finish)) ||
			(event.Finish.Before(start) && event.Finish.After(finish)) {
			events = append(events, event)
		}
	}
	return events, nil
}

func (s *Storage) GetEvent(id string) (storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, e := range s.arrays.Events {
		if e.ID == id {
			return e, nil
		}
	}
	return storage.Event{}, nil
}

func (s *Storage) GetUser(id string) (storage.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, u := range s.arrays.Users {
		if u.ID == id {
			return u, nil
		}
	}
	return storage.User{}, nil
}

func (s *Storage) WriteEvent(event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if event.ID == "" {
		s.arrays.Events = append(s.arrays.Events, event)
	} else {
		return storage.ErrNotNullNewObjID
	}
	// TODO check user ID
	return nil
}

func (s *Storage) WriteUser(user storage.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if user.ID == "" {
		s.arrays.Users = append(s.arrays.Users, user)
	} else {
		return storage.ErrNotNullNewObjID
	}
	return nil
}

func (s *Storage) UpdateEvent(event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if event.ID != "" {
		for i, e := range s.arrays.Events {
			if e.ID == event.ID {
				s.arrays.Events[i] = event
				return nil
			}
		}
	} else {
		return storage.ErrNotNullNewObjID
	}
	return storage.ErrNotFoundUpdateObj
}

func (s *Storage) UpdateUser(user storage.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if user.ID != "" {
		for i, u := range s.arrays.Users {
			if u.ID == user.ID {
				s.arrays.Users[i] = user
				return nil
			}
		}
	} else {
		return storage.ErrNullUpdateObjID
	}
	return storage.ErrNotFoundUpdateObj
}
