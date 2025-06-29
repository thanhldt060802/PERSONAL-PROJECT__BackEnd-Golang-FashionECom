package repository

import (
	"context"
	"fmt"
	"log"
	"thanhldt060802/infrastructure"
	"thanhldt060802/internal/model"
	"thanhldt060802/utils"

	"github.com/google/uuid"
)

func InitTableUsers() {
	ctx := context.Background()

	var exists bool
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM information_schema.tables 
			WHERE table_schema = 'public' AND table_name = ?
		)
	`
	if err := infrastructure.PostgresDB.QueryRowContext(ctx, query, "tb_user").Scan(&exists); err != nil {
		log.Fatal("Check table users on PostgreSQL failed: ", err)
	}

	if !exists {
		if _, err := infrastructure.PostgresDB.NewCreateTable().Model(&model.User{}).Exec(ctx); err != nil {
			log.Fatal("Create table users on PostgreSQL failed: ", err)
		}

		userData := []*model.User{}

		adminHashedPassword, _ := utils.GenerateHashedPassword("123")
		userData = append(userData, &model.User{
			Id:             uuid.New().String(),
			FullName:       "Full Name Of Admin",
			Email:          "admin@gmail.com",
			Username:       "admin",
			HashedPassword: adminHashedPassword,
			Address:        "Củ chi",
			RoleName:       "ADMIN",
		})

		for i := range 20 {
			userHashedPassword, _ := utils.GenerateHashedPassword("123")
			userData = append(userData, &model.User{
				Id:             uuid.New().String(),
				FullName:       fmt.Sprintf("Full Name Of User %v", i+1),
				Email:          fmt.Sprintf("user%v@gmail.com", i+1),
				Username:       fmt.Sprintf("user%v", i+1),
				HashedPassword: userHashedPassword,
				Address:        fmt.Sprintf("Address Of User %v", i+1),
				RoleName:       "CUSTOMER",
			})
		}

		if _, err := infrastructure.PostgresDB.NewInsert().Model(&userData).Exec(ctx); err != nil {
			log.Fatal("Create data for table users on PostgreSQL failed: ", err)
		}
	}
}
