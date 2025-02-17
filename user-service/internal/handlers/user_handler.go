package handlers

import (
	"context"
	"errors"
	"fmt"

	"github.com/abin-saji-2003/MicroService-Clean-Architecture-Go-/tree/main/user-service/api/proto"
	"github.com/abin-saji-2003/MicroService-Clean-Architecture-Go-/tree/main/user-service/internal/models"
	"github.com/abin-saji-2003/MicroService-Clean-Architecture-Go-/tree/main/user-service/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	repo repository.UserRepository
	proto.UnimplementedUserServiceServer
}

func NewUserHandler(repo repository.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

// Register User
func (h *UserHandler) RegisterUser(ctx context.Context, req *proto.RegisterUserRequest) (*proto.RegisterUserResponse, error) {
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

	return &proto.RegisterUserResponse{Message: "User registered successfully"}, nil
}

// Login User
func (h *UserHandler) LoginUser(ctx context.Context, req *proto.LoginUserRequest) (*proto.LoginUserResponse, error) {
	user, err := h.repo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("incorrect password")
	}

	return &proto.LoginUserResponse{
		UserId: uint32(user.ID),
		Name:   user.Name,
		Email:  user.Email,
	}, nil
}

func (h *UserHandler) GetUserByID(ctx context.Context, req *proto.GetUserByIDRequest) (*proto.GetUserByIDResponse, error) {
	user, err := h.repo.GetUserByID(uint(req.UserId))
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return &proto.GetUserByIDResponse{
		UserId: uint32(user.ID),
		Name:   user.Name,
		Email:  user.Email,
	}, nil
}
