package infrastructure

import (
	"database/sql"
	"fmt"
	"log"
	"thanhldt060802/config"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

var PostgresDB *bun.DB

func InitPostgesDB() {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.AppConfig.PostgresUser,
		config.AppConfig.PostgresPassword,
		config.AppConfig.PostgresHost,
		config.AppConfig.PostgresPort,
		config.AppConfig.PostgresDB,
	)

	connection, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Connect to PostgreSQL failed: ", err)
	}

	if err := connection.Ping(); err != nil {
		log.Fatal("Ping to PostgreSQL failed: ", err)
	}

	PostgresDB = bun.NewDB(connection, pgdialect.New())

	if err := PostgresDB.Ping(); err != nil {
		log.Fatal("Ping to PostgreSQL with Bun failed: ", err)
	}

	log.Println("Connect to PostgreSQL with Bun successful")
}
