package storage

import (
    "errors"
    "sync"
    "project/models"
)

type Storage struct {
    mu    sync.RWMutex
    data  map[int]models.User
    nextID int
}

func NewStorage() *Storage {
    return &Storage{
        data:  make(map[int]models.User),
        nextID: 1,
    }
}

func (s *Storage) CreateUser(u models.User) models.User {
    s.mu.Lock()
    defer s.mu.Unlock()
    u.ID = s.nextID
    s.nextID++
    s.data[u.ID] = u
    return u
}

func (s *Storage) GetAll() []models.User {
    s.mu.RLock()
    defer s.mu.RUnlock()
    users := make([]models.User, 0, len(s.data))
    for _, v := range s.data {
        users = append(users, v)
    }
    return users
}

func (s *Storage) GetByID(id int) (models.User, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    u, ok := s.data[id]
    if !ok {
        return models.User{}, errors.New("user not found")
    }
    return u, nil
}

func (s *Storage) Update(id int, u models.User) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    _, ok := s.data[id]
    if !ok {
        return errors.New("user not found")
    }
    u.ID = id
    s.data[id] = u
    return nil
}

func (s *Storage) Delete(id int) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    _, ok := s.data[id]
    if !ok {
        return errors.New("user not found")
    }
    delete(s.data, id)
    return nil
}
