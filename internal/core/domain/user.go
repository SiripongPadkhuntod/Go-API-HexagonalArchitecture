package domain

import (
	"strings"
	"time"
)

// User is the core business entity.
type User struct {
	ID        string
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(id, name, email string, now time.Time) User {
	return User{
		ID:        id,
		Name:      strings.TrimSpace(name),
		Email:     strings.TrimSpace(strings.ToLower(email)),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (u *User) Update(name, email string, now time.Time) {
	u.Name = strings.TrimSpace(name)
	u.Email = strings.TrimSpace(strings.ToLower(email))
	u.UpdatedAt = now
}
