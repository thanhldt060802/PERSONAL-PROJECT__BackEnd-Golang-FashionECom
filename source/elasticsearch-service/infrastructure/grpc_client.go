package infrastructure

import (
	"fmt"
	"log"
	"thanhldt060802/config"
	"thanhldt060802/internal/grpc/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var userServiceGRPCConnection *grpc.ClientConn
var UserServiceGRPCClient pb.UserServiceGRPCClient

func CloseAll() {
	userServiceGRPCConnection.Close()
}

func InitAllServiceGRPCClients() {
	// Kết nối user-service
	userServiceGRPCServerAddress := fmt.Sprintf(
		"%s:%s",
		config.AppConfig.UserServiceGRPCServerHost,
		config.AppConfig.UserServiceGRPCServerPort,
	)

	conn, err := grpc.NewClient(userServiceGRPCServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("connect to user-service failed: %s", err.Error())
	}
	userServiceGRPCConnection = conn
	UserServiceGRPCClient = pb.NewUserServiceGRPCClient(conn)

	log.Println("connect to user-service successful")

	// Kết nối catalog-service

	// Kết nối order-service
}
