package infrastructure

import (
	"log"
	"thanhldt060802/internal/grpc-client/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClientsManager struct {
	UserServiceClient pb.UserServiceClient

	Connections []*grpc.ClientConn
}

func (m *GRPCClientsManager) Close() {
	for _, conn := range m.Connections {
		_ = conn.Close()
	}
}

func InitGRPCClientsManager(manager *GRPCClientsManager, config map[string]string) {
	// Kết nối đến UserService
	userConn, err := grpc.NewClient(config["user-service"], grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("connect to user-service failed: %s", err.Error())
	}
	manager.UserServiceClient = pb.NewUserServiceClient(userConn)
	manager.Connections = append(manager.Connections, userConn)

	log.Println("connect to all out-services successful")
}
