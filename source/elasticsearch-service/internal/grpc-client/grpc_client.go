package grpc_client

import (
	"fmt"
	"thanhldt060802/internal/grpc-client/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClientConfig struct {
	UserServiceConn *grpc.ClientConn
}

func NewGRPCClientConfig(userServiceHost string) (*GRPCClientConfig, pb.UserServiceClient, error) {
	client, conn, err := NewUserGRPCClient(userServiceHost)
	if err != nil {
		return nil, nil, err
	}
	return &GRPCClientConfig{
		UserServiceConn: conn,
	}, client, nil
}

func NewUserGRPCClient(target string) (pb.UserServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to gRPC server: %w", err)
	}

	client := pb.NewUserServiceClient(conn)
	return client, conn, nil
}

func (c *GRPCClientConfig) Close() {
	if c.UserServiceConn != nil {
		c.UserServiceConn.Close()
	}
}
