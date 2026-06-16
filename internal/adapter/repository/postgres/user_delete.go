package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"

	entity "hexagonalarchitecture/internal/adapter/repository/postgres/entity"
	"hexagonalarchitecture/internal/core/domain"
)

func (r *AppRepository) Delete(ctx context.Context, id string) error {
	const query = `
		DELETE FROM ` + entity.UserTable + `
		WHERE ` + entity.UserColumnID + ` = $1
		RETURNING ` + entity.UserColumns + `
	`

	var deletedUserEntity entity.UserEntity
	err := r.db.QueryRow(ctx, query, id).Scan(deletedUserEntity.ScanDest()...)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.ErrUserNotFound
		}
		return err
	}

	return nil
}
