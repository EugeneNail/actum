package auth

import (
	"context"
	"github.com/EugeneNail/actum/internal/database/repository"
	"github.com/EugeneNail/actum/internal/infrastructure/errors"
	"github.com/EugeneNail/actum/internal/infrastructure/hash"
	"github.com/EugeneNail/actum/internal/service/auth/jwt"
	"github.com/EugeneNail/actum/internal/service/auth/refresh"
	"strings"
)

type AuthenticationService interface {
	Login(cxt context.Context, email string, password string) (accessToken string, refreshToken string, err error)
	Register(cxt context.Context, name string, email string, password string) (id int64, err error)
}

type service struct {
	repo    repository.UserRepository
	refresh *refresh.Service
}

func NewAuthenticationService(repo repository.UserRepository, refresh *refresh.Service) AuthenticationService {
	return &service{repo, refresh}
}

func (s *service) Login(ctx context.Context, email string, password string) (string, string, error) {
	user, err := s.repo.Find(ctx, "email", strings.ToLower(email))
	if err != nil {
		return "", "", errors.Wrap(err, "failed to find the user")
	}

	if user.Id == 0 || user.Password != hash.New(password) {
		return "", "", nil
	}

	accessToken, err := jwt.Make(user.Id)
	if err != nil {
		return "", "", errors.Wrap(err, "failed to create jwt token")
	}

	refreshToken, err := s.refresh.MakeToken(user.Id)
	if err != nil {
		return "", "", errors.Wrap(err, "failed to create refresh token")
	}

	return accessToken, refreshToken, nil
}

func (s *service) Register(ctx context.Context, name string, email string, password string) (int64, error) {
	id, err := s.repo.Create(ctx, name, email, password)
	if err != nil {
		return 0, errors.Wrap(err, "failed to create the user")
	}

	return id, nil
}
