package repository

import (
	"context"
	"thanhldt060802/infrastructure"
	"thanhldt060802/internal/model"
)

type userRepository struct {
}

type UserRepository interface {
	GetById(ctx context.Context, id string) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Create(ctx context.Context, newUser *model.User) error
	Update(ctx context.Context, updatedUser *model.User) error
	DeleteById(ctx context.Context, id string) error

	// Elasticsearch integrattion (init data for elasticsearch-service)
	GetAll(ctx context.Context) ([]model.User, error)
}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (userRepository *userRepository) GetById(ctx context.Context, id string) (*model.User, error) {
	var user model.User

	if err := infrastructure.PostgresDB.NewSelect().Model(&user).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	return &user, nil
}

func (userRepository *userRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User

	if err := infrastructure.PostgresDB.NewSelect().Model(&user).Where("username = ?", username).Scan(ctx); err != nil {
		return nil, err
	}

	return &user, nil
}

func (userRepository *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User

	if err := infrastructure.PostgresDB.NewSelect().Model(&user).Where("email = ?", email).Scan(ctx); err != nil {
		return nil, err
	}

	return &user, nil
}

func (userRepository *userRepository) Create(ctx context.Context, newUser *model.User) error {
	_, err := infrastructure.PostgresDB.NewInsert().Model(newUser).Exec(ctx)
	return err
}

func (userRepository *userRepository) Update(ctx context.Context, updatedUser *model.User) error {
	_, err := infrastructure.PostgresDB.NewUpdate().Model(updatedUser).Where("id = ?", updatedUser.Id).Exec(ctx)
	return err
}

func (userRepository *userRepository) DeleteById(ctx context.Context, id string) error {
	_, err := infrastructure.PostgresDB.NewDelete().Model(&model.User{}).Where("id = ?", id).Exec(ctx)
	return err
}

func (userRepository *userRepository) GetAll(ctx context.Context) ([]model.User, error) {
	var users []model.User

	if err := infrastructure.PostgresDB.NewSelect().Model(&users).Scan(ctx); err != nil {
		return nil, err
	}

	return users, nil
}
