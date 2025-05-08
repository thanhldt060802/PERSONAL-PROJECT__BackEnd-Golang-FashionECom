package service

import (
	"context"
	"fmt"
	"thanhldt060802/internal/dto"
	"thanhldt060802/internal/repository"
)

type userService struct {
	userElasticsearchRepository repository.UserElasticsearchRepository
}

type UserService interface {
	StatisticsNumberOfUsersCreated(ctx context.Context, reqDTO *dto.StatisticsNumberOfUsersCreatedRequest) (*dto.NumberOfUsersCreatedReport, error)
}

func NewUserService(userElasticsearchRepository repository.UserElasticsearchRepository) UserService {
	return &userService{
		userElasticsearchRepository: userElasticsearchRepository,
	}
}

func (userService *userService) StatisticsNumberOfUsersCreated(ctx context.Context, reqDTO *dto.StatisticsNumberOfUsersCreatedRequest) (*dto.NumberOfUsersCreatedReport, error) {
	report, err := userService.userElasticsearchRepository.StatisticsNumberOfUsersCreated(ctx,
		reqDTO.TimeInterval,
		reqDTO.CreatedAtGTE,
		reqDTO.CreatedAtLTE,
	)
	if err != nil {
		return nil, fmt.Errorf("query users from elasticsearch failed: %s", err.Error())
	}

	return report, nil
}
