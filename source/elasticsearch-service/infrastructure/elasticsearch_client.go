package infrastructure

import (
	"fmt"
	"log"
	"thanhldt060802/config"

	"github.com/elastic/go-elasticsearch/v8"
)

var ElasticsearchClient *elasticsearch.Client

func InitElasticsearchClient() {
	address := fmt.Sprintf(
		"http://%s:%s",
		config.AppConfig.ElasticsearchHost,
		config.AppConfig.ElasticsearchPort,
	)

	elasticsearchClient, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{
			address,
		},
		Username: config.AppConfig.ElasticsearchUsername,
		Password: config.AppConfig.ElasticsearchPassword,
	})
	if err != nil {
		log.Fatal("Connect to Elasticsearch failed: ", err)
	}

	ElasticsearchClient = elasticsearchClient

	res, err := ElasticsearchClient.Info()
	if err != nil {
		log.Fatal("Ping to Elasticsearch failed: ", err)
	}
	defer res.Body.Close()

	log.Println("Connect to Elasticsearch successful")
}
