package grpc

import (
	"log"
  "net"
	// "strconv"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"location-service/internal/types"
  "location-service/api/proto/driver"
)

type GrpcHandler struct {
	config  GrpcServerConfig
	service types.TrackingService
}

type GrpcServerConfig struct {
	Port     string
	Protocol string
}

func NewGrpcHandler(svc types.TrackingService, cfg types.GrpcServerConfig) *GrpcHandler {
	return &GrpcHandler{
		config:  setConfig(cfg),
		service: svc,
	}
}

func setConfig(cfg types.GrpcServerConfig) GrpcServerConfig {
	return GrpcServerConfig{
		Port:     cfg.Port,
		Protocol: cfg.Protocol,
	}
}

// msg defined in pb generated file
func (gh *GrpcHandler) CheckHealth(ctx context.Context, msg *driver.CheckHealthRequest) (*driver.CheckHealthResponse, error) {
	return &driver.CheckHealthResponse{
    ServiceStatus: 200,
  }, nil
}



func (gh *GrpcHandler) Serve() {
	lis, err := net.Listen(gh.config.Protocol, ":9000")
	if err != nil {
		log.Fatalf(err.Error())
	}

	grpcServer := grpc.NewServer()
	// register gprc server, pass in struct with implemented function that codegen will call
	// when request is received
	driver.RegisterLocationServiceServer(grpcServer, gh)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf(err.Error())
	}
}
