package grpc_server

import (
	"context"
	"fmt"
	"log"
	"net"
	"thanhldt060802/internal/dto"
	"thanhldt060802/internal/grpc-server/pb"
	"thanhldt060802/internal/repository"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	pb.UnimplementedUserServiceServer
	userRepository repository.UserRepository
}

func NewServer(userRepository repository.UserRepository) *Server {
	return &Server{
		userRepository: userRepository,
	}
}

func (s *Server) GetAllUsers(ctx context.Context, req *pb.GetAllUsersRequest) (*pb.GetAllUsersResponse, error) {
	// // Tạo dữ liệu mẫu
	// users := []*pb.User{
	// 	{
	// 		Id:        1,
	// 		FullName:  "John Doe",
	// 		Email:     "john@example.com",
	// 		Username:  "johndoe",
	// 		Address:   "123 Main St",
	// 		RoleName:  "admin",
	// 		CreatedAt: "2025-01-01",
	// 		UpdatedAt: "2025-01-01",
	// 	},
	// 	{
	// 		Id:        2,
	// 		FullName:  "Jane Smith",
	// 		Email:     "jane@example.com",
	// 		Username:  "janesmith",
	// 		Address:   "456 Oak St",
	// 		RoleName:  "user",
	// 		CreatedAt: "2025-02-01",
	// 		UpdatedAt: "2025-02-01",
	// 	},
	// }

	if s.userRepository == nil {
		fmt.Println("asdasdasdasdasd")
	}

	users, err := s.userRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	// Trả về danh sách người dùng
	return &pb.GetAllUsersResponse{
		Users: dto.ToListUserProto(users),
	}, nil
}

func (s *Server) Run() {
	// Lắng nghe trên port 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Tạo gRPC server
	s1 := grpc.NewServer(grpc.Creds(insecure.NewCredentials()))

	// Đăng ký dịch vụ UserService
	pb.RegisterUserServiceServer(s1, &Server{})

	// Bắt đầu server gRPC
	fmt.Println("gRPC server started on port 50051")
	if err := s1.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
