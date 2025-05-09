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

	// Sync all available users
	huma.Register(api, huma.Operation{
		Method:      http.MethodPost,
		Path:        "/user-service/sync",
		Summary:     "/user-service/sync",
		Description: "Sync all available users.",
		Tags:        []string{"Sync Data"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication, jwtAuthMiddleware.RequireAdmin},
	}, userHandler.SyncAllAvailableUsers)

	// // Statistics number of users created
	// huma.Register(api, huma.Operation{
	// 	Method:      http.MethodGet,
	// 	Path:        "/users/statistics-number-of-users-created",
	// 	Summary:     "/users/statistics-number-of-users-created",
	// 	Description: "Statistics number of users created.",
	// 	Tags:        []string{"User"},
	// 	Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication, jwtAuthMiddleware.RequireAdmin},
	// }, userHandler.StatisticsNumberOfUsersCreated)

	return userHandler
}

func (userHandler *UserHandler) SyncAllAvailableUsers(ctx context.Context, reqDTO *dto.SyncAllAvailableUsersRequest) (*dto.SuccessResponse, error) {
	if err := userHandler.userService.SyncAllAvailableUsers(ctx, reqDTO); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusInternalServerError
		res.Code = "ERR_INTERNAL_SERVER"
		res.Message = "Sync all available users failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Sync all available users successful"
	return res, nil
}

// func (userHandler *UserHandler) StatisticsNumberOfUsersCreated(ctx context.Context, reqDTO *dto.StatisticsNumberOfUsersCreatedRequest) (*dto.BodyResponse[dto.NumberOfUsersCreatedReport], error) {
// 	report, err := userHandler.userService.StatisticsNumberOfUsersCreated(ctx, reqDTO)
// 	if err != nil {
// 		res := &dto.ErrorResponse{}
// 		res.Status = http.StatusInternalServerError
// 		res.Code = "ERR_INTERNAL_SERVER"
// 		res.Message = "Statistics number of users create failed"
// 		res.Details = []string{err.Error()}
// 		return nil, res
// 	}

// 	res := &dto.BodyResponse[dto.NumberOfUsersCreatedReport]{}
// 	res.Body.Code = "OK"
// 	res.Body.Message = "Statistics number of users created successful"
// 	res.Body.Data = *report
// 	return res, nil
// }
