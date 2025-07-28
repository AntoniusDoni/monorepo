package service

import (
	"errors"

	"github.com/antoniusDoni/monorepo/shared/auth"
	"github.com/antoniusDoni/monorepo/shared/contract"
	"github.com/antoniusDoni/monorepo/shared/model"
	"github.com/antoniusDoni/monorepo/shared/repository"
	"golang.org/x/crypto/bcrypt"
)

// AuthService handles user registration, login, and retrieval
type AuthService struct {
	repo      repository.UserRepository
	jwtSecret string
}

// NewAuthService creates a new AuthService instance
func NewAuthService(repo repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{repo: repo, jwtSecret: jwtSecret}
}

// Register creates a new user with hashed password
func (s *AuthService) Register(username, password string) error {
	_, err := s.repo.FindByUsername(username)
	if err == nil {
		return errors.New("username already exists")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &model.User{
		Username:     username,
		PasswordHash: string(hashed),
	}

	return s.repo.Create(user)
}

// Login validates user credentials and returns a JWT token
func (s *AuthService) Login(username, password string) (*contract.LoginResponse, error) {
	user, err := s.repo.FindByUsername(username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	if len(user.Roles) == 0 {
		return nil, errors.New("user has no roles assigned")
	}

	token, err := auth.CreateJWTToken(user.ID, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &contract.LoginResponse{
		UserID: user.ID,
		Role:   user.Roles[0].Name,
		Token:  token,
	}, nil
}

// GetUser retrieves a user by ID
func (s *AuthService) GetUser(id uint) (*model.User, error) {
	return s.repo.FindByID(id)
}
