package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"

	entity "hexagonalarchitecture/internal/adapter/outbound/repository/postgres/entity"
	"hexagonalarchitecture/internal/core/domain"
)

func (r *UserRepository) Create(ctx context.Context, user domain.User) (domain.User, error) {
	const query = `
		INSERT INTO ` + entity.UserTable + ` (` + entity.UserColumns + `)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING ` + entity.UserColumns + `
	`

	userEntity := entity.FromDomain(user)
	var createdUserEntity entity.UserEntity
	err := r.db.QueryRow(ctx, query,
		userEntity.ID,
		userEntity.Name,
		userEntity.Email,
		userEntity.CreatedAt,
		userEntity.UpdatedAt,
	).Scan(createdUserEntity.ScanDest()...)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, mapPostgresError(err)
	}

	return createdUserEntity.ToDomain(), nil
}
