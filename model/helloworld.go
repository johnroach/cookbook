package model

import (
	"context"
	"cookbook/proto-gen/helloworld"
	log "github.com/sirupsen/logrus"
)


type GreetingServer struct{}

// SayHello implements helloworld.GreeterServer
func (s *GreetingServer) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	log.Infof("Received: %v", in.Name)
	return &helloworld.HelloReply{Message: "Hello " + in.Name}, nil
}