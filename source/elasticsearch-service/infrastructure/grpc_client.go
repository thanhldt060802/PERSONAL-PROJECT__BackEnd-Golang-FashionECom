package infrastructure

import (
	"fmt"
	"log"
	"thanhldt060802/config"
	"thanhldt060802/internal/grpc/client/userservicepb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var UserServiceGRPCClient userservicepb.UserServiceGRPCClient

type serviceGRPCConnectionManager struct {
	userServiceGRPCConnection *grpc.ClientConn
}

func (serviceGRPCConnectionManager *serviceGRPCConnectionManager) CloseAll() {
	serviceGRPCConnectionManager.userServiceGRPCConnection.Close()
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

	conn, err := grpc.NewClient(userServiceGRPCServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("connect to user-service failed: %s", err.Error())
	}
	ServiceGRPCConnectionManager.userServiceGRPCConnection = conn
	UserServiceGRPCClient = userservicepb.NewUserServiceGRPCClient(conn)

	log.Printf("connect to user-service successful")

	// Kết nối catalog-service

	// Kết nối order-service
}
