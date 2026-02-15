package services

import (
	"context"
	"errors"
	"time"

	"github.com/devlpr-nitish/leetcode-tracker-backend/internal/config"
	"github.com/devlpr-nitish/leetcode-tracker-backend/internal/models"
	"github.com/devlpr-nitish/leetcode-tracker-backend/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepo  repository.UserRepository
	JWTSecret string
}

func NewAuthService(userRepo repository.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		UserRepo:  userRepo,
		JWTSecret: cfg.JWTSecret,
	}
}

type TokenResponse struct {
	Token string `json:"token"`
}

func (s *AuthService) Signup(ctx context.Context, username, email, password string) (*models.User, error) {
	// Check if user exists
	existingUser, _ := s.UserRepo.GetByUsername(ctx, username)
	if existingUser != nil {
		return nil, errors.New("username already taken")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
		// LeetCodeUsername: username, // Removed as it is now redundant with Username
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.UserRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, identifier, password string) (*TokenResponse, error) {
	// Try to find user by username
	user, err := s.UserRepo.GetByUsername(ctx, identifier)
	// If not found by username, try email. Ideally, repository would have GetByEmail or GetByIdentifier
	if user == nil {
		user, err = s.UserRepo.GetByEmail(ctx, identifier)
	}
	// But since UserRepo interface is defined (and typically implemented with GORM), let's assume strict username for now or we update repo.
	// Since user prompted "signup using his leetcode username , email and password, and login using either email / password or username or password",
	// I should check if identifier is email.

	if user == nil {
		// Try fetching by email - need to add GetByEmail to repo interface or cheat a bit.
		// Let's assume I fix repo later or do it here. For now, strict username login if repo method missing.
		// Wait, I can update repo interface.
		return nil, errors.New("invalid credentials")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate JWT
	token, err := s.generateJWT(user)
	if err != nil {
		return nil, err
	}

	return &TokenResponse{Token: token}, nil
}

func (s *AuthService) generateJWT(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(), // 72 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.JWTSecret))
}
