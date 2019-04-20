package model

import (
	"context"
	"cookbook/proto-gen/helloworld"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"testing"
)

const hwBufSize = 1024 * 1024

var hwLis *bufconn.Listener

func init() {
	hwLis = bufconn.Listen(hwBufSize)
	healthS := grpc.NewServer()

	// Registering servers
	helloworld.RegisterGreeterServer(healthS, &GreetingServer{})

	go func() {
		if err := healthS.Serve(hwLis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func hwBufDialer(context.Context, string) (net.Conn, error) {
	return hwLis.Dial()
}

// Tests SayHello Method of HelloWorld
func TestGreetingServer_SayHello(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(hwBufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := helloworld.NewGreeterClient(conn)
	resp, err := client.SayHello(ctx, &helloworld.HelloRequest{
		Name: "Dr. Seuss",
	})
	if err != nil {
		t.Fatalf("SayHello failed: %v", err)
	}

	assert.Equal(t, resp.Message, "Hello Dr. Seuss")
}