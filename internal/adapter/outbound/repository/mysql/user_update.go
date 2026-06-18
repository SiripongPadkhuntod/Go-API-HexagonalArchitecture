package mysql

import (
	"context"

	entity "hexagonalarchitecture/internal/adapter/outbound/repository/mysql/entity"
	"hexagonalarchitecture/internal/core/domain"
)

func (r *UserRepository) Update(ctx context.Context, user domain.User) (domain.User, error) {
	const query = `
		UPDATE ` + entity.UserTable + `
		SET ` + entity.UserColumnName + ` = ?, ` + entity.UserColumnEmail + ` = ?, ` + entity.UserColumnUpdatedAt + ` = ?
		WHERE ` + entity.UserColumnID + ` = ?
	`

	userEntity := entity.FromDomain(user)
	res, err := r.db.ExecContext(ctx, query,
		userEntity.Name,
		userEntity.Email,
		userEntity.UpdatedAt,
		userEntity.ID,
	)
	if err != nil {
		return domain.User{}, mapMysqlError(err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return domain.User{}, err
	}

	if rowsAffected == 0 {
		return domain.User{}, domain.ErrUserNotFound
	}

	return user, nil
}
