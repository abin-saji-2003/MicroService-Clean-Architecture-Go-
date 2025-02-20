package handlers

import (
	"context"
	"errors"
	"fmt"

	userProto "github.com/abin-saji-2003/GRPC-Pkg/proto/userpb"
	"user-service/internal/models"
	"user-service/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	repo repository.UserRepository
	userProto.UnimplementedUserServiceServer
}

func NewUserHandler(repo repository.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

// Register User
func (h *UserHandler) RegisterUser(ctx context.Context, req *userProto.RegisterUserRequest) (*userProto.RegisterUserResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := h.repo.CreateUser(user); err != nil {
		return nil, errors.New("failed to create user")
	}

	return &userProto.RegisterUserResponse{Message: "User registered successfully"}, nil
}

// Login User
func (h *UserHandler) LoginUser(ctx context.Context, req *userProto.LoginUserRequest) (*userProto.LoginUserResponse, error) {
	user, err := h.repo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("incorrect password")
	}

	return &userProto.LoginUserResponse{
		UserId: uint32(user.ID),
		Name:   user.Name,
		Email:  user.Email,
	}, nil
}

func (h *UserHandler) GetUserByID(ctx context.Context, req *userProto.GetUserByIDRequest) (*userProto.GetUserByIDResponse, error) {
	user, err := h.repo.GetUserByID(uint(req.UserId))
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return &userProto.GetUserByIDResponse{
		UserId: uint32(user.ID),
		Name:   user.Name,
		Email:  user.Email,
	}, nil
}
