package grpcimpl

import (
	"fmt"
	"log"
	"net"
	"thanhldt060802/config"
	"thanhldt060802/internal/grpc/service/userservicepb"

	"google.golang.org/grpc"
)

func StartGRPCServer(userServiceGRPCImpl userservicepb.UserServiceGRPCServer) {
	address := fmt.Sprintf(
		"%s:%s",
		config.AppConfig.UserServiceGRPCHost,
		config.AppConfig.UserServiceGRPCPort,
	)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Start gRPC server failed: %s", err.Error())
	}

	server := grpc.NewServer()

	userservicepb.RegisterUserServiceGRPCServer(server, userServiceGRPCImpl)

	log.Printf("Start gRPC server successful")

	go func() {
		if err := server.Serve(lis); err != nil {
			log.Fatalf("Serve failed: %s", err.Error())
		}
	}()
}
