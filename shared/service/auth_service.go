package service

import (
	"errors"
	"time"

	warehouseModel "github.com/antoniusDoni/monorepo/modules/warehouse/model"
	warehouseRepo "github.com/antoniusDoni/monorepo/modules/warehouse/repository"
	"github.com/antoniusDoni/monorepo/shared/auth"
	"github.com/antoniusDoni/monorepo/shared/contract"
	"github.com/antoniusDoni/monorepo/shared/model"
	"github.com/antoniusDoni/monorepo/shared/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthService handles user registration, login, and retrieval
type AuthService struct {
	repo       repository.UserRepository
	officeRepo warehouseRepo.OfficeRepository
	db         *gorm.DB
	jwtSecret  string
}

// NewAuthService creates a new AuthService instance
func NewAuthService(repo repository.UserRepository, officeRepo warehouseRepo.OfficeRepository, db *gorm.DB, jwtSecret string) *AuthService {
	return &AuthService{repo: repo, officeRepo: officeRepo, db: db, jwtSecret: jwtSecret}
}

// Register creates a new user with hashed password
func (s *AuthService) Register(username, password, email, officeID string) error {
	// Check if username already exists
	_, err := s.repo.FindByUsername(username)
	if err == nil {
		return errors.New("username already exists")
	}

	// Check if email already exists
	_, err = s.repo.FindByEmail(email)
	if err == nil {
		return errors.New("email already exists")
	}

	// Validate office exists
	officeUUID, err := uuid.Parse(officeID)
	if err != nil {
		return errors.New("invalid office ID format")
	}

	_, err = s.officeRepo.GetByID(officeID)
	if err != nil {
		return errors.New("office not found")
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &model.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hashed),
		OfficeID:     officeUUID,
	}

	err = s.repo.Create(user)
	if err != nil {
		return err
	}

	// Assign admin role to the newly created user
	err = s.assignAdminRole(user.ID)
	if err != nil {
		// Log the error but don't fail the registration
		// The user is created successfully, just without role assignment
		return errors.New("user created but role assignment failed: " + err.Error())
	}

	return nil
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

// RegisterWithOffice creates a new office and user in a single transaction
func (s *AuthService) RegisterWithOffice(req *contract.RegisterWithOfficeRequest) (*contract.RegisterWithOfficeResponse, error) {
	// Check if username already exists
	_, err := s.repo.FindByUsername(req.Username)
	if err == nil {
		return nil, errors.New("username already exists")
	}

	// Check if email already exists
	_, err = s.repo.FindByEmail(req.Email)
	if err == nil {
		return nil, errors.New("email already exists")
	}

	// Check if office code already exists
	existingOffice, _ := s.officeRepo.GetByCode(req.OfficeCode)
	if existingOffice != nil {
		return nil, errors.New("office code already exists")
	}

	// Create office first
	office := &warehouseModel.Office{
		ID:        uuid.New(),
		Code:      req.OfficeCode,
		Name:      req.OfficeName,
		Address:   req.OfficeAddress,
		City:      req.OfficeCity,
		Phone:     req.OfficePhone,
		Status:    "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = s.officeRepo.Create(office)
	if err != nil {
		return nil, errors.New("failed to create office: " + err.Error())
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user with the new office ID
	user := &model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashed),
		OfficeID:     office.ID,
		CreatedAt:    time.Now(),
	}

	err = s.repo.Create(user)
	if err != nil {
		// If user creation fails, we should ideally rollback the office creation
		// For now, we'll just return the error
		return nil, errors.New("failed to create user: " + err.Error())
	}

	// Assign admin role to the newly created user
	err = s.assignAdminRole(user.ID)
	if err != nil {
		// Log the error but don't fail the registration
		// The user is created successfully, just without role assignment
		return &contract.RegisterWithOfficeResponse{
			Message:  "Office and user registered successfully, but role assignment failed",
			OfficeID: office.ID.String(),
			UserID:   user.ID,
		}, nil
	}

	return &contract.RegisterWithOfficeResponse{
		Message:  "Office and user registered successfully",
		OfficeID: office.ID.String(),
		UserID:   user.ID,
	}, nil
}

// assignAdminRole assigns the admin role to a user
func (s *AuthService) assignAdminRole(userID uint) error {
	// Find the admin role
	var adminRole model.Role
	err := s.db.Where("name = ?", "admin").First(&adminRole).Error
	if err != nil {
		return errors.New("admin role not found")
	}

	// Find the user
	var user model.User
	err = s.db.First(&user, userID).Error
	if err != nil {
		return errors.New("user not found")
	}

	// Assign the role to the user
	err = s.db.Model(&user).Association("Roles").Append(&adminRole)
	if err != nil {
		return errors.New("failed to assign role to user")
	}

	return nil
}

// GetUser retrieves a user by ID
func (s *AuthService) GetUser(id uint) (*model.User, error) {
	return s.repo.FindByID(id)
}

func (s *AuthService) GetListOffice(page, pageSize int, searchTerm string) ([]warehouseModel.Office, int64, error) {
	return s.officeRepo.GetAll(page, pageSize, searchTerm)
}
