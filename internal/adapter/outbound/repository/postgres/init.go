package postgres

import (
	"hexagonalarchitecture/internal/core/port"
)

var _ port.AppRepository = (*AppRepository)(nil)

type AppRepository struct {
	db Executor
}

func NewAppRepository(db Executor) *AppRepository {
	return &AppRepository{db: db}
}
