package main

import (
	// "log"

	"log"

	"google.golang.org/grpc"
	// "google.golang.org/grpc"
	// "location-service/api/proto"
)

func main() {
	var conn *grpc.ClientConn

	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("error time! %s", err)
	}

	defer conn.Close()

	// client := proto.NewCourierServiceClient(conn)
	// msg := &proto.GetCourierRequest{Id: "id"}

	// res, err := client.CheckHealth(context.Background(), msg)
	if err != nil {
		log.Fatalf(err.Error())
	}

	// log.Printf("Location service response: %d", res.ServiceStatus)
}
