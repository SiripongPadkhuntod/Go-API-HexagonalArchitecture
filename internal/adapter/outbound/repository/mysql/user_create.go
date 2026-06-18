package mysql

import (
	"context"

	entity "hexagonalarchitecture/internal/adapter/outbound/repository/mysql/entity"
	"hexagonalarchitecture/internal/core/domain"
)

func (r *UserRepository) Create(ctx context.Context, user domain.User) (domain.User, error) {
	const query = `
		INSERT INTO ` + entity.UserTable + ` (` + entity.UserColumns + `)
		VALUES (?, ?, ?, ?, ?)
	`

	userEntity := entity.FromDomain(user)
	_, err := r.db.ExecContext(ctx, query,
		userEntity.ID,
		userEntity.Name,
		userEntity.Email,
		userEntity.CreatedAt,
		userEntity.UpdatedAt,
	)
	if err != nil {
		return domain.User{}, mapMysqlError(err)
	}

	return user, nil
}
