package main

type User struct {
	ID    int
	Email string
}

func NewUser(id int, email string) (*User, error) {
	if id == 0 || email == "" {
		return nil, errInvalidUser
	}
	return &User{ID: id, Email: email}, nil
}
