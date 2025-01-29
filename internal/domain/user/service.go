package user

import (
	"context"
	"fmt"

	"github.com/Gabukuro/gymratz-api/internal/infra/ports/repo"
	"github.com/Gabukuro/gymratz-api/internal/pkg/entity/user"
	"github.com/Gabukuro/gymratz-api/internal/pkg/jwt"
)

type (
	Service struct {
		userRepo     repo.UserRepository
		tokenService *jwt.TokenService
	}

	ServiceParams struct {
		UserRepo     repo.UserRepository
		TokenService *jwt.TokenService
	}
)

func NewService(params ServiceParams) *Service {
	return &Service{
		userRepo:     params.UserRepo,
		tokenService: params.TokenService,
	}
}

func (s *Service) CreateUser(ctx context.Context, name, email, password string) error {
	userMode := user.Model{
		Name:     name,
		Email:    email,
		Password: password,
	}

	if err := userMode.HashPassword(); err != nil {
		return fmt.Errorf("could not hash password: %w", err)
	}

	return s.userRepo.Create(ctx, userMode)
}

func (s *Service) LoginUser(ctx context.Context, email, password string) (*string, error) {
	userModel, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("could not find user: %w", err)
	}

	if !userModel.CheckPassword(password) {
		return nil, fmt.Errorf("invalid password")
	}

	token, err := s.tokenService.GenerateToken(userModel.Email)
	if err != nil {
		return nil, fmt.Errorf("could not generate token: %w", err)
	}

	return &token, nil
}
