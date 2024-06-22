package auth

import (
	"context"
	pb "github.com/EugeneNail/actum/grpc/gen/auth"
)

type Client interface {
	Register(ctx context.Context, name string, email string, password string) (id int64, err error)
	Login(ctx context.Context, email string, password string) (accessToken string, refreshToken string, err error)
}

type client struct {
	stub pb.AuthClient
}

func NewClient(grpcClient pb.AuthClient) Client {
	return &client{grpcClient}
}

func (ac *client) Register(ctx context.Context, name string, email string, password string) (id int64, err error) {
	resp, err := ac.stub.Register(ctx, &pb.RegisterRequest{
		Name:     name,
		Email:    email,
		Password: password,
	})

	if err != nil {
		return 0, err
	}

	return resp.Id, nil
}

func (ac *client) Login(ctx context.Context, email string, password string) (accessToken string, refreshToken string, err error) {
	resp, err := ac.stub.Login(ctx, &pb.LoginRequest{
		Email:    email,
		Password: password,
	})

	if err != nil {
		return "", "", err
	}

	return resp.AccessToken, resp.RefreshToken, nil
}
