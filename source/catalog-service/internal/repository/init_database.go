package repository

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"thanhldt060802/infrastructure"
	"thanhldt060802/internal/model"

	"github.com/google/uuid"
)

func InitTableCategories() {
	ctx := context.Background()

	var exists bool
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM information_schema.tables 
			WHERE table_schema = 'public' AND table_name = ?
		)
	`
	if err := infrastructure.PostgresDB.QueryRowContext(ctx, query, "categories").Scan(&exists); err != nil {
		log.Fatal("Check table categories on PostgreSQL failed: ", err)
	}

	if !exists {
		if _, err := infrastructure.PostgresDB.NewCreateTable().Model(&model.Category{}).Exec(ctx); err != nil {
			log.Fatal("Create table categories on PostgreSQL failed: ", err)
		}

		categoryData := []*model.Category{}

		for i := range 10 {
			categoryData = append(categoryData, &model.Category{
				Id:   uuid.New().String(),
				Name: fmt.Sprintf("Name Of Category %v", i+1),
			})
		}

		if _, err := infrastructure.PostgresDB.NewInsert().Model(&categoryData).Exec(ctx); err != nil {
			log.Fatal("Create data for table categories on PostgreSQL failed: ", err)
		}
	}
}

func InitTableBrands() {
	ctx := context.Background()

	var exists bool
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM information_schema.tables 
			WHERE table_schema = 'public' AND table_name = ?
		)
	`
	if err := infrastructure.PostgresDB.QueryRowContext(ctx, query, "brands").Scan(&exists); err != nil {
		log.Fatal("Check table brands on PostgreSQL failed: ", err)
	}

	if !exists {
		if _, err := infrastructure.PostgresDB.NewCreateTable().Model(&model.Brand{}).Exec(ctx); err != nil {
			log.Fatal("Create table brands on PostgreSQL failed: ", err)
		}

		brandData := []*model.Brand{}

		for i := range 10 {
			brandData = append(brandData, &model.Brand{
				Id:          uuid.New().String(),
				Name:        fmt.Sprintf("Name Of Brand %v", i+1),
				Description: fmt.Sprintf("Description Of Brand %v", i+1),
			})
		}

		if _, err := infrastructure.PostgresDB.NewInsert().Model(&brandData).Exec(ctx); err != nil {
			log.Fatal("Create data for table brands on PostgreSQL failed: ", err)
		}
	}
}

func InitTableProducts() {
	ctx := context.Background()

	var exists bool
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM information_schema.tables 
			WHERE table_schema = 'public' AND table_name = ?
		)
	`
	if err := infrastructure.PostgresDB.QueryRowContext(ctx, query, "products").Scan(&exists); err != nil {
		log.Fatal("Check table products on PostgreSQL failed: ", err)
	}

	if !exists {
		if _, err := infrastructure.PostgresDB.NewCreateTable().Model(&model.Product{}).Exec(ctx); err != nil {
			log.Fatal("Create table products on PostgreSQL failed: ", err)
		}

		productData := []*model.Product{}

		sexs := []string{"MALE", "FEMALE", "UNISEX"}
		var categoryIds []string
		var brandIds []string

		if err := infrastructure.PostgresDB.NewSelect().Model(&model.Category{}).Column("id").Scan(ctx, &categoryIds); err != nil {
			log.Fatal("Get all category ids from table categories on PostgreSQL failed: ", err)
		}
		if err := infrastructure.PostgresDB.NewSelect().Model(&model.Brand{}).Column("id").Scan(ctx, &brandIds); err != nil {
			log.Fatal("Get all brand ids from table brands on PostgreSQL failed: ", err)
		}

		for i := range 50 {
			productData = append(productData, &model.Product{
				Id:                 uuid.New().String(),
				Name:               fmt.Sprintf("Name Of Product %v", i+1),
				Description:        fmt.Sprintf("Description Of Product %v", i+1),
				Sex:                sexs[rand.Intn(len(sexs))],
				Price:              50000 + int64(rand.Intn(91))*5000,
				DiscountPercentage: int32(rand.Intn(6)) * 5,
				Stock:              int32(rand.Intn(11)) + 10,
				ImageURL:           "image.png",
				CategoryId:         categoryIds[rand.Intn(len(categoryIds))],
				BrandId:            brandIds[rand.Intn(len(brandIds))],
			})
		}

		if _, err := infrastructure.PostgresDB.NewInsert().Model(&productData).Exec(ctx); err != nil {
			log.Fatal("Create data for table products on PostgreSQL failed: ", err)
		}
	}
}
