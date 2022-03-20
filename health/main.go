package main

import (
	"log"
	"os"

	"github.com/chiru1221/microservice-knapsack/client"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	client := client.ClientImpl{
		Addr: os.Getenv("SERVER_ADDR"),
	}
	if err := client.Construct(); err != nil {
		log.Fatal(err)
	}
	resp, err := client.CallHealth("healthz")
	if err != nil {
		log.Fatal(err)
	}
	if resp.Status != healthpb.HealthCheckResponse_SERVING {
		log.Fatal(resp)
	}
}
