package infrastructure

import (
	"log"
	"net"
	"thanhldt060802/internal/grpc/pb"

	"google.golang.org/grpc"
)

func StartGRPCServer(address string, productServiceImpl pb.ProductServiceServer) {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("gRPC server start failed: %s", err.Error())
	}

	server := grpc.NewServer()

	pb.RegisterProductServiceServer(server, productServiceImpl)

	log.Printf("gRPC server is starting on %s\n", address)

	if err := server.Serve(lis); err != nil {
		log.Fatalf("Serve failed: %s", err.Error())
	}
}
