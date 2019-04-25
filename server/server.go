package server

import (
	"cookbook/config"
	"cookbook/model"
	health "cookbook/proto-gen/health"
	"cookbook/proto-gen/helloworld"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

// Initializes the service by pulling in the configurations and registering the GRPC servers
func Init() {
	c := config.GetConfig()

	lis, err := net.Listen("tcp", ":"+config.GetWithDefault("PORT", "8080").(string))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Infof("Listening on :%v", c.Get("PORT").(string))

	s := grpc.NewServer()

	helloworld.RegisterGreeterServer(s, &model.GreetingServer{})
	health.RegisterHealthServer(s, &model.HealthServer{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
