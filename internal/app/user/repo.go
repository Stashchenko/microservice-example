package user

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stashchenko/microservice-example/internal/repository"
)

type RepositoryInterface interface {
	FindAccountByID(ctx context.Context, id int64) (*Account, error)
	AddAccount(ctx context.Context, a *Account) (*Account, error)
}

type Repository struct {
	pg *pgxpool.Pool
}

func NewUserRepository(pg *pgxpool.Pool) *Repository {
	return &Repository{pg: pg}
}

func (r *Repository) FindAccountByID(ctx context.Context, id int64) (*Account, error) {
	row := r.pg.QueryRow(ctx, "SELECT id, email, created_at FROM accounts WHERE id = $1", id)
	var a Account
	err := row.Scan(&a.ID, &a.Email, &a.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", repository.ErrNotFound, err)
	}
	return &a, nil
}

func (r *Repository) AddAccount(ctx context.Context, a *Account) (*Account, error) {
	err := r.pg.QueryRow(ctx, "INSERT INTO accounts (email, password, created_at) VALUES ($1, $2, $3) RETURNING id", a.Email, a.Password, a.CreatedAt).Scan(&a.ID)
	if err != nil {
		return nil, err
	}
	return a, nil
}
