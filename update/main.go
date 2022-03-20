package main

import (
	"context"
	"log"
	"net"

	pb "github.com/chiru1221/microservice-knapsack/knapsack"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type serverImpl struct {
	pb.UnimplementedUpdateTableServer
}

func (server *serverImpl) Update(ctx context.Context, knapsackDP *pb.KnapsackDP) (*pb.Result, error) {
	var result int32 = knapsackDP.Dp[knapsackDP.I+1].Row[knapsackDP.W]
	if knapsackDP.W-knapsackDP.Weight >= 0 {
		if knapsackDP.Dp[knapsackDP.I+1].Row[knapsackDP.W] <
			knapsackDP.Dp[knapsackDP.I].Row[knapsackDP.W-knapsackDP.Weight]+knapsackDP.Value {
			result = knapsackDP.Dp[knapsackDP.I].Row[knapsackDP.W-knapsackDP.Weight] + knapsackDP.Value
		}
	}
	if result < knapsackDP.Dp[knapsackDP.I].Row[knapsackDP.W] {
		result = knapsackDP.Dp[knapsackDP.I].Row[knapsackDP.W]
	}
	return &pb.Result{Result: result}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	pb.RegisterUpdateTableServer(s, &serverImpl{})
	srv := health.NewServer()
	srv.SetServingStatus("healthz", healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(s, srv)
	log.Println("Start server")
	if err = s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
