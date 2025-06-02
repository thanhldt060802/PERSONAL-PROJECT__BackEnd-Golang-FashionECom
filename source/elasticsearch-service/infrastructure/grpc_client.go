package infrastructure

import (
	"fmt"
	"log"
	"thanhldt060802/config"
	"thanhldt060802/internal/grpc/client/catalogservicepb"
	"thanhldt060802/internal/grpc/client/userservicepb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var UserServiceGRPCClient userservicepb.UserServiceGRPCClient
var CatalogServiceGRPCClient catalogservicepb.CatalogServiceGRPCClient

type serviceGRPCConnectionManager struct {
	userServiceGRPCConnection    *grpc.ClientConn
	catalogServiceGRPCConnection *grpc.ClientConn
}

func (serviceGRPCConnectionManager *serviceGRPCConnectionManager) CloseAll() {
	serviceGRPCConnectionManager.userServiceGRPCConnection.Close()
	serviceGRPCConnectionManager.catalogServiceGRPCConnection.Close()
}

var ServiceGRPCConnectionManager *serviceGRPCConnectionManager

func InitAllServiceGRPCClients() {
	ServiceGRPCConnectionManager = &serviceGRPCConnectionManager{}

	// Kết nối user-service
	userServiceGRPCServerAddress := fmt.Sprintf(
		"%s:%s",
		config.AppConfig.UserServiceGRPCHost,
		config.AppConfig.UserServiceGRPCPort,
	)

	conn1, err := grpc.NewClient(userServiceGRPCServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("connect to user-service failed: %s", err.Error())
	}
	ServiceGRPCConnectionManager.userServiceGRPCConnection = conn1
	UserServiceGRPCClient = userservicepb.NewUserServiceGRPCClient(conn1)

	log.Printf("connect to user-service successful")

	// Kết nối catalog-service
	catalogServiceGRPCServerAddress := fmt.Sprintf(
		"%s:%s",
		config.AppConfig.CatalogServiceGRPCHost,
		config.AppConfig.CatalogServiceGRPCPort,
	)

	conn2, err := grpc.NewClient(catalogServiceGRPCServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("connect to catalog-service failed: %s", err.Error())
	}
	ServiceGRPCConnectionManager.catalogServiceGRPCConnection = conn2
	CatalogServiceGRPCClient = catalogservicepb.NewCatalogServiceGRPCClient(conn2)

	log.Printf("connect to catalog-service successful")

	// Kết nối order-service
}
