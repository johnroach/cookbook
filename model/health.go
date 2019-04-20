package model

import (
	"context"
	health "cookbook/proto-gen/health"
	log "github.com/sirupsen/logrus"
)

type HealthServer struct{}

func (s *HealthServer) Check(ctx context.Context, in *health.HealthCheckRequest) (*health.HealthCheckResponse, error) {
	log.Debug("Health check requested")
	return &health.HealthCheckResponse{Status: health.HealthCheckResponse_SERVING}, nil
}
