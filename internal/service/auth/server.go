package auth

import (
	"context"
	pb "github.com/EugeneNail/actum/grpc/gen/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	authentication AuthenticationService
	pb.UnimplementedAuthServer
}

func NewServer(authentication AuthenticationService) *Server {
	return &Server{authentication: authentication}
}

func (s *Server) Login(context.Context, *pb.LoginRequest) (*pb.LoginResponse, error) {
	return &pb.LoginResponse{}, nil
}

func (s *Server) Register(ctx context.Context, r *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	id, err := s.authentication.Register(ctx, r.Name, r.Email, r.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.RegisterResponse{Id: id}, nil
}

func (s *Server) Logout(context.Context, *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	return &pb.LogoutResponse{}, nil
}
