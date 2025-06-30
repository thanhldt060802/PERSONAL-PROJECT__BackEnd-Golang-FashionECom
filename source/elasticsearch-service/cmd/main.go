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

	grpcimpl.StartGRPCServer(grpcimpl.NewElasticsearchServiceGRPCImpl(
		service.NewUserService(config.AppConfig.SyncAvailableDataFromUserService),
		service.NewCatalogService(config.AppConfig.SyncAvailableDataFromCatalogService),
		service.NewOrderService(config.AppConfig.SyncAvailableDataFromOrderService),
	))

	select {}

}
