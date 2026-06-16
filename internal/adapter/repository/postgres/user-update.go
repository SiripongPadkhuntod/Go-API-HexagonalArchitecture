package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"

	entity "hexagonalarchitecture/internal/adapter/repository/postgres/entity"
	"hexagonalarchitecture/internal/core/domain"
)

func (r *UserRepository) Update(ctx context.Context, user domain.User) (domain.User, error) {
	const query = `
		UPDATE ` + entity.UserTable + `
		SET ` + entity.UserColumnName + ` = $2,
			` + entity.UserColumnEmail + ` = $3,
			` + entity.UserColumnUpdatedAt + ` = $4
		WHERE ` + entity.UserColumnID + ` = $1
		RETURNING ` + entity.UserColumns + `
	`

	userEntity := entity.FromDomain(user)
	var updatedUserEntity entity.UserEntity
	err := r.db.QueryRow(ctx, query, userEntity.ID, userEntity.Name, userEntity.Email, userEntity.UpdatedAt).
		Scan(updatedUserEntity.ScanDest()...)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, err
	}

	return updatedUserEntity.ToDomain(), nil
}
