package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"thanhldt060802/infrastructure"
	"thanhldt060802/internal/dto"
	"thanhldt060802/utils"
)

type userElasticsearchRepository struct {
}

type UserElasticsearchRepository interface {
	StatisticsNumberOfUsersCreated(ctx context.Context,
		calendarInterval string,
		createdAtGTE string,
		createdAtLTE string,
	) (*dto.NumberOfUsersCreatedReport, error)
}

func NewUserElasticsearchRepository() UserElasticsearchRepository {
	return &userElasticsearchRepository{}
}

func (userElasticsearchRepository *userElasticsearchRepository) StatisticsNumberOfUsersCreated(ctx context.Context,
	calendarInterval string,
	createdAtGTE string,
	createdAtLTE string,
) (*dto.NumberOfUsersCreatedReport, error) {
	report := &dto.NumberOfUsersCreatedReport{}
	report.TimeInterval = calendarInterval

	mustConditions := []map[string]interface{}{}

	// If filtering by created_at in range or partial range
	createdAtRange := map[string]interface{}{}
	if createdAtGTE != "" {
		createdAtRange["gte"] = createdAtGTE
		report.StartTime = createdAtGTE
	}
	if createdAtLTE != "" {
		createdAtRange["lte"] = createdAtLTE
		report.EndTime = createdAtLTE
	}
	if len(createdAtRange) > 0 {
		createdAtRange["format"] = "strict_date_optional_time" // For format YYYY-MM-ddTHH:mm:ss
		mustConditions = append(mustConditions, map[string]interface{}{
			"range": map[string]interface{}{
				"created_at": createdAtRange,
			},
		})
	}

	// If not filtering -> get all
	if len(mustConditions) == 0 {
		mustConditions = append(mustConditions, map[string]interface{}{
			"match_all": map[string]interface{}{},
		})
	}

	// Setup query
	query := map[string]interface{}{
		"size": 0,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": mustConditions,
			},
		},
		"aggs": map[string]interface{}{
			"total_users_per_interval": map[string]interface{}{
				"date_histogram": map[string]interface{}{
					"field":             "created_at",
					"calendar_interval": calendarInterval,
					"format":            "yyyy-MM-dd'T'HH:mm:ss",
				},
			},
			"total_users": map[string]interface{}{
				"value_count": map[string]interface{}{
					"field": "id",
				},
			},
			"avg_users_per_interval": map[string]interface{}{
				"avg_bucket": map[string]interface{}{
					"buckets_path": "total_users_per_interval>_count",
				},
			},
		},
	}

	// Convert query to JSON query
	queryJSON, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	// Send request to Elasticsearch
	res, err := infrastructure.ElasticsearchClient.Search(
		infrastructure.ElasticsearchClient.Search.WithContext(ctx),
		infrastructure.ElasticsearchClient.Search.WithIndex("users"),
		infrastructure.ElasticsearchClient.Search.WithBody(bytes.NewReader(queryJSON)),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Parse Elasticsearch response
	if res.IsError() {
		return nil, fmt.Errorf("some thing wrong when querying users")
	}

	// Declare Elasticsearch response
	var elasticsearchResponse struct {
		Aggregations struct {
			TotalUsersPerInterval struct {
				Buckets []struct {
					KeyAsString string `json:"key_as_string"`
					DocCount    int64  `json:"doc_count"`
				} `json:"buckets"`
			} `json:"total_users_per_interval"`
			TotalUsers struct {
				Value float64 `json:"value"`
			} `json:"total_users"`
			AvgUsersPerInterval struct {
				Value float64 `json:"value"`
			} `json:"avg_users_per_interval"`
		} `json:"aggregations"`
	}

	// Unmarshal Elasticsearch response body to Elasticsearch response
	elasticsearchResponseBody := json.NewDecoder(res.Body)
	if err := elasticsearchResponseBody.Decode(&elasticsearchResponse); err != nil {
		return nil, err
	}

	// Extract data from Elasticsearch response
	report.Total = elasticsearchResponse.Aggregations.TotalUsers.Value
	report.Average = elasticsearchResponse.Aggregations.AvgUsersPerInterval.Value
	for _, bucket := range elasticsearchResponse.Aggregations.TotalUsersPerInterval.Buckets {
		startTime := bucket.KeyAsString
		endTime, err := utils.GenerateEndTimeString(startTime, calendarInterval)
		if err != nil {
			return nil, err
		}

		report.Details = append(report.Details, struct {
			StartTime string  `json:"start_time"`
			EndTime   string  `json:"end_time"`
			Total     float64 `json:"total"`
		}{
			StartTime: startTime,
			EndTime:   endTime,
			Total:     float64(bucket.DocCount),
		})
	}

	return report, nil
}
