package user

import (
	"time"
)

type Account struct {
	ID       string
	Email    string
	Password string

	CreatedAt  time.Time
	LoggedInAt time.Time
}

func NewAccount(email Email, password Password) *Account {
	return &Account{
		Email:     email.String(),
		Password:  password.String(),
		CreatedAt: time.Now().UTC(),
	}
}

func (a *Account) Login() {
	a.LoggedInAt = time.Now().UTC()
}
