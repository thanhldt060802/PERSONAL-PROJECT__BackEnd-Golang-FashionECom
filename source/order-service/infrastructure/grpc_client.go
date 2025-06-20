package infrastructure

import (
	"fmt"
	"log"
	"net"
	"thanhldt060802/config"
	"thanhldt060802/internal/grpc/client/catalogservicepb"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var CatalogServiceGRPCClient catalogservicepb.CatalogServiceGRPCClient

type serviceGRPCConnectionManager struct {
	catalogServiceGRPCConnection *grpc.ClientConn
}

func (serviceGRPCConnectionManager *serviceGRPCConnectionManager) CloseAll() {
	serviceGRPCConnectionManager.catalogServiceGRPCConnection.Close()
}

var ServiceGRPCConnectionManager *serviceGRPCConnectionManager

func InitAllServiceGRPCClients() {
	ServiceGRPCConnectionManager = &serviceGRPCConnectionManager{}

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
					log.Fatalf("connect to elasticsearch-service failed: %s", err.Error())
				}
				ServiceGRPCConnectionManager.catalogServiceGRPCConnection = conn
				CatalogServiceGRPCClient = catalogservicepb.NewCatalogServiceGRPCClient(conn)

				log.Printf("Connect to elasticsearch-service successful")

				return
			}

			log.Printf("Waiting for catalog-service (%s) to be ready...", catalogServiceGRPCServerAddress)
			time.Sleep(1 * time.Second)
		}
	}()

	// Kết nối elasticsearch-service
}
