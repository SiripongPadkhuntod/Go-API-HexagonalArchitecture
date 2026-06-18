package mysql

import (
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"

	"hexagonalarchitecture/internal/core/domain"
)

func mapMysqlError(err error) error {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 { // 1062 is Duplicate entry for key
		return fmt.Errorf("%w: email already exists", domain.ErrUserAlreadyExists)
	}

	return err
}
