package grpcimpl

import (
	"fmt"
	"log"
	"net"
	"thanhldt060802/config"
	"thanhldt060802/internal/grpc/service/elasticsearchservicepb"

	"google.golang.org/grpc"
)

func StartGRPCServer(elasticsearchServiceGRPCImpl elasticsearchservicepb.ElasticsearchServiceGRPCServer) {
	address := fmt.Sprintf(
		"%s:%s",
		config.AppConfig.ElasticsearchServiceGRPCHost,
		config.AppConfig.ElasticsearchServiceGRPCPort,
	)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Start gRPC server failed: %s", err.Error())
	}

	server := grpc.NewServer()

	elasticsearchservicepb.RegisterElasticsearchServiceGRPCServer(server, elasticsearchServiceGRPCImpl)

	log.Printf("Start gRPC server successful")

	go func() {
		if err := server.Serve(lis); err != nil {
			log.Fatalf("Serve failed: %s", err.Error())
		}
	}()
}
