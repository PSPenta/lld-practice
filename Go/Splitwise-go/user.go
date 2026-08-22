package main

type User struct {	Email string
}

func NewUser(email string) *User {
	return &User{Email: email}
}
