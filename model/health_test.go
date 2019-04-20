package model

import (
	"context"
	cbhealth "cookbook/proto-gen/health"
	"fmt"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"testing"
)

const healthBufSize = 1024 * 1024

var healthLis *bufconn.Listener

func init() {
	healthLis = bufconn.Listen(healthBufSize)
	healthS := grpc.NewServer()

	// Registering servers
	cbhealth.RegisterHealthServer(healthS, &HealthServer{})

	go func() {
		if err := healthS.Serve(healthLis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func healthBufDialer(context.Context, string) (net.Conn, error) {
	return healthLis.Dial()
}

func TestHealthServer_Check(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(healthBufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := cbhealth.NewHealthClient(conn)
	resp, err := client.Check(ctx, &cbhealth.HealthCheckRequest{})
	if err != nil {
		t.Fatalf("Health Check failed: %v", err)
	}
	fmt.Print(resp.Status)
	assert.Equal(t, resp.Status.String(), "SERVING")
}