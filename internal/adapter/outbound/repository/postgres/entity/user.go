package entity

import (
	"time"

	"hexagonalarchitecture/internal/core/domain"
)

const (
	UserTable = "users"

	UserColumnID        = "id"
	UserColumnName      = "name"
	UserColumnEmail     = "email"
	UserColumnCreatedAt = "created_at"
	UserColumnUpdatedAt = "updated_at"

	UserColumns = UserColumnID + ", " +
		UserColumnName + ", " +
		UserColumnEmail + ", " +
		UserColumnCreatedAt + ", " +
		UserColumnUpdatedAt
)

// UserEntity is the PostgreSQL persistence shape for the users table.
// UserEntity เปรียบเสมือนรูปร่างของข้อมูลที่จะถูกเก็บในฐานข้อมูล
// UserEntity จะถูกใช้เพื่อแปลงข้อมูลระหว่าง domain.User และ PostgreSQL
type UserEntity struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func FromDomain(user domain.User) UserEntity {
	return UserEntity{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func (e UserEntity) ToDomain() domain.User {
	return domain.User{
		ID:        e.ID,
		Name:      e.Name,
		Email:     e.Email,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

func (e *UserEntity) ScanDest() []any {
	return []any{
		&e.ID,
		&e.Name,
		&e.Email,
		&e.CreatedAt,
		&e.UpdatedAt,
	}
}
