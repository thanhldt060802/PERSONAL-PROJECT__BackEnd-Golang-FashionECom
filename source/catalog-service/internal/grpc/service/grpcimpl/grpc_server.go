package grpcimpl

import (
	"fmt"
	"log"
	"net"
	"thanhldt060802/config"
	"thanhldt060802/internal/grpc/service/catalogservicepb"

	"google.golang.org/grpc"
)

func StartGRPCServer(catalogServiceGRPCImpl catalogservicepb.CatalogServiceGRPCServer) {
	address := fmt.Sprintf(
		"%s:%s",
		config.AppConfig.CatalogServiceGRPCHost,
		config.AppConfig.CatalogServiceGRPCPort,
	)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Start gRPC server failed: %s", err.Error())
	}

	server := grpc.NewServer()

	catalogservicepb.RegisterCatalogServiceGRPCServer(server, catalogServiceGRPCImpl)

	log.Printf("Start gRPC server successful")

	go func() {
		if err := server.Serve(lis); err != nil {
			log.Fatalf("Serve failed: %s", err.Error())
		}
	}()
}
