package services

import (
	"errors"
	"sync"

	"task2-go-microservice/models"
)

type UserService struct {
	mu    sync.RWMutex
	store map[int]models.User
	next  int
}

func NewUserService() *UserService {
	return &UserService{store: make(map[int]models.User), next: 1}
}

func (s *UserService) List() []models.User {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]models.User, 0, len(s.store))
	for _, u := range s.store {
		result = append(result, u)
	}
	return result
}

func (s *UserService) Get(id int) (models.User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	u, ok := s.store[id]
	return u, ok
}

func (s *UserService) Create(user models.User) models.User {
	s.mu.Lock()
	defer s.mu.Unlock()

	user.ID = s.next
	s.store[user.ID] = user
	s.next++
	return user
}

func (s *UserService) Update(id int, user models.User) (models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.store[id]; !ok {
		return models.User{}, errors.New("user not found")
	}
	user.ID = id
	s.store[id] = user
	return user, nil
}

func (s *UserService) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.store[id]; !ok {
		return false
	}
	delete(s.store, id)
	return true
}
