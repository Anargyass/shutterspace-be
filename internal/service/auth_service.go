package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"shutterspace/internal/domain"
	"shutterspace/internal/repository"
	"shutterspace/pkg/hash"
	pkgjwt "shutterspace/pkg/jwt"
)

type AuthService interface {
	Register(ctx context.Context, name, email, password string) (*domain.User, error)
	Login(ctx context.Context, email, password string) (*domain.User, string, error)
}

type authService struct {
	userRepo    repository.UserRepository
	jwtSecret   string
	expiryHours int
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string, expiryHours int) AuthService {
	return &authService{
		userRepo:    userRepo,
		jwtSecret:   jwtSecret,
		expiryHours: expiryHours,
	}
}

func (s *authService) Register(ctx context.Context, name, email, password string) (*domain.User, error) {
	// Cek apakah email sudah terdaftar
	_, err := s.userRepo.FindByEmail(ctx, email)
	if err == nil {
		// Ditemukan → email sudah ada
		return nil, domain.ErrEmailExists
	}
	if !errors.Is(err, domain.ErrNotFound) {
		return nil, err
	}

	// Hash password sebelum disimpan
	passwordHash, err := hash.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		ID:           uuid.New(),
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         domain.RoleUser,
		IsActive:     true,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (*domain.User, string, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		// SECURITY: Jangan beri tahu apakah email yang salah atau password yang salah
		return nil, "", domain.ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, "", domain.ErrInactive
	}

	if err := hash.CheckPassword(password, user.PasswordHash); err != nil {
		return nil, "", domain.ErrInvalidCredentials
	}

	token, err := pkgjwt.SignToken(user.ID, user.Email, string(user.Role), s.jwtSecret, s.expiryHours)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}
