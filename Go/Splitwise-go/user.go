package main

import "fmt"

type User struct {
	ID    int
	Name  string
	Email string
}

func NewUser(id int, name, email string) (*User, error) {
	if id <= 0 {
		return nil, fmt.Errorf("user ID must be a positive integer")
	}
	if name == "" {
		return nil, fmt.Errorf("user name must be a non-empty string")
	}
	if email == "" || !containsAt(email) {
		return nil, fmt.Errorf("user email must be a valid email address")
	}
	return &User{ID: id, Name: name, Email: email}, nil
}

func containsAt(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == '@' {
			return true
		}
	}
	return false
}
