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

	// Login account
	huma.Register(api, huma.Operation{
		Method:      http.MethodPost,
		Path:        "/login",
		Summary:     "/login",
		Description: "Login account.",
		Tags:        []string{"Account"},
	}, userHandler.LoginAccount)

	// Logout account
	huma.Register(api, huma.Operation{
		Method:      http.MethodDelete,
		Path:        "/logout",
		Summary:     "/logout",
		Description: "Logout account.",
		Tags:        []string{"Account"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication},
	}, userHandler.LogoutAccount)

	// Register account
	huma.Register(api, huma.Operation{
		Method:      http.MethodPost,
		Path:        "/register",
		Summary:     "/register",
		Description: "Register account.",
		Tags:        []string{"Account"},
	}, userHandler.RegisterAccount)

	// Get my account
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/my-account",
		Summary:     "/my-account",
		Description: "Get account.",
		Tags:        []string{"Account"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication},
	}, userHandler.GetMyAccount)

	// Update my account
	huma.Register(api, huma.Operation{
		Method:      http.MethodPut,
		Path:        "/my-account",
		Summary:     "/my-account",
		Description: "Update account.",
		Tags:        []string{"Account"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication},
	}, userHandler.UpdateMyAccount)

	// Get all logged in accounts
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/logged-in-accounts",
		Summary:     "/logged-in-accounts",
		Description: "Show all logged in accounts.",
		Tags:        []string{"Account"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication, jwtAuthMiddleware.RequireAdmin},
	}, userHandler.GetAllLoggedInAccounts)

	// Delete logged in account
	huma.Register(api, huma.Operation{
		Method:      http.MethodDelete,
		Path:        "/logged-in-accounts/{id}",
		Summary:     "/logged-in-accounts/{id}",
		Description: "Delete logged in account.",
		Tags:        []string{"Account"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication, jwtAuthMiddleware.RequireAdmin},
	}, userHandler.DeleteLoggedInAccount)

	return userHandler
}

func (userHandler *UserHandler) GetUsers(ctx context.Context, reqDTO *dto.GetUsersRequest) (*dto.PaginationBodyResponseList[*dto.UserView], error) {
	users, err := userHandler.userService.GetUsers(ctx, reqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusInternalServerError
		res.Code = "ERR_INTERNAL_SERVER"
		res.Message = "Get users failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.PaginationBodyResponseList[*dto.UserView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get users successful"
	res.Body.Data = users
	res.Body.Total = len(users)
	return res, nil
}

func (userHandler *UserHandler) GetUserById(ctx context.Context, reqDTO *dto.GetUserByIdRequest) (*dto.BodyResponse[*dto.UserView], error) {
	if reqDTO.Id == "{id}" {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Get user by id failed"
		res.Details = []string{"missing path parameters: id"}
		return nil, res
	}

	foundUser, err := userHandler.userService.GetUserById(ctx, reqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Get user by id failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.BodyResponse[*dto.UserView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get user by id successful"
	res.Body.Data = foundUser
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
	if reqDTO.Id == "{id}" {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Update user by id failed"
		res.Details = []string{"missing path parameters: id"}
		return nil, res
	}

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
	if reqDTO.Id == "{id}" {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Delete user by id failed"
		res.Details = []string{"missing path parameters: id"}
		return nil, res
	}

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

func (userHandler *UserHandler) LoginAccount(ctx context.Context, reqDTO *dto.LoginAccountRequest) (*dto.BodyResponse[string], error) {
	token, err := userHandler.userService.LoginAccount(ctx, reqDTO)
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
	res.Body.Data = token
	return res, nil
}

func (userHandler *UserHandler) LogoutAccount(ctx context.Context, _ *struct{}) (*dto.SuccessResponse, error) {
	if err := userHandler.userService.LogoutAccount(ctx, ctx.Value("user_id").(string)); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusInternalServerError
		res.Code = "ERR_INTERNAL_SERVER"
		res.Message = "Logout account failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Logout account successful"
	return res, nil
}

func (userHandler *UserHandler) RegisterAccount(ctx context.Context, reqDTO *dto.RegisterAccountRequest) (*dto.SuccessResponse, error) {
	convertReqDTO := &dto.CreateUserRequest{}
	convertReqDTO.Body.FullName = reqDTO.Body.FullName
	convertReqDTO.Body.Email = reqDTO.Body.Email
	convertReqDTO.Body.Username = reqDTO.Body.Username
	convertReqDTO.Body.Password = reqDTO.Body.Password
	convertReqDTO.Body.Address = reqDTO.Body.Address
	convertReqDTO.Body.RoleName = "CUSTOMER"

	if err := userHandler.userService.CreateUser(ctx, convertReqDTO); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Register account failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Register account successful"
	return res, nil
}

func (userHandler *UserHandler) GetMyAccount(ctx context.Context, _ *struct{}) (*dto.BodyResponse[*dto.UserView], error) {
	convertReqDTO := &dto.GetUserByIdRequest{}
	convertReqDTO.Id = ctx.Value("user_id").(string)

	foundUser, err := userHandler.userService.GetUserById(ctx, convertReqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Get account failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.BodyResponse[*dto.UserView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get account successful"
	res.Body.Data = foundUser
	return res, nil
}

func (userHandler *UserHandler) UpdateMyAccount(ctx context.Context, reqDTO *dto.UpdateAccountRequest) (*dto.SuccessResponse, error) {
	convertReqDTO := &dto.UpdateUserByIdRequest{}
	convertReqDTO.Id = ctx.Value("user_id").(string)
	convertReqDTO.Body.FullName = reqDTO.Body.FullName
	convertReqDTO.Body.Email = reqDTO.Body.Email
	convertReqDTO.Body.Password = reqDTO.Body.Password
	convertReqDTO.Body.Address = reqDTO.Body.Address

	if err := userHandler.userService.UpdateUserById(ctx, convertReqDTO); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Update account failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Update account successful"
	return res, nil
}

func (userHandler *UserHandler) GetAllLoggedInAccounts(ctx context.Context, _ *struct{}) (*dto.BodyResponse[[]string], error) {
	loggedInAccounts, err := userHandler.userService.GetAllLoggedInAccounts(ctx)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusInternalServerError
		res.Code = "ERR_INTERNAL_SERVER"
		res.Message = "Get all logged in accounts failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.BodyResponse[[]string]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get all logged in accounts successful"
	res.Body.Data = loggedInAccounts
	return res, nil
}

func (userHandler *UserHandler) DeleteLoggedInAccount(ctx context.Context, reqDTO *dto.DeleteLoggedInAccountRequest) (*dto.SuccessResponse, error) {
	if reqDTO.Id == "{id}" {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Delete logged in account failed"
		res.Details = []string{"missing path parameters: id"}
		return nil, res
	}

	if err := userHandler.userService.LogoutAccount(ctx, reqDTO.Id); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Delete logged in account failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Delete logged in account successful"
	return res, nil
}
