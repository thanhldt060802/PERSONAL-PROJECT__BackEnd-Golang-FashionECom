package main

import (
	"net/http"
	"thanhldt060802/config"
	"thanhldt060802/infrastructure"
	"thanhldt060802/internal/dto"
	grpc_server "thanhldt060802/internal/grpc-server"
	"thanhldt060802/internal/handler"
	"thanhldt060802/internal/middleware"
	"thanhldt060802/internal/repository"
	"thanhldt060802/internal/service"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
)

// Huma Docs UI template by Scalar
var humaDocsEmbedded = `<!doctype html>
<html>
  <head>
    <title>FashionECom APIs</title>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
  </head>
  <body>
    <script
      id="api-reference"
      data-url="/openapi.json"></script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
  </body>
</html>`

func main() {

	config.InitConfig()
	infrastructure.InitPostgesDB()
	defer infrastructure.PostgresDB.Close()
	infrastructure.InitRedisClient()
	defer infrastructure.RedisClient.Close()

	humaCfg := huma.DefaultConfig("FashionECom - User Service", "v1.0.0")
	humaCfg.DocsPath = ""
	humaCfg.JSONSchemaDialect = ""
	humaCfg.CreateHooks = nil
	humaCfg.Components = &huma.Components{
		SecuritySchemes: map[string]*huma.SecurityScheme{
			"BearerAuth": {
				Type:         "http",
				Scheme:       "bearer",
				BearerFormat: "JWT",
			},
		},
	}

	huma.NewError = func(status int, msg string, errs ...error) huma.StatusError {
		details := make([]string, len(errs))
		for i, err := range errs {
			details[i] = err.Error()
		}
		res := &dto.ErrorResponse{}
		res.Status = status
		res.Message = msg
		res.Details = details
		return res
	}

	r := gin.Default()
	r.GET("/fashionecom/api-document", func(ctx *gin.Context) {
		ctx.Data(http.StatusOK, "text/html", []byte(humaDocsEmbedded))
	})

	api := humagin.New(r, humaCfg)

	jwtAuthMiddleware := middleware.NewAuthMiddleware()

	userRepository := repository.NewUserRepository()

	userService := service.NewUserService(userRepository)

	userHandler := handler.NewUserHandler(api, userService, jwtAuthMiddleware)

	grpcConfig := &grpc_server.GRPCServerConfig{}
	grpcConfig.Port = 50051
	grpcConfig.GetAllUsersImpl = grpc_server.NewGRPCUserService(userHandler)
	grpc_server.StartGRPCServer(*grpcConfig)

	r.Run(":" + config.AppConfig.AppPort)

}
