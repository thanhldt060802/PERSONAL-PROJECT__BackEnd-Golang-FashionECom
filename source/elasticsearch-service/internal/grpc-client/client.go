package grpc_client

import (
	"context"
	"fmt"
	"log"
	"thanhldt060802/internal/dto"
	"thanhldt060802/internal/grpc-client/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Run() {
	// Kết nối tới server gRPC bằng NewClient (thay thế grpc.Dial)
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Tạo client gRPC
	c := pb.NewUserServiceClient(conn)

	// Gọi hàm GetAllUsers
	resp, err := c.GetAllUsers(context.Background(), &pb.GetAllUsersRequest{})
	if err != nil {
		log.Fatalf("could not fetch users: %v", err)
	}

	// In danh sách người dùng
	fmt.Println("Users:")
	for _, u := range resp.Users {
		fmt.Printf("ID: %d, Name: %s, Email: %s\n", u.GetId(), u.GetFullName(), u.GetEmail())
	}
}

func GetAllUsers() []dto.UserView {
	// Kết nối tới server gRPC bằng NewClient (thay thế grpc.Dial)
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Tạo client gRPC
	c := pb.NewUserServiceClient(conn)

	// Gọi hàm GetAllUsers
	resp, err := c.GetAllUsers(context.Background(), &pb.GetAllUsersRequest{})
	if err != nil {
		log.Fatalf("could not fetch users: %v", err)
	}

	// In danh sách người dùng
	return dto.ToListUserViewFromProto(resp.Users)
}
