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
	GetAllUsers(ctx context.Context) ([]dto.UserView, error)

	GetUsers(ctx context.Context, reqDTO *dto.GetUsersRequest) ([]dto.UserView, error)
	GetUserById(ctx context.Context, reqDTO *dto.GetUserByIdRequest) (*dto.UserView, error)
	CreateUser(ctx context.Context, reqDTO *dto.CreateUserRequest) error
	UpdateUserById(ctx context.Context, reqDTO *dto.UpdateUserByIdRequest) error
	DeleteUserById(ctx context.Context, reqDTO *dto.DeleteUserByIdRequest) error
	KillAccountToken(ctx context.Context, reqDTO *dto.KillUserTokenRequest) error
	LoginUserAccount(ctx context.Context, reqDTO *dto.LoginUserAccountRequest) (*string, error)
	LogoutUserAccount(ctx context.Context) error
	RegisterUserAccount(ctx context.Context, reqDTO *dto.RegisterUserAccountRequest) error
	GetUserAccount(ctx context.Context) (*dto.UserView, error)
	UpdateUserAccount(ctx context.Context, reqDTO *dto.UpdateUserAccountRequest) error
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}

func (userService *userService) GetAllUsers(ctx context.Context) ([]dto.UserView, error) {
	users, err := userService.userRepository.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("query users from postgresql failed: %s", err.Error())
	}

	return dto.ToListUserView(users), nil
}

func (userService *userService) GetUsers(ctx context.Context, reqDTO *dto.GetUsersRequest) ([]dto.UserView, error) {
	sortFields := utils.ParseSorter(reqDTO.SortBy)

	users, err := userService.userRepository.Get(ctx, reqDTO.Offset, reqDTO.Limit, sortFields)
	if err != nil {
		return nil, fmt.Errorf("query users from postgresql failed: %s", err.Error())
	}

	return dto.ToListUserView(users), nil
}

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

	// Create cart for user
	// ...

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

	return nil
}

func (userService *userService) DeleteUserById(ctx context.Context, reqDTO *dto.DeleteUserByIdRequest) error {
	if _, err := userService.userRepository.GetById(ctx, reqDTO.Id); err != nil {
		return fmt.Errorf("id of user is not valid")
	}

	// Delete cart of user
	// ...

	if err := userService.userRepository.DeleteById(ctx, reqDTO.Id); err != nil {
		return fmt.Errorf("delete user from postgresql failed: %s", err.Error())
	}

	return nil
}

func (userService *userService) KillAccountToken(ctx context.Context, reqDTO *dto.KillUserTokenRequest) error {
	redisKey := fmt.Sprintf("token:%s", reqDTO.Token)

	deleted, err := infrastructure.RedisClient.Del(ctx, redisKey).Result()
	if err != nil {
		return fmt.Errorf("delete token from redis failed: %s", err.Error())
	}
	if deleted == 0 {
		return fmt.Errorf("token is not valid or expired")
	}

	return nil
}

func (userService *userService) LoginUserAccount(ctx context.Context, reqDTO *dto.LoginUserAccountRequest) (*string, error) {
	foundUser, err := userService.userRepository.GetByUsername(ctx, reqDTO.Body.Username)
	if err != nil {
		return nil, fmt.Errorf("username of user is not valid")
	}

	if utils.ValidatePassword(foundUser.HashedPassword, reqDTO.Body.Password) != nil {
		return nil, fmt.Errorf("password of user does not match")
	}

	// Get cart of user for storing cart_id in token
	// ...

	tokenStr, err := utils.GenerateToken(foundUser.Id, foundUser.RoleName, 0)
	if err != nil {
		return nil, fmt.Errorf("generate token failed: %s", err.Error())
	}

	redisKey := fmt.Sprintf("token:%s", tokenStr)
	userData := map[string]interface{}{
		"access_token": tokenStr,
		"user_id":      foundUser.Id,
		"role_name":    foundUser.RoleName,
		"cart_id":      0,
	}
	userDataJSON, _ := json.Marshal(userData)

	status, err := infrastructure.RedisClient.SetEx(ctx, redisKey, userDataJSON, *config.AppConfig.TokenExpireMinutesValue()).Result()
	if err != nil {
		return nil, fmt.Errorf("save token to redis failed: %w", err)
	}
	if status != "OK" {
		return nil, fmt.Errorf("unexpected response from redis: %s", status)
	}

	return &tokenStr, nil
}

func (userService *userService) LogoutUserAccount(ctx context.Context) error {
	convertReqDTO := &dto.KillUserTokenRequest{}
	convertReqDTO.Token = ctx.Value("access_token").(string)

	return userService.KillAccountToken(ctx, convertReqDTO)
}

func (userService *userService) RegisterUserAccount(ctx context.Context, reqDTO *dto.RegisterUserAccountRequest) error {
	convertReqDTO := &dto.CreateUserRequest{}
	convertReqDTO.Body.FullName = reqDTO.Body.FullName
	convertReqDTO.Body.Email = reqDTO.Body.Email
	convertReqDTO.Body.Username = reqDTO.Body.Username
	convertReqDTO.Body.Password = reqDTO.Body.Password
	convertReqDTO.Body.Address = reqDTO.Body.Address
	convertReqDTO.Body.RoleName = "CUSTOMER"

	return userService.CreateUser(ctx, convertReqDTO)
}

func (userService *userService) GetUserAccount(ctx context.Context) (*dto.UserView, error) {
	convertReqDTO := &dto.GetUserByIdRequest{}
	convertReqDTO.Id = ctx.Value("user_id").(int64)

	return userService.GetUserById(ctx, convertReqDTO)
}

func (userService *userService) UpdateUserAccount(ctx context.Context, reqDTO *dto.UpdateUserAccountRequest) error {
	convertReqDTO := &dto.UpdateUserByIdRequest{}
	convertReqDTO.Id = ctx.Value("user_id").(int64)
	convertReqDTO.Body.FullName = reqDTO.Body.FullName
	convertReqDTO.Body.Email = reqDTO.Body.Email
	convertReqDTO.Body.Password = reqDTO.Body.Password
	convertReqDTO.Body.Address = reqDTO.Body.Address

	return userService.UpdateUserById(ctx, convertReqDTO)
}
