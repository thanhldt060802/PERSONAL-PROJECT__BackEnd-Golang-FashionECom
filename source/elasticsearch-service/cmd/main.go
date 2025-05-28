package main

import (
	"thanhldt060802/config"
	"thanhldt060802/infrastructure"
	"thanhldt060802/internal/grpc/service/grpcimpl"
	"thanhldt060802/internal/service"
)

func main() {

	config.InitConfig()
	infrastructure.InitElasticsearchClient()
	infrastructure.InitRedisClient()
	defer infrastructure.RedisClient.Close()
	infrastructure.InitAllServiceGRPCClients()
	defer infrastructure.ServiceGRPCConnectionManager.CloseAll()

	grpcimpl.StartGRPCServer(grpcimpl.NewElasticsearchServiceGRPCImpl(service.NewUserService()))

	select {}

}
