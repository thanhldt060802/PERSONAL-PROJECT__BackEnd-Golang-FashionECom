package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string

	ElasticsearchHost     string
	ElasticsearchPort     string
	ElasticsearchUsername string
	ElasticsearchPassword string

	RedisHost     string
	RedisPort     string
	RedisPassword string

	ElasticsearchServiceGRPCHost string
	ElasticsearchServiceGRPCPort string
	UserServiceGRPCHost          string
	UserServiceGRPCPort          string
}

var AppConfig *Config

func InitConfig() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Load file .env failed: ", err)
	}

	AppConfig = &Config{
		AppPort: GetEnv("APP_PORT", "8080"),

		ElasticsearchHost:     GetEnv("ELASTICSEARCH_HOST", "localhost"),
		ElasticsearchPort:     GetEnv("ELASTICSEARCH_PORT", "9200"),
		ElasticsearchUsername: GetEnv("ELASTICSEARCH_USERNAME", "elastic"),
		ElasticsearchPassword: GetEnv("ELASTICSEARCH_PASSWORD", ""),

		RedisHost:     GetEnv("REDIS_HOST", "localhost"),
		RedisPort:     GetEnv("REDIS_PORT", "6379"),
		RedisPassword: GetEnv("REDIS_PASSWORD", ""),

		ElasticsearchServiceGRPCHost: GetEnv("ELASTICSEARCH_SERVICE_GRPC_HOST", "localhost"),
		ElasticsearchServiceGRPCPort: GetEnv("ELASTICSEARCH_SERVICE_GRPC_PORT", "50050"),
		UserServiceGRPCHost:          GetEnv("USER_SERVICE_GRPC_HOST", "localhost"),
		UserServiceGRPCPort:          GetEnv("USER_SERVICE_GRPC_PORT", "50050"),
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
