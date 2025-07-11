package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/mutoulbj/gocsms/internal/dto"
	"github.com/mutoulbj/gocsms/internal/models"
	"github.com/mutoulbj/gocsms/internal/repository"
	"github.com/mutoulbj/gocsms/internal/utils"
	"github.com/sirupsen/logrus"
)

var (
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrUserNotFound          = errors.New("user not found")
)

type UserService struct {
	repository.UserRepository
	log *logrus.Logger
}

func NewUserService(repo *repository.UserRepository, log *logrus.Logger) *UserService {
	return &UserService{
		UserRepository: *repo,
		log:            log,
	}
}

// GetUserById retrieves a user by their ID
func (s *UserService) GetUserById(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error) {
	s.log.Infof("Retrieving user by ID: %s", id)

	user, err := s.UserRepository.GetUserById(ctx, id)
	if err != nil {
		s.log.WithError(err).Error("Failed to retrieve user by ID")
		return nil, err
	}
	return dto.ToUserResponse(user), nil
}

// GetUserByUsername retrieves a user by their username
func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*dto.UserResponse, error) {
	s.log.Infof("Retrieving user by username: %s", username)

	user, err := s.UserRepository.GetUserByUsername(ctx, username)
	if err != nil {
		s.log.WithError(err).Error("Failed to retrieve user by username")
		return nil, err
	}
	return dto.ToUserResponse(user), nil
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, req *dto.RegisterRequest) (*dto.UserResponse, error) {
	s.log.Infof("Creating user: %s", req.Username)
	// Check if the username already exists
	existingUser, err := s.UserRepository.GetUserByUsername(ctx, req.Username)
	if err != nil {
		s.log.WithError(err).Error("Failed to check existing user by username")
		return nil, err
	}
	if existingUser != nil {
		s.log.Warnf("Username %s already exists", req.Username)
		return nil, ErrUsernameAlreadyExists
	}

	existEmailUser, err := s.UserRepository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		s.log.WithError(err).Error("Failed to check existing user by email")
		return nil, err
	}
	if existEmailUser != nil {
		s.log.Warnf("Email %s already exists", req.Email)
		return nil, ErrEmailAlreadyExists
	}

	salt := utils.GenerateSalt()                                // Generate a random salt for password hashing
	hashedPassword, _ := utils.HashPassword(req.Password, salt) // Hash the password with the salt

	newUser := &models.User{
		Username:     req.Username,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		PhoneNumber:  req.PhoneNumber,
		PasswordHash: hashedPassword,
		Salt:         salt,
	}
	if err := s.UserRepository.Create(ctx, newUser); err != nil {
		s.log.WithError(err).Error("Failed to create user")
		return nil, err
	}
	return dto.ToUserResponse(newUser), nil
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(ctx context.Context, id uuid.UUID, req *dto.UpdateProfileRequest) (*dto.UserResponse, error) {
	s.log.Infof("Updating user with ID: %s", id)

	user, err := s.UserRepository.GetUserById(ctx, id)
	if err != nil {
		s.log.WithError(err).Error("Failed to retrieve user for update")
		return nil, ErrUserNotFound
	}

	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.PhoneNumber != "" {
		user.PhoneNumber = req.PhoneNumber
	}

	if err := s.UserRepository.Update(ctx, user); err != nil {
		s.log.WithError(err).Error("Failed to update user")
		return nil, err
	}

	return dto.ToUserResponse(user), nil
}

// ListUsers retrieves all users with pagination
func (s *UserService) ListUsers(ctx context.Context, username, email string, page, pageSize int) (*dto.UserListResponse, error) {
	s.log.Infof("Listing users with page: %d, pageSize: %d", page, pageSize)

	users, total, err := s.UserRepository.List(ctx, username, email, page, pageSize)
	if err != nil {
		s.log.WithError(err).Error("Failed to list users")
		return nil, err
	}

	var userResponses []*dto.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, dto.ToUserResponse(user))
	}

	return &dto.UserListResponse{
		Users:    userResponses,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}
