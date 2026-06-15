package domain

import (
	"strings"
	"time"
)

// User is the core business entity.
type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewUser(id, name, email string) User {
	now := time.Now().UTC()
	return User{
		ID:        id,
		Name:      strings.TrimSpace(name),
		Email:     strings.TrimSpace(strings.ToLower(email)),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (u *User) Update(name, email string) {
	u.Name = strings.TrimSpace(name)
	u.Email = strings.TrimSpace(strings.ToLower(email))
	u.UpdatedAt = time.Now().UTC()
}
