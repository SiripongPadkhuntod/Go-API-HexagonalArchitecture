package postgres

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"

	"hexagonalarchitecture/internal/core/domain"
)

const uniqueViolationCode = "23505"

func mapPostgresError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == uniqueViolationCode {
		return fmt.Errorf("%w: email already exists", domain.ErrUserAlreadyExists)
	}

	return err
}
