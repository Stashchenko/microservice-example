package userhandler

import (
	"context"
	"github.com/stashchenko/microservice-example/internal/grpc/handler"
	"github.com/stashchenko/microservice-example/pkg/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log/slog"
	"strconv"
)

type Handler struct {
	proto.UnimplementedUserServer

	*handler.Handler
}

func NewHandler(handler *handler.Handler) *Handler {
	return &Handler{
		Handler: handler,
	}
}
func (h *Handler) GetUser(ctx context.Context, p *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	user, err := h.Svc.User.FindAccountByID(ctx, p.UserId)
	if err != nil {
		slog.ErrorContext(ctx, "failed to find account", "error", err)
		return nil, h.Error(ctx, err)
	}
	num, _ := strconv.ParseInt(user.ID, 10, 64)
	return &proto.GetUserResponse{
		UserId:    num,
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}, nil
}

func (h *Handler) Signup(ctx context.Context, p *proto.SignUpRequest) (*proto.SignupResponse, error) {
	user, err := h.Svc.User.CreateAccount(ctx, p.Email, p.Password)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create account", "error", err)

		return nil, h.Error(ctx, err)
	}

	return &proto.SignupResponse{
		UserId: user.ID,
	}, nil
}
