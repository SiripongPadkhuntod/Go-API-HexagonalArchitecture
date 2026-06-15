package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"hexagonalarchitecture/internal/core/domain"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(ctx context.Context, databaseURL string) (*UserRepository, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, err
	}

	repo := &UserRepository{pool: pool}
	if err := repo.ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	if err := repo.migrate(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return repo, nil
}

func (r *UserRepository) Close() {
	r.pool.Close()
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

func (r *UserRepository) ping(ctx context.Context) error {
	return r.pool.Ping(ctx)
}

func (r *UserRepository) migrate(ctx context.Context) error {
	const query = `
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL,
			updated_at TIMESTAMPTZ NOT NULL
		);

		CREATE UNIQUE INDEX IF NOT EXISTS users_email_unique ON users (email);
	`

	_, err := r.pool.Exec(ctx, query)
	return err
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
