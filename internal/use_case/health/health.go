package health

import (
	"context"
	"log/slog"

	"github.com/abdivasiyev/rester/pkg/httpx"
)

type UseCase interface {
	Health(ctx context.Context, request httpx.DefaultRequest) (httpx.DefaultResponse, error)
}

type useCase struct {
	logger *slog.Logger
}

func New(logger *slog.Logger) UseCase {
	return &useCase{
		logger: logger,
	}
}

func (u *useCase) Health(ctx context.Context, request httpx.DefaultRequest) (httpx.DefaultResponse, error) {
	u.logger.Info("Health check")
	return httpx.DefaultResponse{Message: "OK"}, nil
}
