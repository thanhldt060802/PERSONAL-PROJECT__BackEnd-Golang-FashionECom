package repository

import (
	"context"
	"log"
	"thanhldt060802/infrastructure"
	"thanhldt060802/internal/model"
)

func InitTableCartItem() {
	ctx := context.Background()

	var exists bool
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM information_schema.tables 
			WHERE table_schema = 'public' AND table_name = ?
		)
	`
	if err := infrastructure.PostgresDB.QueryRowContext(ctx, query, "tb_cart_item").Scan(&exists); err != nil {
		log.Fatal("Check table tb_cart_item on PostgreSQL failed: ", err)
	}

	if !exists {
		if _, err := infrastructure.PostgresDB.NewCreateTable().Model(&model.CartItem{}).Exec(ctx); err != nil {
			log.Fatal("Create table tb_cart_item on PostgreSQL failed: ", err)
		}
	}
}

func InitTableInvoice() {
	ctx := context.Background()

	var exists bool
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM information_schema.tables 
			WHERE table_schema = 'public' AND table_name = ?
		)
	`
	if err := infrastructure.PostgresDB.QueryRowContext(ctx, query, "tb_invoice").Scan(&exists); err != nil {
		log.Fatal("Check table tb_invoice on PostgreSQL failed: ", err)
	}

	if !exists {
		if _, err := infrastructure.PostgresDB.NewCreateTable().Model(&model.Invoice{}).Exec(ctx); err != nil {
			log.Fatal("Create table tb_invoice on PostgreSQL failed: ", err)
		}
	}
}

func InitTableInvoiceDetail() {
	ctx := context.Background()

	var exists bool
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM information_schema.tables 
			WHERE table_schema = 'public' AND table_name = ?
		)
	`
	if err := infrastructure.PostgresDB.QueryRowContext(ctx, query, "tb_invoice_detail").Scan(&exists); err != nil {
		log.Fatal("Check table tb_invoice_detail on PostgreSQL failed: ", err)
	}

	if !exists {
		if _, err := infrastructure.PostgresDB.NewCreateTable().Model(&model.InvoiceDetail{}).Exec(ctx); err != nil {
			log.Fatal("Create table tb_invoice_detail on PostgreSQL failed: ", err)
		}
	}
}
