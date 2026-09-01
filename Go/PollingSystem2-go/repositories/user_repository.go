package repositories

import (
	"fmt"
	"sync"

	"lld-practice/pollingsystem2-go/models"
)

type UserRepository struct {
	mu     sync.RWMutex
	users  []*models.User
	nextID int
}

func NewUserRepository() *UserRepository {
	return &UserRepository{users: []*models.User{}, nextID: 1}
}

func (r *UserRepository) Create(email string) (*models.User, error) {
	if email == "" {
		return nil, fmt.Errorf("invalid user email!")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for _, u := range r.users {
		if u.Email == email {
			return nil, fmt.Errorf("user already exists!")
		}
	}

	user, err := models.NewUser(r.nextID, email)
	if err != nil {
		return nil, err
	}
	r.nextID++
	r.users = append(r.users, user)
	return user, nil
}
