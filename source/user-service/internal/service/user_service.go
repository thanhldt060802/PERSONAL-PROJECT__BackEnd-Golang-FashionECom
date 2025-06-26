package service

import (
	"context"
	"encoding/json"
	"fmt"
	"thanhldt060802/config"
	"thanhldt060802/infrastructure"
	"thanhldt060802/internal/dto"
	"thanhldt060802/internal/grpc/client/elasticsearchservicepb"
	"thanhldt060802/internal/model"
	"thanhldt060802/internal/repository"
	"thanhldt060802/utils"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type userService struct {
	userRepository repository.UserRepository
}

type UserService interface {
	GetUserById(ctx context.Context, reqDTO *dto.GetUserByIdRequest) (*dto.UserView, error)
	CreateUser(ctx context.Context, reqDTO *dto.CreateUserRequest) error
	UpdateUserById(ctx context.Context, reqDTO *dto.UpdateUserByIdRequest) error
	DeleteUserById(ctx context.Context, reqDTO *dto.DeleteUserByIdRequest) error

	LoginAccount(ctx context.Context, reqDTO *dto.LoginAccountRequest) (string, error)
	LogoutAccount(ctx context.Context, id string) error
	GetAllLoggedInAccounts(ctx context.Context) ([]string, error)

	// Elasticsearch integration (init data for elasticsearch-service)
	GetAllUsers(ctx context.Context) ([]*dto.UserView, error)

	// Elasticsearch integration features
	GetUsers(ctx context.Context, reqDTO *dto.GetUsersRequest) ([]*dto.UserView, error)
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
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
		Id:             uuid.New().String(),
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

	newUserView := dto.ToUserView(&newUser)
	payload, _ := json.Marshal(newUserView)
	if err := infrastructure.RedisClient.Publish(ctx, "user-service.created-user", payload).Err(); err != nil {
		return fmt.Errorf("pulish event user-service.created-user failed: %s", err.Error())
	}

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
	timeUpdate := time.Now().UTC()
	foundUser.UpdatedAt = &timeUpdate

	if err := userService.userRepository.Update(ctx, foundUser); err != nil {
		return fmt.Errorf("update user on postgresql failed: %s", err.Error())
	}

	updatedUserView := dto.ToUserView(foundUser)
	payload, _ := json.Marshal(updatedUserView)
	if err := infrastructure.RedisClient.Publish(ctx, "user-service.updated-user", payload).Err(); err != nil {
		return fmt.Errorf("pulish event user-service.updated-user failed: %s", err.Error())
	}

	return nil
}

func (userService *userService) DeleteUserById(ctx context.Context, reqDTO *dto.DeleteUserByIdRequest) error {
	if _, err := userService.userRepository.GetById(ctx, reqDTO.Id); err != nil {
		return fmt.Errorf("id of user is not valid")
	}

	userService.LogoutAccount(ctx, reqDTO.Id)

	if err := userService.userRepository.DeleteById(ctx, reqDTO.Id); err != nil {
		return fmt.Errorf("delete user from postgresql failed: %s", err.Error())
	}

	if err := infrastructure.RedisClient.Publish(ctx, "user-service.deleted-user", reqDTO.Id).Err(); err != nil {
		return fmt.Errorf("pulish event user-service.deleted-user failed: %s", err.Error())
	}

	return nil
}

func (userService *userService) LoginAccount(ctx context.Context, reqDTO *dto.LoginAccountRequest) (string, error) {
	foundUser, err := userService.userRepository.GetByUsername(ctx, reqDTO.Body.Username)
	if err != nil {
		return "", fmt.Errorf("username of user is not valid")
	}

	redisKey := fmt.Sprintf("%s:token", foundUser.Id)
	tokenStr, err := infrastructure.RedisClient.Get(ctx, redisKey).Result()
	if err == redis.Nil {
		if utils.ValidatePassword(foundUser.HashedPassword, reqDTO.Body.Password) != nil {
			return "", fmt.Errorf("password of user does not match")
		}

		tokenStr, err = utils.GenerateToken(foundUser.Id, foundUser.RoleName)
		if err != nil {
			return "", fmt.Errorf("generate token failed: %s", err.Error())
		}

		redisKey = fmt.Sprintf("token:%s", tokenStr)
		userData := map[string]interface{}{
			"user_id":   foundUser.Id,
			"role_name": foundUser.RoleName,
		}
		userDataJSON, _ := json.Marshal(userData)
		status, err := infrastructure.RedisClient.SetEx(ctx, redisKey, userDataJSON, config.AppConfig.TokenExpireMinutesValue()).Result()
		if err != nil {
			return "", fmt.Errorf("save token to redis failed: %s", err.Error())
		}
		if status != "OK" {
			return "", fmt.Errorf("unexpected response from redis - status: %s", status)
		}

		redisKey = fmt.Sprintf("%s:token", foundUser.Id)
		status, err = infrastructure.RedisClient.SetEx(ctx, redisKey, tokenStr, config.AppConfig.TokenExpireMinutesValue()).Result()
		if err != nil {
			return "", fmt.Errorf("save logged in account to redis failed: %s", err.Error())
		}
		if status != "OK" {
			return "", fmt.Errorf("unexpected response from redis - status: %s", status)
		}
	} else if err != nil {
		return "", fmt.Errorf("check logged in account from redis failed: %s", err.Error())
	}

	return tokenStr, nil
}

func (userService *userService) LogoutAccount(ctx context.Context, id string) error {
	redisKey := fmt.Sprintf("%s:token", id)
	tokenStr, err := infrastructure.RedisClient.Get(ctx, redisKey).Result()
	if err == redis.Nil {
		return fmt.Errorf("logged in account is not valid")
	} else if err != nil {
		return fmt.Errorf("check logged in account from redis failed: %s", err.Error())
	}

	deleted, err := infrastructure.RedisClient.Del(ctx, redisKey).Result()
	if err != nil {
		return fmt.Errorf("delete logged id account from redis failed: %s", err.Error())
	}
	if deleted == 0 {
		return fmt.Errorf("logged in account is not valid")
	}

	redisKey = fmt.Sprintf("token:%s", tokenStr)
	deleted, err = infrastructure.RedisClient.Del(ctx, redisKey).Result()
	if err != nil {
		return fmt.Errorf("delete token from redis failed: %s", err.Error())
	}
	if deleted == 0 {
		return fmt.Errorf("token is not valid or expired")
	}

	return nil
}

func (userService *userService) GetAllLoggedInAccounts(ctx context.Context) ([]string, error) {
	var cursor uint64
	var loggedInAccounts []string

	for {
		keys, nextCursor, err := infrastructure.RedisClient.Scan(ctx, cursor, "token:*", 100).Result()
		if err != nil {
			return nil, fmt.Errorf("scan keys on redis failed: %s", err.Error())
		}

		for _, key := range keys {
			userDataJson, err := infrastructure.RedisClient.Get(ctx, key).Result()
			if err != nil {
				return nil, fmt.Errorf("check token on redis failed: %s", err.Error())
			}

			var userData struct {
				UserId   string `json:"user_id"`
				RoleName string `json:"role_name"`
			}
			json.Unmarshal([]byte(userDataJson), &userData)

			loggedInAccounts = append(loggedInAccounts, userData.UserId)
		}

		if nextCursor == 0 {
			break
		}
		cursor = nextCursor
	}

	return loggedInAccounts, nil
}

func (userService *userService) GetAllUsers(ctx context.Context) ([]*dto.UserView, error) {
	users, err := userService.userRepository.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("query users from postgresql failed: %s", err.Error())
	}

	return dto.ToListUserView(users), nil
}

func (userService *userService) GetUsers(ctx context.Context, reqDTO *dto.GetUsersRequest) ([]*dto.UserView, error) {
	if infrastructure.ElasticsearchServiceGRPCClient != nil {
		convertReqDTO := &elasticsearchservicepb.GetUsersRequest{}
		convertReqDTO.Offset = reqDTO.Offset
		convertReqDTO.Limit = reqDTO.Limit
		convertReqDTO.SortBy = reqDTO.SortBy
		convertReqDTO.FullName = reqDTO.FullName
		convertReqDTO.Email = reqDTO.Email
		convertReqDTO.Username = reqDTO.Username
		convertReqDTO.Address = reqDTO.Address
		convertReqDTO.RoleName = reqDTO.RoleName
		convertReqDTO.CreatedAtGte = reqDTO.CreatedAtGTE
		convertReqDTO.CreatedAtLte = reqDTO.CreatedAtLTE

		grpcRes, err := infrastructure.ElasticsearchServiceGRPCClient.GetUsers(ctx, convertReqDTO)
		if err != nil {
			return nil, fmt.Errorf("get users from user-service failed: %s", err.Error())
		}

		return dto.FromListUserProtoToListUserView(grpcRes.Users), nil
	} else {
		return nil, fmt.Errorf("elasticsearch-service is not running")
	}
}
