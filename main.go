package main

import (
	"cookbook/config"
	"cookbook/utils"
	"google.golang.org/grpc"
	helloworld "cookbook/proto-gen/helloworld"
	health "cookbook/proto-gen/health"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	log.Printf("Received: %v", in.Name)
	return &helloworld.HelloReply{Message: "Hello " + in.Name}, nil
}

func(s *server) Check(ctx context.Context, in *health.HealthCheckRequest) (*health.HealthCheckResponse, error){
	log.Printf("I am healthy!")
	return &health.HealthCheckResponse{Status: health.HealthCheckResponse_SERVING}, nil
}

func main() {
	environment := flag.String("e",
		utils.GetEnv("ENVIRONMENT", "dev"), "Sets the environment and pulls in config accordingly.")

	flag.Usage = func() {
		fmt.Println("Usage: server -e {mode}")
		os.Exit(1)
	}
	flag.Parse()

	// Configuration gets initialized here
	config.Init(*environment)

	c := config.GetConfig()

	lis, err := net.Listen("tcp",  ":"+ c.Get("PORT").(string))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	helloworld.RegisterGreeterServer(s, &server{})
	health.RegisterHealthServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}