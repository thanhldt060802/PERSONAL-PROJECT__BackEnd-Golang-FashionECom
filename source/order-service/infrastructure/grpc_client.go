package infrastructure

import (
	"fmt"
	"log"
	"net"
	"thanhldt060802/config"
	"thanhldt060802/internal/grpc/client/catalogservicepb"
	"thanhldt060802/internal/grpc/client/elasticsearchservicepb"
	"thanhldt060802/internal/grpc/client/userservicepb"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var UserServiceGRPCClient userservicepb.UserServiceGRPCClient
var CatalogServiceGRPCClient catalogservicepb.CatalogServiceGRPCClient
var ElasticsearchServiceGRPCClient elasticsearchservicepb.ElasticsearchServiceGRPCClient

type serviceGRPCConnectionManager struct {
	userServiceGRPCConnection          *grpc.ClientConn
	catalogServiceGRPCConnection       *grpc.ClientConn
	elasticsearchServiceGRPCConnection *grpc.ClientConn
}

func (serviceGRPCConnectionManager *serviceGRPCConnectionManager) CloseAll() {
	serviceGRPCConnectionManager.userServiceGRPCConnection.Close()
	serviceGRPCConnectionManager.catalogServiceGRPCConnection.Close()
	serviceGRPCConnectionManager.elasticsearchServiceGRPCConnection.Close()
}

var ServiceGRPCConnectionManager *serviceGRPCConnectionManager

func InitAllServiceGRPCClients() {
	ServiceGRPCConnectionManager = &serviceGRPCConnectionManager{}

	// Kết nối user-service
	go func() {
		userServiceGRPCServerAddress := net.JoinHostPort(config.AppConfig.UserServiceGRPCHost, config.AppConfig.UserServiceGRPCPort)
		for {
			testingConn, err := net.DialTimeout("tcp", userServiceGRPCServerAddress, 2*time.Second)
			if err == nil {
				testingConn.Close()

				userServiceGRPCServerAddress = fmt.Sprintf(
					"%s:%s",
					config.AppConfig.UserServiceGRPCHost,
					config.AppConfig.UserServiceGRPCPort,
				)

				conn, err := grpc.NewClient(userServiceGRPCServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
				if err != nil {
					log.Fatalf("connect to user-service failed: %s", err.Error())
				}
				ServiceGRPCConnectionManager.userServiceGRPCConnection = conn
				UserServiceGRPCClient = userservicepb.NewUserServiceGRPCClient(conn)

				log.Printf("Connect to user-service successful")

				return
			}

			log.Printf("Waiting for user-service (%s) to be ready...", userServiceGRPCServerAddress)
			time.Sleep(1 * time.Second)
		}
	}()

	// Kết nối catalog-service
	go func() {
		catalogServiceGRPCServerAddress := net.JoinHostPort(config.AppConfig.CatalogServiceGRPCHost, config.AppConfig.CatalogServiceGRPCPort)
		for {
			testingConn, err := net.DialTimeout("tcp", catalogServiceGRPCServerAddress, 2*time.Second)
			if err == nil {
				testingConn.Close()

				catalogServiceGRPCServerAddress = fmt.Sprintf(
					"%s:%s",
					config.AppConfig.CatalogServiceGRPCHost,
					config.AppConfig.CatalogServiceGRPCPort,
				)

				conn, err := grpc.NewClient(catalogServiceGRPCServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
				if err != nil {
					log.Fatalf("connect to catalog-service failed: %s", err.Error())
				}
				ServiceGRPCConnectionManager.catalogServiceGRPCConnection = conn
				CatalogServiceGRPCClient = catalogservicepb.NewCatalogServiceGRPCClient(conn)

				log.Printf("Connect to catalog-service successful")

				return
			}

			log.Printf("Waiting for catalog-service (%s) to be ready...", catalogServiceGRPCServerAddress)
			time.Sleep(1 * time.Second)
		}
	}()

	// Kết nối elasticsearch-service
	go func() {
		elasticsearchServiceGRPCServerAddress := net.JoinHostPort(config.AppConfig.ElasticsearchServiceGRPCHost, config.AppConfig.ElasticsearchServiceGRPCPort)
		for {
			testingConn, err := net.DialTimeout("tcp", elasticsearchServiceGRPCServerAddress, 2*time.Second)
			if err == nil {
				testingConn.Close()

				elasticsearchServiceGRPCServerAddress = fmt.Sprintf(
					"%s:%s",
					config.AppConfig.ElasticsearchServiceGRPCHost,
					config.AppConfig.ElasticsearchServiceGRPCPort,
				)

				conn, err := grpc.NewClient(elasticsearchServiceGRPCServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
				if err != nil {
					log.Fatalf("connect to elasticsearch-service failed: %s", err.Error())
				}
				ServiceGRPCConnectionManager.elasticsearchServiceGRPCConnection = conn
				ElasticsearchServiceGRPCClient = elasticsearchservicepb.NewElasticsearchServiceGRPCClient(conn)

				log.Printf("Connect to elasticsearch-service successful")

				return
			}

			log.Printf("Waiting for elasticsearch-service (%s) to be ready...", elasticsearchServiceGRPCServerAddress)
			time.Sleep(1 * time.Second)
		}
	}()
}
