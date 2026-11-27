package service

import (
	"context"
	"time"

	"github.com/ananda1farelas/sewakos-backend/app/model"
	"github.com/ananda1farelas/sewakos-backend/app/repository"
	"github.com/ananda1farelas/sewakos-backend/helper"
)

type AuthService interface {
	Register(ctx context.Context, req model.RegisterRequest) (*model.AuthResponse, *helper.AppError)
	Login(ctx context.Context, req model.LoginRequest) (*model.AuthResponse, *helper.AppError)
	Refresh(ctx context.Context, refreshToken string) (*model.AuthResponse, *helper.AppError)
	Logout(ctx context.Context, refreshToken string) *helper.AppError
	Me(ctx context.Context, userID int) (*model.UserResponse, *helper.AppError)
}

type authService struct {
	userRepo   repository.UserRepository
	tokenRepo  repository.TokenRepository
	jwtSecret  string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewAuthService(userRepo repository.UserRepository, tokenRepo repository.TokenRepository, jwtSecret string, accessTTL, refreshTTL time.Duration) AuthService {
	return &authService{
		userRepo:   userRepo,
		tokenRepo:  tokenRepo,
		jwtSecret:  jwtSecret,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (s *authService) Register(ctx context.Context, req model.RegisterRequest) (*model.AuthResponse, *helper.AppError) {
	if appErr := ValidateRegisterInput(req.Name, req.Email, req.Password); appErr != nil {
		return nil, appErr
	}

	email := NormalizeEmail(req.Email)

	hashed, err := helper.HashPassword(req.Password)
	if err != nil {
		return nil, helper.ErrInternal("gagal memproses password")
	}

	user := &model.User{
		Name:     req.Name,
		Email:    email,
		Phone:    req.Phone,
		Password: hashed,
		Role:     DefaultRole,
		IsActive: true,
	}

	created, err := s.userRepo.Create(ctx, user)
	if err != nil {
		if err == repository.ErrDuplicate {
			return nil, helper.ErrDuplicate("email sudah terdaftar")
		}
		return nil, helper.ErrInternal("gagal membuat akun")
	}

	return s.issueTokens(ctx, created)
}

func (s *authService) Login(ctx context.Context, req model.LoginRequest) (*model.AuthResponse, *helper.AppError) {
	if appErr := ValidateLoginInput(req.Email, req.Password); appErr != nil {
		return nil, appErr
	}

	email := NormalizeEmail(req.Email)

	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, helper.ErrUnauthorized("email atau password salah")
		}
		return nil, helper.ErrInternal("gagal login")
	}

	if !user.IsActive {
		return nil, helper.ErrForbidden("akun tidak aktif")
	}

	if !helper.CheckPasswordHash(req.Password, user.Password) {
		return nil, helper.ErrUnauthorized("email atau password salah")
	}

	return s.issueTokens(ctx, user)
}

func (s *authService) Refresh(ctx context.Context, refreshToken string) (*model.AuthResponse, *helper.AppError) {
	tokenHash := helper.HashToken(refreshToken)

	stored, err := s.tokenRepo.FindActive(ctx, tokenHash)
	if err != nil {
		return nil, helper.ErrUnauthorized("refresh token tidak valid")
	}

	user, err := s.userRepo.FindByID(ctx, stored.UserID)
	if err != nil {
		return nil, helper.ErrUnauthorized("user tidak ditemukan")
	}

	if !user.IsActive {
		return nil, helper.ErrForbidden("akun tidak aktif")
	}

	// rotasi: cabut token lama, terbitkan baru
	if err := s.tokenRepo.Revoke(ctx, tokenHash); err != nil {
		return nil, helper.ErrInternal("gagal rotasi token")
	}

	return s.issueTokens(ctx, user)
}

func (s *authService) Logout(ctx context.Context, refreshToken string) *helper.AppError {
	tokenHash := helper.HashToken(refreshToken)

	if err := s.tokenRepo.Revoke(ctx, tokenHash); err != nil {
		return helper.ErrInternal("gagal logout")
	}

	return nil
}

func (s *authService) Me(ctx context.Context, userID int) (*model.UserResponse, *helper.AppError) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, helper.ErrNotFound("user tidak ditemukan")
		}
		return nil, helper.ErrInternal("gagal mengambil data user")
	}

	resp := user.ToResponse()
	return &resp, nil
}

func (s *authService) issueTokens(ctx context.Context, user *model.User) (*model.AuthResponse, *helper.AppError) {
	accessToken, err := helper.GenerateAccessToken(user.ID, user.Role, s.jwtSecret, s.accessTTL)
	if err != nil {
		return nil, helper.ErrInternal("gagal membuat access token")
	}

	refreshToken, err := helper.GenerateRandomToken()
	if err != nil {
		return nil, helper.ErrInternal("gagal membuat refresh token")
	}

	tokenHash := helper.HashToken(refreshToken)
	expiresAt := time.Now().Add(s.refreshTTL)

	if err := s.tokenRepo.Store(ctx, user.ID, tokenHash, expiresAt); err != nil {
		return nil, helper.ErrInternal("gagal menyimpan refresh token")
	}

	return &model.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user.ToResponse(),
	}, nil
}
