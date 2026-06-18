package mysql

import (
	"context"

	entity "hexagonalarchitecture/internal/adapter/outbound/repository/mysql/entity"
	"hexagonalarchitecture/internal/core/domain"
)

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	const query = `
		DELETE FROM ` + entity.UserTable + `
		WHERE ` + entity.UserColumnID + ` = ?
	`

	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}
