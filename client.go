package main

import (
  "log"

  "golang.org/x/net/context"
  "google.golang.org/grpc"

  "location-service/api/proto/driver"
)

func main() {
  var conn *grpc.ClientConn

  conn, err := grpc.Dial(":9000", grpc.WithInsecure())
  if err != nil {
    log.Fatalf("error time! %s", err)
  }

  defer conn.Close()

  client := driver.NewLocationServiceClient(conn)
  msg := &driver.CheckHealthRequest{}

  res, err := client.CheckHealth(context.Background(), msg)
  if err != nil {
    log.Fatalf(err.Error())
  }

  log.Printf("Location service response: %d", res.ServiceStatus)
}
