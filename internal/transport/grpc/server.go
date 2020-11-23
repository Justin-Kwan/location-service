package grpc

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

	// "location-service/api/proto/courier"
	// "location-service/api/proto/health"
	"location-service/api/proto/order"
	"location-service/internal"
	"location-service/internal/transport"
	"location-service/internal/types"
)

type GrpcHandler struct {
	config  *GrpcServerConfig
	service transport.TrackingService
}

type GrpcServerConfig struct {
	Port     string
	Protocol string
}

func NewGrpcHandler(svc transport.TrackingService, cfg *types.GrpcServerConfig) *GrpcHandler {
	return &GrpcHandler{
		config:  setConfig(cfg),
		service: svc,
	}
}

func setConfig(cfg *types.GrpcServerConfig) *GrpcServerConfig {
	return &GrpcServerConfig{
		Port:     cfg.Port,
		Protocol: cfg.Protocol,
	}
}

func (h *GrpcHandler) Serve() error {
	lis, err := net.Listen(h.config.Protocol, ":9000")
	if err != nil {
		log.Fatalf(err.Error())
	}

	grpcServer := grpc.NewServer()
	h.registerServices(grpcServer)

	log.Printf("Grpc server started...")
	return grpcServer.Serve(lis)
}

func (h *GrpcHandler) registerServices(grpcServer *grpc.Server) {
	// health.RegisterHealthServiceServer(grpcServer, gh)
	order.RegisterOrderServiceServer(grpcServer, h)
	// courier.RegisterCourierServiceServer(grpcServer, gh)
}

func catchPanic() {
	if p := recover(); p != nil {
		log.Println("recovered from panic", p)
	}
}

// FindAllNearbyOrderIDs ...
// TODO: build adapters seperately
func (h *GrpcHandler) FindAllNearbyOrderIDs(ctx context.Context, req *order.FindAllNearbyOrderIDsRequest) (*order.FindAllNearbyOrderIDsReply, error) {
	defer catchPanic()

	res := &order.FindAllNearbyOrderIDsReply{}
	coord := req.GetLocation()
	radius := req.GetRadius()

	coordo := internal.NewLocation(coord.GetLon(), coord.GetLat())

	orderIDs, err := h.service.FindAllNearbyOrderIDs(&coordo, radius)

	// TODO: Move rest status updates to central error handler
	if err != nil {
		res.Status = 400
		log.Println("ERROR MAN: ", err)
	} else {
		res.Status = 200
		res.OrderIds = orderIDs
	}

	return res, nil
}

// func (h *GrpcHandler) CheckHealth(ctx context.Context, req *health.CheckHealthRequest) (*health.CheckHealthResponse, error) {
// 	return &health.CheckHealthResponse{
// 		ServiceStatus: 200,
// 	}, nil
// }

// func (h *GrpcHandler) GetCourier(ctx context.Context, req *courier.GetCourierRequest) (*courier.GetCourierResponse, error) {
// 	res := &courier.GetCourierResponse{}

// 	c, err := service.GetCourier(req.Id)
// 	if err != nil {
// 		res.Status = 400
// 	} else {
// 		res.Status = 200
// 	}

// 	return nil, res
// }

// func (h *GrpcHandler) DeleteCourier(ctx context.Context, req *courier.DeleteCourierRequest) (*courier.DeleteCourierResponse, error) {
// 	res := &courier.GetCourierResponse{}

// 	if err := service.DeleteCourier(req.CourierID); err != nil {
// 		res.Status = 500
// 	} else {
// 		res.Status = 200
// 	}

// 	res.Status = 200
// 	return res, nil
// }

// func GetAllNearbyCouriers(context.Context, req *courier.GetAllNearbyCouriersRequest) (*courier.GetAllNearbyCouriersResponse, error) {
//   res := &courier.GetAllNearbyCouriersResponse{}

// 	couriers, err != service.GetAllNearbyCouriers(coord, radius)

// 	if err != nil {
// 		res.Status = 400
// 		res.Couriers = []
// 	} else {
// 		res.Status = 200
// 		res.Couriers = couriers
// 	}

// 	return res, nil
// }
