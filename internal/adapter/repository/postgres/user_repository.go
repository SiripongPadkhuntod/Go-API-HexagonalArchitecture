package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"hexagonalarchitecture/internal/core/domain"
	"hexagonalarchitecture/internal/core/port"
)

var _ port.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) Create(ctx context.Context, user domain.User) (domain.User, error) {
	const query = `
		INSERT INTO users (id, name, email, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, email, created_at, updated_at
	`

	return r.scanUser(ctx, query, user.ID, user.Name, user.Email, user.CreatedAt, user.UpdatedAt)
}

func (r *UserRepository) FindAll(ctx context.Context) ([]domain.User, error) {
	const query = `
		SELECT id, name, email, created_at, updated_at
		FROM users
		ORDER BY created_at ASC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]domain.User, 0)
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, rows.Err()
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (domain.User, error) {
	const query = `
		SELECT id, name, email, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	return r.scanUser(ctx, query, id)
}

func (r *UserRepository) Update(ctx context.Context, user domain.User) (domain.User, error) {
	const query = `
		UPDATE users
		SET name = $2, email = $3, updated_at = $4
		WHERE id = $1
		RETURNING id, name, email, created_at, updated_at
	`

	return r.scanUser(ctx, query, user.ID, user.Name, user.Email, user.UpdatedAt)
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	const query = `DELETE FROM users WHERE id = $1`

	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) scanUser(ctx context.Context, query string, args ...any) (domain.User, error) {
	var user domain.User
	err := r.pool.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, err
	}

	return user, nil
}
