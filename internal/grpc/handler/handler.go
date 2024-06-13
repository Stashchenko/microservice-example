package handler

import (
	"context"
	"errors"

	"github.com/stashchenko/microservice-example/internal/app"
	"github.com/stashchenko/microservice-example/internal/app/user"
	"github.com/stashchenko/microservice-example/internal/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	User *user.Service
}

type Handler struct {
	Svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{
		Svc: svc,
	}
}

func (h *Handler) Error(ctx context.Context, err error) error {
	switch {
	case err == nil:
		return nil

	case errors.Is(err, repository.ErrNotFound):
		return status.Error(codes.NotFound, "The requested resource was not found.")

	case errors.Is(err, app.ErrInvalidInput):
		return status.Error(codes.InvalidArgument, err.Error())

	case errors.Is(err, app.ErrForbidden):
		return status.Error(codes.PermissionDenied, "You do not have permission to perform this action.")

	default:
		return status.Error(codes.Internal, "An internal error occurred.")
	}
}
