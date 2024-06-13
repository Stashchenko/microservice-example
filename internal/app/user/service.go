package user

import (
	"context"
	"fmt"
	"github.com/stashchenko/microservice-example/internal/app"
)

type Service struct {
	repo RepositoryInterface
}

func NewService(repo RepositoryInterface) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) FindAccountByID(ctx context.Context, id int64) (*Account, error) {
	account, err := s.repo.FindAccountByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("find account by id: %w", err)
	}
	return account, nil
}

func (s *Service) CreateAccount(ctx context.Context, email, password string) (*Account, error) {
	var input struct {
		email    Email
		password Password
	}
	{
		var err error

		if input.email, err = NewEmail(email); err != nil {
			return nil, app.NewInvalidInputErr(err)
		}

		if input.password, err = NewPassword(password); err != nil {
			return nil, app.NewInvalidInputErr(err)
		}
	}

	account := NewAccount(input.email, input.password)
	account, err := s.repo.AddAccount(ctx, account)
	if err != nil {
		return nil, fmt.Errorf("add account: %w", err)
	}
	return account, nil
}
