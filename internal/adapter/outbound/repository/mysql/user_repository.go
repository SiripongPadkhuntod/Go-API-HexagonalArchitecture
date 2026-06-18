package mysql

import (
	"hexagonalarchitecture/internal/core/port"
)

var _ port.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
	db Executor
}

func NewUserRepository(db Executor) *UserRepository {
	return &UserRepository{db: db}
}
