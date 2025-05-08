package handler

import (
	"context"
	"net/http"
	"thanhldt060802/internal/dto"
	"thanhldt060802/internal/middleware"
	"thanhldt060802/internal/service"

	"github.com/danielgtaylor/huma/v2"
)

type UserHandler struct {
	userService       service.UserService
	jwtAuthMiddleware *middleware.JWTAuthMiddleware
}

func NewUserHandler(api huma.API, userService service.UserService, jwtAuthMiddleware *middleware.JWTAuthMiddleware) *UserHandler {
	userHandler := &UserHandler{
		userService:       userService,
		jwtAuthMiddleware: jwtAuthMiddleware,
	}

	// Get users
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/users",
		Summary:     "/users",
		Description: "Get users.",
		Tags:        []string{"User"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication, jwtAuthMiddleware.RequireAdmin},
	}, userHandler.GetUsers)

	// Get user by id
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/users/id/{id}",
		Summary:     "/users/id/{id}",
		Description: "Get user by id.",
		Tags:        []string{"User"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication, jwtAuthMiddleware.RequireAdmin},
	}, userHandler.GetUserById)

	// Get user by username
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/users/username/{username}",
		Summary:     "/users/username/{username}",
		Description: "Get user by username.",
		Tags:        []string{"User"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication, jwtAuthMiddleware.RequireAdmin},
	}, userHandler.GetUserByUsername)

	// Get user by email
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/users/email/{email}",
		Summary:     "/users/email/{email}",
		Description: "Get user by email.",
		Tags:        []string{"User"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication, jwtAuthMiddleware.RequireAdmin},
	}, userHandler.GetUserByEmail)

	// Create user
	huma.Register(api, huma.Operation{
		Method:      http.MethodPost,
		Path:        "/users",
		Summary:     "/users",
		Description: "Create user.",
		Tags:        []string{"User"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication, jwtAuthMiddleware.RequireAdmin},
	}, userHandler.CreateUser)

	// Update user by id
	huma.Register(api, huma.Operation{
		Method:      http.MethodPut,
		Path:        "/users/id/{id}",
		Summary:     "/users/id/{id}",
		Description: "Update user by id.",
		Tags:        []string{"User"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication, jwtAuthMiddleware.RequireAdmin},
	}, userHandler.UpdateUserById)

	// Delete user by id
	huma.Register(api, huma.Operation{
		Method:      http.MethodDelete,
		Path:        "/users/id/{id}",
		Summary:     "/users/id/{id}",
		Description: "Delete user by id.",
		Tags:        []string{"User"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication, jwtAuthMiddleware.RequireAdmin},
	}, userHandler.DeleteUserById)

	// Kill user token
	huma.Register(api, huma.Operation{
		Method:      http.MethodDelete,
		Path:        "/users/kill-token/{token}",
		Summary:     "/users/kill-token/{token}",
		Description: "Kill user token.",
		Tags:        []string{"Account"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication, jwtAuthMiddleware.RequireAdmin},
	}, userHandler.KillUserToken)

	// Login user account
	huma.Register(api, huma.Operation{
		Method:      http.MethodPost,
		Path:        "/login",
		Summary:     "/login",
		Description: "Login user account.",
		Tags:        []string{"Account"},
	}, userHandler.LoginUserAccount)

	// Logout user account
	huma.Register(api, huma.Operation{
		Method:      http.MethodDelete,
		Path:        "/logout",
		Summary:     "/logout",
		Description: "Logout user account.",
		Tags:        []string{"Account"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication},
	}, userHandler.LogoutUserAccount)

	// Register user
	huma.Register(api, huma.Operation{
		Method:      http.MethodPost,
		Path:        "/register",
		Summary:     "/register",
		Description: "Register user account.",
		Tags:        []string{"Account"},
	}, userHandler.RegisterUserAccount)

	// Get user account
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/my-account",
		Summary:     "/my-account",
		Description: "Get user account.",
		Tags:        []string{"Account"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication},
	}, userHandler.GetUserAccount)

	// Update user account
	huma.Register(api, huma.Operation{
		Method:      http.MethodPut,
		Path:        "/my-account",
		Summary:     "/my-account",
		Description: "Update user account.",
		Tags:        []string{"Account"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication},
	}, userHandler.UpdateUserAccount)

	// Integrate with Elasticsearch

	// Sync all available users to Elasticsearch
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/user/sync-to-elasticsearch",
		Summary:     "/user/sync-to-elasticsearch",
		Description: "Sync all available users to Elasticsearch.",
		Tags:        []string{"Sync Data"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication, jwtAuthMiddleware.RequireAdmin},
	}, userHandler.SyncAllAvailableUsersToElasticsearch)

	// Get users with Elasticsearch
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/users/elasticsearch",
		Summary:     "/users/elasticsearch",
		Description: "Get users with Elasticsearch.",
		Tags:        []string{"User"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication, jwtAuthMiddleware.RequireAdmin},
	}, userHandler.GetUsersWithElasticsearch)

	return userHandler
}

func (userHandler *UserHandler) GetUsers(ctx context.Context, reqDTO *dto.GetUsersRequest) (*dto.PaginationBodyResponseList[dto.UserView], error) {
	users, err := userHandler.userService.GetUsers(ctx, reqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusInternalServerError
		res.Code = "ERR_INTERNAL_SERVER"
		res.Message = "Get users failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.PaginationBodyResponseList[dto.UserView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get users successful"
	res.Body.Data = users
	res.Body.Total = len(users)
	return res, nil
}

func (userHandler *UserHandler) GetUserById(ctx context.Context, reqDTO *dto.GetUserByIdRequest) (*dto.BodyResponse[dto.UserView], error) {
	foundUser, err := userHandler.userService.GetUserById(ctx, reqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Get user by id failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.BodyResponse[dto.UserView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get user by id successful"
	res.Body.Data = *foundUser
	return res, nil
}

func (userHandler *UserHandler) GetUserByUsername(ctx context.Context, reqDTO *dto.GetUserByUsernameRequest) (*dto.BodyResponse[dto.UserView], error) {
	foundUser, err := userHandler.userService.GetUserByUsername(ctx, reqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Get user by username failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.BodyResponse[dto.UserView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get user by username successful"
	res.Body.Data = *foundUser
	return res, nil
}

func (userHandler *UserHandler) GetUserByEmail(ctx context.Context, reqDTO *dto.GetUserByEmailRequest) (*dto.BodyResponse[dto.UserView], error) {
	foundUser, err := userHandler.userService.GetUserByEmail(ctx, reqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Get user by email failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.BodyResponse[dto.UserView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get user by email successful"
	res.Body.Data = *foundUser
	return res, nil
}

func (userHandler *UserHandler) CreateUser(ctx context.Context, reqDTO *dto.CreateUserRequest) (*dto.SuccessResponse, error) {
	if err := userHandler.userService.CreateUser(ctx, reqDTO); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Create user failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Create user successful"
	return res, nil
}

func (userHandler *UserHandler) UpdateUserById(ctx context.Context, reqDTO *dto.UpdateUserByIdRequest) (*dto.SuccessResponse, error) {
	if err := userHandler.userService.UpdateUserById(ctx, reqDTO); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Update user by id failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Update user by id successful"
	return res, nil
}

func (userHandler *UserHandler) DeleteUserById(ctx context.Context, reqDTO *dto.DeleteUserByIdRequest) (*dto.SuccessResponse, error) {
	if err := userHandler.userService.DeleteUserById(ctx, reqDTO); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Delete user by id failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Delete user by id successful"
	return res, nil
}

func (userHandler *UserHandler) LoginUserAccount(ctx context.Context, reqDTO *dto.LoginUserAccountRequest) (*dto.BodyResponse[string], error) {
	token, err := userHandler.userService.LoginUserAccount(ctx, reqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Login user account failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.BodyResponse[string]{}
	res.Body.Code = "OK"
	res.Body.Message = "Login user account successful"
	res.Body.Data = *token
	return res, nil
}

func (userHandler *UserHandler) KillUserToken(ctx context.Context, reqDTO *dto.KillUserTokenRequest) (*dto.SuccessResponse, error) {
	if err := userHandler.userService.KillUserToken(ctx, reqDTO); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Kill user token failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Kill user token successful"
	return res, nil
}

func (userHandler *UserHandler) LogoutUserAccount(ctx context.Context, _ *struct{}) (*dto.SuccessResponse, error) {
	if err := userHandler.userService.LogoutUserAccount(ctx); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Logout user account failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Logout user account successful"
	return res, nil
}

func (userHandler *UserHandler) RegisterUserAccount(ctx context.Context, reqDTO *dto.RegisterUserAccountRequest) (*dto.SuccessResponse, error) {
	if err := userHandler.userService.RegisterUserAccount(ctx, reqDTO); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Register user account failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Register user account successful"
	return res, nil
}

func (userHandler *UserHandler) GetUserAccount(ctx context.Context, _ *struct{}) (*dto.BodyResponse[dto.UserView], error) {
	foundUser, err := userHandler.userService.GetUserAccount(ctx)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Get user account failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.BodyResponse[dto.UserView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get user account successful"
	res.Body.Data = *foundUser
	return res, nil
}

func (userHandler *UserHandler) UpdateUserAccount(ctx context.Context, reqDTO *dto.UpdateUserAccountRequest) (*dto.SuccessResponse, error) {
	if err := userHandler.userService.UpdateUserAccount(ctx, reqDTO); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Update user account failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Update user account successful"
	return res, nil
}

func (userHandler *UserHandler) SyncAllAvailableUsersToElasticsearch(ctx context.Context, _ *struct{}) (*dto.SuccessResponse, error) {
	if err := userHandler.userService.SyncAllAvailableUsersToElasticsearch(ctx); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusInternalServerError
		res.Code = "ERR_INTERNAL_SERVER"
		res.Message = "Sync all available users to Elasticsearch failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Sync all available users to Elasticsearch successful"
	return res, nil
}

func (userHandler *UserHandler) GetUsersWithElasticsearch(ctx context.Context, reqDTO *dto.GetUsersWithElasticsearchRequest) (*dto.PaginationBodyResponseList[dto.UserView], error) {
	users, err := userHandler.userService.GetUsersWithElasticsearch(ctx, reqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusInternalServerError
		res.Code = "ERR_INTERNAL_SERVER"
		res.Message = "Get users with Elasticsearch failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.PaginationBodyResponseList[dto.UserView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get users with Elasticsearch successful"
	res.Body.Data = users
	res.Body.Total = len(users)
	return res, nil
}
