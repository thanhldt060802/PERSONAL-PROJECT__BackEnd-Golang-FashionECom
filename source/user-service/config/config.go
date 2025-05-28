package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string

	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string

	JWTSecret          string
	TokenExpireMinutes string

	RedisHost     string
	RedisPort     string
	RedisPassword string

	UserServiceGRPCHost          string
	UserServiceGRPCPort          string
	ElasticSearchServiceGRPCHost string
	ElasticSearchServiceGRPCPort string
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

		JWTSecret:          GetEnv("JWT_SECRET", "123"),
		TokenExpireMinutes: GetEnv("TOKEN_EXPIRE_MINUTES", "30"),

		RedisHost:     GetEnv("REDIS_HOST", "localhost"),
		RedisPort:     GetEnv("REDIS_PORT", "6379"),
		RedisPassword: GetEnv("REDIS_PASSWORD", ""),

		UserServiceGRPCHost:          GetEnv("USER_SERVICE_GRPC_HOST", "localhost"),
		UserServiceGRPCPort:          GetEnv("USER_SERVICE_GRPC_PORT", "50050"),
		ElasticSearchServiceGRPCHost: GetEnv("ELASTICSEARCH_SERVICE_GRPC_HOST", "localhost"),
		ElasticSearchServiceGRPCPort: GetEnv("ELASTICSEARCH_SERVICE_GRPC_PORT", "50050"),
	}

	// Validate constraint environment variable value
	if _, err := strconv.Atoi(AppConfig.TokenExpireMinutes); err != nil {
		log.Fatal("Evironment variable TOKEN_EXPIRE_MINUTES is not valid number (must int): ", err)
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

func (config *Config) TokenExpireMinutesValue() time.Duration {
	tokenExpireMinutes, _ := strconv.Atoi(config.TokenExpireMinutes)
	expireDuration := time.Duration(tokenExpireMinutes) * time.Minute
	return expireDuration
}
