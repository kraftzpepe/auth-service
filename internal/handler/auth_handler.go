package handler

import (
	"context"

	"github.com/kraftzpepe/auth-service/internal/service"
	pb "github.com/kraftzpepe/auth-service/proto/generated"
)

type AuthHandler struct {
	pb.UnimplementedAuthServiceServer
	AuthService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{AuthService: authService}
}

// Implement the Register method
func (h *AuthHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	user, accessToken, refreshToken, err := h.AuthService.Signup(ctx, req.GetUsername(), req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, err
	}

	return &pb.RegisterResponse{
		Id:           user.ID.String(),
		Username:     user.Username,
		Email:        user.Email,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// Implement the Login method
func (h *AuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	accessToken, refreshToken, err := h.AuthService.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// gRPC endpoint for refreshing access tokens
func (h *AuthHandler) RefreshAccessToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	accessToken, refreshToken, err := h.AuthService.RefreshAccessToken(req.GetRefreshToken())
	if err != nil {
		return nil, err
	}

	return &pb.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (h *AuthHandler) RequestPasswordReset(ctx context.Context, req *pb.RequestPasswordResetRequest) (*pb.RequestPasswordResetResponse, error) {
	message, err := h.AuthService.RequestPasswordReset(ctx, req.GetEmail())
	if err != nil {
		return nil, err
	}

	return &pb.RequestPasswordResetResponse{
		Message: message,
	}, nil
}

func (h *AuthHandler) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	message, err := h.AuthService.ResetPassword(ctx, req.GetToken(), req.GetNewPassword())
	if err != nil {
		return nil, err
	}

	return &pb.ResetPasswordResponse{
		Message: message,
	}, nil
}

// gRPC endpoint for getting a user by email
func (h *AuthHandler) GetUserByEmail(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := h.AuthService.GetUserByEmail(ctx, req.GetIdentifier())
	if err != nil {
		return nil, err
	}

	return &pb.GetUserResponse{
		Id:        user.ID.String(),
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	}, nil
}

// gRPC endpoint for getting a user by UUID
func (h *AuthHandler) GetUserByUUID(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := h.AuthService.GetUserByUUID(ctx, req.GetIdentifier())
	if err != nil {
		return nil, err
	}

	return &pb.GetUserResponse{
		Id:        user.ID.String(),
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	}, nil
}

// gRPC endpoint for getting a user by username
func (h *AuthHandler) GetUserByUsername(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := h.AuthService.GetUserByUsername(ctx, req.GetIdentifier())
	if err != nil {
		return nil, err
	}

	return &pb.GetUserResponse{
		Id:        user.ID.String(),
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	}, nil
}
