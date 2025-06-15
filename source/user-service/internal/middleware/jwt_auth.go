package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"thanhldt060802/infrastructure"

	"github.com/danielgtaylor/huma/v2"
	"github.com/redis/go-redis/v9"
)

type JWTAuthMiddleware struct {
}

func NewAuthMiddleware() *JWTAuthMiddleware {
	return &JWTAuthMiddleware{}
}

func (jwtAuthMiddleware *JWTAuthMiddleware) Authentication(ctx huma.Context, next func(huma.Context)) {
	authHeader := ctx.Header("Authorization")
	if authHeader == "" {
		CustomHumaWriteErr(ctx, http.StatusUnauthorized, "ERR_UNAUTHORIZED", "Authorization header missing", []string{"invalid credentials"})
		return
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	redisKey := fmt.Sprintf("token:%s", tokenStr)
	userDataJson, err := infrastructure.RedisClient.Get(ctx.Context(), redisKey).Result()
	if err == redis.Nil {
		CustomHumaWriteErr(ctx, http.StatusUnauthorized, "ERR_UNAUTHORIZED", "Token not found or expired", []string{"invalid token"})
		return
	} else if err != nil {
		CustomHumaWriteErr(ctx, http.StatusUnauthorized, "ERR_UNAUTHORIZED", "Check token on Redis failed", []string{err.Error()})
		return
	}

	var userData struct {
		UserId   string `json:"user_id"`
		RoleName string `json:"role_name"`
	}
	json.Unmarshal([]byte(userDataJson), &userData)

	ctx = huma.WithValue(ctx, "user_id", userData.UserId)
	ctx = huma.WithValue(ctx, "role_name", userData.RoleName)

	next(ctx)
}

func (jwtAuthMiddleware *JWTAuthMiddleware) RequireAdmin(ctx huma.Context, next func(huma.Context)) {
	if roleName, _ := ctx.Context().Value("role_name").(string); roleName != "ADMIN" {
		CustomHumaWriteErr(ctx, http.StatusForbidden, "ERR_FORBIDDEN", "Access denied", []string{"no permission"})
		return
	}

	next(ctx)
}
