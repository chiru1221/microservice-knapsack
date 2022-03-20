package client

import (
	"context"
	"time"

	pb "github.com/chiru1221/microservice-knapsack/knapsack"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type Client interface {
	Construct() error
	Destruct()
	CallUpdate(knapsackDP *pb.KnapsackDP) (*pb.Result, error)
	CallHealth(service string) (*healthpb.HealthCheckResponse, error)
}

type ClientImpl struct {
	Addr         string
	Conn         *grpc.ClientConn
	grpcClient   pb.UpdateTableClient
	healthClient healthpb.HealthClient
}

func (client *ClientImpl) Construct() error {
	conn, err := grpc.Dial(client.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	client.Conn = conn
	client.grpcClient = pb.NewUpdateTableClient(client.Conn)
	client.healthClient = healthpb.NewHealthClient(client.Conn)
	return nil
}

func (client *ClientImpl) Destruct() {
	client.Conn.Close()
}

func (client *ClientImpl) CallUpdate(knapsackDP *pb.KnapsackDP) (*pb.Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return client.grpcClient.Update(ctx, knapsackDP)
}

func (client *ClientImpl) CallHealth(service string) (*healthpb.HealthCheckResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return client.healthClient.Check(ctx, &healthpb.HealthCheckRequest{Service: service})
}
