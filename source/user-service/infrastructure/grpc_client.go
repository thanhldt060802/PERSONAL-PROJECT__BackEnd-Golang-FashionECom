package infrastructure

import (
	"fmt"
	"log"
	"net"
	"thanhldt060802/config"
	"thanhldt060802/internal/grpc/client/elasticsearchservicepb"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var ElasticsearchServiceGRPCClient elasticsearchservicepb.ElasticsearchServiceGRPCClient

type serviceGRPCConnectionManager struct {
	elasticsearchServiceGRPCConnection *grpc.ClientConn
}

func (serviceGRPCConnectionManager *serviceGRPCConnectionManager) CloseAll() {
	serviceGRPCConnectionManager.elasticsearchServiceGRPCConnection.Close()
}

var ServiceGRPCConnectionManager *serviceGRPCConnectionManager

func InitAllServiceGRPCClients() {
	ServiceGRPCConnectionManager = &serviceGRPCConnectionManager{}

	// Kết nối elasticsearch-service
	go func() {
		elasticsearchServiceGRPCServerAddress := net.JoinHostPort(config.AppConfig.ElasticSearchServiceGRPCHost, config.AppConfig.ElasticSearchServiceGRPCPort)
		for {
			testingConn, err := net.DialTimeout("tcp", elasticsearchServiceGRPCServerAddress, 2*time.Second)
			if err == nil {
				testingConn.Close()

				elasticsearchServiceGRPCServerAddress = fmt.Sprintf(
					"%s:%s",
					config.AppConfig.ElasticSearchServiceGRPCHost,
					config.AppConfig.ElasticSearchServiceGRPCPort,
				)

				conn, err := grpc.NewClient(elasticsearchServiceGRPCServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
				if err != nil {
					log.Fatalf("connect to elasticsearch-service failed: %s", err.Error())
				}
				ServiceGRPCConnectionManager.elasticsearchServiceGRPCConnection = conn
				ElasticsearchServiceGRPCClient = elasticsearchservicepb.NewElasticsearchServiceGRPCClient(conn)

				log.Printf("connect to elasticsearch-service successful")

				return
			}

			log.Printf("Waiting for user-service (%s) to be ready...", elasticsearchServiceGRPCServerAddress)
			time.Sleep(1 * time.Second)
		}
	}()
}
