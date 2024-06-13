package health

import (
	"context"
	"github.com/stashchenko/microservice-example/pkg/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type Handler struct {
	proto.UnimplementedHealthServer
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Check(_ context.Context, _ *proto.HealthCheckRequest) (*proto.HealthCheckResponse, error) {
	return &proto.HealthCheckResponse{
		Status: proto.HealthCheckResponse_SERVING,
	}, nil
}

func (h *Handler) Watch(_ *proto.HealthCheckRequest, stream proto.Health_WatchServer) error {
	response := &proto.HealthCheckResponse{
		Status: proto.HealthCheckResponse_SERVING,
	}

	for {
		select {
		case <-stream.Context().Done():
			return status.Error(codes.Canceled, "Request canceled")
		default:
			if err := stream.Send(response); err != nil {
				return err
			}
			time.Sleep(1 * time.Second)
		}
	}
}
