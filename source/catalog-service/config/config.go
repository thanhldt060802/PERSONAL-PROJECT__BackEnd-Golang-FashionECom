package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string

	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string

	RedisHost     string
	RedisPort     string
	RedisPassword string

	CatalogServiceGRPCHost       string
	CatalogServiceGRPCPort       string
	ElasticsearchServiceGRPCHost string
	ElasticsearchServiceGRPCPort string
}

var AppConfig *Config

func InitConfig() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Load file .env failed: ", err)
	}

	AppConfig = &Config{
		AppPort: GetEnv("APP_PORT", "8080"),

		PostgresHost:     GetEnv("POSTGRES_HOST", "localhost"),
		PostgresPort:     GetEnv("POSTGRES_PORT", "5432"),
		PostgresUser:     GetEnv("POSTGRES_USER", "postgres"),
		PostgresPassword: GetEnv("POSTGRES_PASSWORD", ""),
		PostgresDB:       GetEnv("POSTGRES_DB", "my_db"),

		RedisHost:     GetEnv("REDIS_HOST", "localhost"),
		RedisPort:     GetEnv("REDIS_PORT", "6379"),
		RedisPassword: GetEnv("REDIS_PASSWORD", ""),

		CatalogServiceGRPCHost:       GetEnv("CATALOG_SERVICE_GRPC_HOST", "localhost"),
		CatalogServiceGRPCPort:       GetEnv("CATALOG_SERVICE_GRPC_PORT", "50050"),
		ElasticsearchServiceGRPCHost: GetEnv("ELASTICSEARCH_SERVICE_GRPC_HOST", "localhost"),
		ElasticsearchServiceGRPCPort: GetEnv("ELASTICSEARCH_SERVICE_GRPC_PORT", "50050"),
	}

	log.Println("Load .env file successful")
}

func GetEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	} else {
		return defaultValue
	}
}
