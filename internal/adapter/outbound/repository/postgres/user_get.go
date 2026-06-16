package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"

	entity "hexagonalarchitecture/internal/adapter/outbound/repository/postgres/entity"
	"hexagonalarchitecture/internal/core/domain"
)

func (r *AppRepository) FindAll(ctx context.Context) ([]domain.User, error) {
	const query = `
		SELECT ` + entity.UserColumns + `
		FROM ` + entity.UserTable + `
		ORDER BY ` + entity.UserColumnCreatedAt + ` ASC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]domain.User, 0)
	for rows.Next() {
		var userEntity entity.UserEntity
		if err := rows.Scan(userEntity.ScanDest()...); err != nil {
			return nil, err
		}
		users = append(users, userEntity.ToDomain())
	}

	return users, rows.Err()
}

func (r *AppRepository) FindByID(ctx context.Context, id string) (domain.User, error) {
	const query = `
		SELECT ` + entity.UserColumns + `
		FROM ` + entity.UserTable + `
		WHERE ` + entity.UserColumnID + ` = $1
	`

	var userEntity entity.UserEntity
	err := r.db.QueryRow(ctx, query, id).Scan(userEntity.ScanDest()...)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, err
	}

	return userEntity.ToDomain(), nil
}
