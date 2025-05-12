package grpc_server

import (
	"fmt"
	"log"
	"net"
	"thanhldt060802/internal/grpc-server/pb"

	"google.golang.org/grpc"
)

type GRPCServerConfig struct {
	Port            int
	GetAllUsersImpl pb.UserServiceServer
}

func StartGRPCServer(cfg GRPCServerConfig) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalf("gRPC server start failed: %s", err.Error())
	}

	server := grpc.NewServer()

	pb.RegisterUserServiceServer(server, cfg.GetAllUsersImpl)

	log.Printf("gRPC server is starting on port %d\n", cfg.Port)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Serve failed: %s", err.Error())
	}
}
