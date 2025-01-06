package service

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kraftzpepe/auth-service/internal/models"
	"github.com/kraftzpepe/auth-service/internal/repositories"
	"github.com/kraftzpepe/auth-service/internal/utils"
)

type AuthService struct {
	UserRepo               *repositories.UserRepository
	RefreshTokenRepo       *repositories.RefreshTokenRepository
	PasswordResetTokenRepo *repositories.PasswordResetTokenRepository
}

func NewAuthService(
	userRepo *repositories.UserRepository,
	refreshTokenRepo *repositories.RefreshTokenRepository,
	passwordResetTokenRepo *repositories.PasswordResetTokenRepository,
) *AuthService {
	return &AuthService{
		UserRepo:               userRepo,
		RefreshTokenRepo:       refreshTokenRepo,
		PasswordResetTokenRepo: passwordResetTokenRepo,
	}
}

// Signup creates a new user, generates tokens, and saves them in the database
func (s *AuthService) Signup(ctx context.Context, username, email, password string) (*models.User, string, string, error) {
	// Validate inputs
	if err := utils.ValidateEmail(email); err != nil {
		return nil, "", "", err
	}
	if err := utils.ValidatePassword(password); err != nil {
		return nil, "", "", err
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, "", "", err
	}

	// Initialize user
	now := time.Now()
	user := &models.User{
		ID:        uuid.New(),
		Username:  username,
		Email:     email,
		Password:  hashedPassword,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Save user to the database
	err = s.UserRepo.CreateUser(ctx, user)
	if err != nil {
		if strings.Contains(err.Error(), "users_username_key") {
			return nil, "", "", errors.New("username is already taken")
		}
		if strings.Contains(err.Error(), "users_email_key") {
			return nil, "", "", errors.New("email is already registered")
		}
		return nil, "", "", errors.New("failed to save user")
	}

	// Generate tokens
	accessToken, err := utils.GenerateJWT(user.ID.String())
	if err != nil {
		return nil, "", "", errors.New("failed to generate access token")
	}

	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return nil, "", "", errors.New("failed to generate refresh token")
	}

	// Save refresh token to the database
	err = s.RefreshTokenRepo.SaveRefreshToken(user.ID, refreshToken, time.Now().Add(7*24*time.Hour)) // 7 days expiration
	if err != nil {
		return nil, "", "", errors.New("failed to save refresh token")
	}

	return user, accessToken, refreshToken, nil
}

// Login authenticates a user and issues tokens
func (s *AuthService) Login(ctx context.Context, email, password string) (string, string, error) {
	user, err := s.UserRepo.GetUserByEmail(ctx, email)
	if err != nil || user == nil {
		return "", "", errors.New("invalid email or password")
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return "", "", errors.New("invalid email or password")
	}

	accessToken, err := utils.GenerateJWT(user.ID.String())
	if err != nil {
		return "", "", errors.New("failed to generate access token")
	}

	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return "", "", errors.New("failed to generate refresh token")
	}

	err = s.RefreshTokenRepo.SaveRefreshToken(user.ID, refreshToken, time.Now().Add(7*24*time.Hour))
	if err != nil {
		return "", "", errors.New("failed to save refresh token")
	}

	return accessToken, refreshToken, nil
}

// RequestPasswordReset generates a reset token and sends it to the user's email
func (s *AuthService) RequestPasswordReset(ctx context.Context, email string) (string, error) {
	// Fetch the user by email
	user, err := s.UserRepo.GetUserByEmail(ctx, email)
	if err != nil || user == nil {
		return "", errors.New("user not found")
	}

	// Generate a reset token
	resetToken, err := utils.GenerateRefreshToken() // Reuse the token generation logic
	if err != nil {
		return "", errors.New("failed to generate reset token")
	}

	// Set expiration time (e.g., 15 minutes)
	expiresAt := time.Now().Add(15 * time.Minute)

	// Save the token to the database
	err = s.PasswordResetTokenRepo.SaveToken(user.ID, resetToken, expiresAt)
	if err != nil {
		return "", errors.New("failed to save reset token")
	}

	// Send the reset token via email
	err = utils.SendPasswordResetEmail(email, resetToken)
	if err != nil {
		return "", errors.New("failed to send password reset email")
	}

	return "Password reset email sent successfully.", nil
}

// ResetPassword verifies the reset token and updates the user's password
func (s *AuthService) ResetPassword(ctx context.Context, token, newPassword string) (string, error) {
	resetToken, err := s.PasswordResetTokenRepo.FindToken(token)
	if err != nil || resetToken == nil {
		log.Printf("Token not found or error: %v", err)
		return "", errors.New("invalid or expired token")
	}

	log.Printf("Reset Token: %+v", resetToken)

	if resetToken.ExpiresAt.Before(time.Now()) {
		log.Printf("Token expired at: %s", resetToken.ExpiresAt)
		return "", errors.New("invalid or expired token")
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return "", errors.New("failed to hash password")
	}

	err = s.UserRepo.UpdatePassword(ctx, resetToken.UserID.String(), hashedPassword)
	if err != nil {
		return "", errors.New("failed to update password")
	}

	err = s.PasswordResetTokenRepo.DeleteToken(token)
	if err != nil {
		return "", errors.New("failed to delete reset token")
	}

	return "Password has been reset successfully.", nil
}

// RefreshAccessToken generates a new AccessToken and RefreshToken
func (s *AuthService) RefreshAccessToken(refreshToken string) (string, string, error) {
	tokenData, err := s.RefreshTokenRepo.FindRefreshToken(refreshToken)
	if err != nil || tokenData.ExpiresAt.Before(time.Now()) {
		return "", "", errors.New("invalid or expired refresh token")
	}

	accessToken, err := utils.GenerateJWT(tokenData.UserID.String())
	if err != nil {
		return "", "", errors.New("failed to generate access token")
	}

	newRefreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return "", "", errors.New("failed to generate refresh token")
	}

	err = s.RefreshTokenRepo.UpdateRefreshToken(tokenData.UserID, newRefreshToken, time.Now().Add(7*24*time.Hour))
	if err != nil {
		return "", "", errors.New("failed to update refresh token")
	}

	return accessToken, newRefreshToken, nil
}

// GetUserByEmail retrieves a user by their email
func (s *AuthService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := s.UserRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// GetUserByUUID retrieves a user by their UUID
func (s *AuthService) GetUserByUUID(ctx context.Context, uuid string) (*models.User, error) {
	user, err := s.UserRepo.GetUserByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// GetUserByUsername retrieves a user by their username
func (s *AuthService) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	user, err := s.UserRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}
