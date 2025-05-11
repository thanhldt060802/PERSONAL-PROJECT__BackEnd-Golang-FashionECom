package service

import (
	"context"
	"encoding/json"
	"fmt"
	"thanhldt060802/config"
	"thanhldt060802/infrastructure"
	"thanhldt060802/internal/dto"
	"thanhldt060802/internal/model"
	"thanhldt060802/internal/repository"
	"thanhldt060802/utils"
	"time"
)

type userService struct {
	userRepository repository.UserRepository
}

type UserService interface {
	// Integrate with Elasticsearch
	GetAllUsers(ctx context.Context) ([]dto.UserView, error)

	// Main features
	GetUserById(ctx context.Context, reqDTO *dto.GetUserByIdRequest) (*dto.UserView, error)
	CreateUser(ctx context.Context, reqDTO *dto.CreateUserRequest) error
	UpdateUserById(ctx context.Context, reqDTO *dto.UpdateUserByIdRequest) error
	DeleteUserById(ctx context.Context, reqDTO *dto.DeleteUserByIdRequest) error

	// Extra feature
	LoginAccount(ctx context.Context, reqDTO *dto.LoginAccountRequest) (*string, error)
	LogoutAccount(ctx context.Context, token string) error
	// ShowAllLoggedInAccounts()
	// DestroyLoggedInAccountToken()

	// Elasticsearch integration features
	// GetUsers()
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}

//
//
// Integrate with Elasticsearch
// ######################################################################################

func (userService *userService) GetAllUsers(ctx context.Context) ([]dto.UserView, error) {
	users, err := userService.userRepository.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("query users from postgresql failed: %s", err.Error())
	}

	return dto.ToListUserView(users), nil
}

//
//
// Main features
// ######################################################################################

func (userService *userService) GetUserById(ctx context.Context, reqDTO *dto.GetUserByIdRequest) (*dto.UserView, error) {
	foundUser, err := userService.userRepository.GetById(ctx, reqDTO.Id)
	if err != nil {
		return nil, fmt.Errorf("id of user is not valid: %s", err.Error())
	}

	return dto.ToUserView(foundUser), nil
}

func (userService *userService) CreateUser(ctx context.Context, reqDTO *dto.CreateUserRequest) error {
	if _, err := userService.userRepository.GetByUsername(ctx, reqDTO.Body.Username); err == nil {
		return fmt.Errorf("username of user is already exists")
	}
	if _, err := userService.userRepository.GetByEmail(ctx, reqDTO.Body.Email); err == nil {
		return fmt.Errorf("email of user is already exists")
	}

	hashedPassword, err := utils.GenerateHashedPassword(reqDTO.Body.Password)
	if err != nil {
		return fmt.Errorf("generate hashed password failed: %s", err.Error())
	}

	newUser := model.User{
		FullName:       reqDTO.Body.FullName,
		Email:          reqDTO.Body.Email,
		Username:       reqDTO.Body.Username,
		HashedPassword: hashedPassword,
		Address:        reqDTO.Body.Address,
		RoleName:       reqDTO.Body.RoleName,
	}
	if err := userService.userRepository.Create(ctx, &newUser); err != nil {
		return fmt.Errorf("insert user to postgresql failed: %s", err.Error())
	}

	// Missing->SyncCreatingToElasticsearch

	return nil
}

func (userService *userService) UpdateUserById(ctx context.Context, reqDTO *dto.UpdateUserByIdRequest) error {
	foundUser, err := userService.userRepository.GetById(ctx, reqDTO.Id)
	if err != nil {
		return fmt.Errorf("id of user is not valid: %s", err.Error())
	}

	if reqDTO.Body.FullName != nil {
		foundUser.FullName = *reqDTO.Body.FullName
	}
	if reqDTO.Body.Email != nil && reqDTO.Body.Email != &foundUser.Email {
		if _, err = userService.userRepository.GetByEmail(ctx, *reqDTO.Body.Email); err == nil {
			return fmt.Errorf("email of user is already exists")
		}
		foundUser.Email = *reqDTO.Body.Email
	}
	if reqDTO.Body.Password != nil {
		hashedPassword, err := utils.GenerateHashedPassword(*reqDTO.Body.Password)
		if err != nil {
			return fmt.Errorf("generate hashed password failed: %s", err.Error())
		}
		foundUser.HashedPassword = hashedPassword
	}
	if reqDTO.Body.Address != nil {
		foundUser.Address = *reqDTO.Body.Address
	}
	if reqDTO.Body.RoleName != nil {
		foundUser.RoleName = *reqDTO.Body.RoleName
	}
	foundUser.UpdatedAt = time.Now().UTC()

	if err := userService.userRepository.Update(ctx, foundUser); err != nil {
		return fmt.Errorf("update user on postgresql failed: %s", err.Error())
	}

	// Missing->SyncUpdatingToElasticsearch

	return nil
}

func (userService *userService) DeleteUserById(ctx context.Context, reqDTO *dto.DeleteUserByIdRequest) error {
	if _, err := userService.userRepository.GetById(ctx, reqDTO.Id); err != nil {
		return fmt.Errorf("id of user is not valid")
	}

	if err := userService.userRepository.DeleteById(ctx, reqDTO.Id); err != nil {
		return fmt.Errorf("delete user from postgresql failed: %s", err.Error())
	}

	// Missing->SyncDeletingToElasticsearch

	return nil
}

//
//
// Extra features
// ######################################################################################

func (userService *userService) LoginAccount(ctx context.Context, reqDTO *dto.LoginAccountRequest) (*string, error) {
	foundUser, err := userService.userRepository.GetByUsername(ctx, reqDTO.Body.Username)
	if err != nil {
		return nil, fmt.Errorf("username of user is not valid")
	}

	if utils.ValidatePassword(foundUser.HashedPassword, reqDTO.Body.Password) != nil {
		return nil, fmt.Errorf("password of user does not match")
	}

	tokenStr, err := utils.GenerateToken(foundUser.Id, foundUser.RoleName)
	if err != nil {
		return nil, fmt.Errorf("generate token failed: %s", err.Error())
	}

	redisKey := fmt.Sprintf("token:%s", tokenStr)
	userData := map[string]interface{}{
		"user_id":   foundUser.Id,
		"role_name": foundUser.RoleName,
	}
	userDataJSON, _ := json.Marshal(userData)

	status, err := infrastructure.RedisClient.SetEx(ctx, redisKey, userDataJSON, *config.AppConfig.TokenExpireMinutesValue()).Result()
	if err != nil {
		return nil, fmt.Errorf("save token to redis failed: %s", err.Error())
	}
	if status != "OK" {
		return nil, fmt.Errorf("unexpected response from redis - status: %s", status)
	}

	return &tokenStr, nil
}

func (userService *userService) LogoutAccount(ctx context.Context, token string) error {
	redisKey := fmt.Sprintf("token:%s", token)

	deleted, err := infrastructure.RedisClient.Del(ctx, redisKey).Result()
	if err != nil {
		return fmt.Errorf("delete token from redis failed: %s", err.Error())
	}
	if deleted == 0 {
		return fmt.Errorf("token is not valid or expired")
	}

	return nil
}
