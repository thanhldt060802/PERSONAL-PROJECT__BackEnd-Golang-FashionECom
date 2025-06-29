package grpcimpl

import (
	"fmt"
	"log"
	"net"
	"thanhldt060802/config"
	"thanhldt060802/internal/grpc/service/orderservicepb"

	"google.golang.org/grpc"
)

func StartGRPCServer(orderServiceGRPCImpl orderservicepb.OrderServiceGRPCServer) {
	address := fmt.Sprintf(
		"%s:%s",
		config.AppConfig.OrderServiceGRPCHost,
		config.AppConfig.OrderServiceGRPCPort,
	)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Start gRPC server failed: %s", err.Error())
	}

	server := grpc.NewServer()

	orderservicepb.RegisterOrderServiceGRPCServer(server, orderServiceGRPCImpl)

	log.Printf("Start gRPC server successful")

	go func() {
		if err := server.Serve(lis); err != nil {
			log.Fatalf("Serve failed: %s", err.Error())
		}
	}()
}
