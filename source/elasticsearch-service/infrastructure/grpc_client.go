package infrastructure

import (
	"fmt"
	"log"
	"net"
	"thanhldt060802/config"
	"thanhldt060802/internal/grpc/client/catalogservicepb"
	"thanhldt060802/internal/grpc/client/orderservicepb"
	"thanhldt060802/internal/grpc/client/userservicepb"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var UserServiceGRPCConnection *grpc.ClientConn
var CatalogServiceGRPCConnection *grpc.ClientConn
var OrderServiceGRPCConnection *grpc.ClientConn

var UserServiceGRPCClient userservicepb.UserServiceGRPCClient
var CatalogServiceGRPCClient catalogservicepb.CatalogServiceGRPCClient
var OrderServiceGRPCClient orderservicepb.OrderServiceGRPCClient

var UserServiceGRPCClientConnectionEvent chan struct{} = make(chan struct{}, 1)
var CatalogServiceGRPCClientConnectionEvent chan struct{} = make(chan struct{}, 1)
var OrderServiceGRPCClientConnectionEvent chan struct{} = make(chan struct{}, 1)

func InitAllServiceGRPCClients() {
	// Kết nối user-service
	if config.AppConfig.SyncAvailableDataFromUserService == "true" {
		go func() {
			userServiceGRPCServerAddress := net.JoinHostPort(config.AppConfig.UserServiceGRPCHost, config.AppConfig.UserServiceGRPCPort)
			for {
				testingConn, err := net.DialTimeout("tcp", userServiceGRPCServerAddress, 2*time.Second)
				if err == nil {
					testingConn.Close()

					userServiceGRPCServerAddress := fmt.Sprintf(
						"%s:%s",
						config.AppConfig.UserServiceGRPCHost,
						config.AppConfig.UserServiceGRPCPort,
					)

					conn, err := grpc.NewClient(userServiceGRPCServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
					if err != nil {
						log.Fatalf("connect to user-service failed: %s", err.Error())
					}
					UserServiceGRPCConnection = conn
					UserServiceGRPCClient = userservicepb.NewUserServiceGRPCClient(conn)

					log.Printf("Connect to user-service successful")

					UserServiceGRPCClientConnectionEvent <- struct{}{}

					return
				}

				log.Printf("Waiting for user-service (%s) to be ready...", userServiceGRPCServerAddress)
				time.Sleep(1 * time.Second)
			}
		}()
	}

	// Kết nối catalog-service
	if config.AppConfig.SyncAvailableDataFromCatalogService == "true" {
		go func() {
			catalogServiceGRPCServerAddress := net.JoinHostPort(config.AppConfig.CatalogServiceGRPCHost, config.AppConfig.CatalogServiceGRPCPort)
			for {
				testingConn, err := net.DialTimeout("tcp", catalogServiceGRPCServerAddress, 2*time.Second)
				if err == nil {
					testingConn.Close()

					catalogServiceGRPCServerAddress := fmt.Sprintf(
						"%s:%s",
						config.AppConfig.CatalogServiceGRPCHost,
						config.AppConfig.CatalogServiceGRPCPort,
					)

					conn, err := grpc.NewClient(catalogServiceGRPCServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
					if err != nil {
						log.Fatalf("connect to catalog-service failed: %s", err.Error())
					}
					CatalogServiceGRPCConnection = conn
					CatalogServiceGRPCClient = catalogservicepb.NewCatalogServiceGRPCClient(conn)

					log.Printf("Connect to catalog-service successful")

					CatalogServiceGRPCClientConnectionEvent <- struct{}{}

					return
				}

				log.Printf("Waiting for catalog-service (%s) to be ready...", catalogServiceGRPCServerAddress)
				time.Sleep(1 * time.Second)
			}
		}()
	}

	// Kết nối order-service
	if config.AppConfig.SyncAvailableDataFromOrderService == "true" {
		go func() {
			orderServiceGRPCServerAddress := net.JoinHostPort(config.AppConfig.OrderServiceGRPCHost, config.AppConfig.OrderServiceGRPCPort)
			for {
				testingConn, err := net.DialTimeout("tcp", orderServiceGRPCServerAddress, 2*time.Second)
				if err == nil {
					testingConn.Close()

					orderServiceGRPCServerAddress := fmt.Sprintf(
						"%s:%s",
						config.AppConfig.OrderServiceGRPCHost,
						config.AppConfig.OrderServiceGRPCPort,
					)

					conn, err := grpc.NewClient(orderServiceGRPCServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
					if err != nil {
						log.Fatalf("connect to order-service failed: %s", err.Error())
					}
					OrderServiceGRPCConnection = conn
					OrderServiceGRPCClient = orderservicepb.NewOrderServiceGRPCClient(conn)

					log.Printf("Connect to order-service successful")

					OrderServiceGRPCClientConnectionEvent <- struct{}{}

					return
				}

				log.Printf("Waiting for order-service (%s) to be ready...", orderServiceGRPCServerAddress)
				time.Sleep(1 * time.Second)
			}
		}()
	}
}
