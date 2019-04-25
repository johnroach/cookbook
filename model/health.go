package model

import (
	"context"
	health "cookbook/proto-gen/health"
	log "github.com/sirupsen/logrus"
)

// HealthServer is a struct defined for and used by health proto
type HealthServer struct{}

// Check checks the health of the running application
func (s *HealthServer) Check(ctx context.Context, in *health.HealthCheckRequest) (*health.HealthCheckResponse, error) {
	log.Debug("Health check requested")
	return &health.HealthCheckResponse{Status: health.HealthCheckResponse_SERVING}, nil
}
