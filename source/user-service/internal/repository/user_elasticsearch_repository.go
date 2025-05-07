package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"thanhldt060802/infrastructure"
	"thanhldt060802/internal/dto"
	"thanhldt060802/internal/model"
	"thanhldt060802/utils"

	"github.com/elastic/go-elasticsearch/v8/esutil"
)

type userElasticsearchRepository struct {
}

type UserElasticsearchRepository interface {
	SyncAllAvailableUsers(ctx context.Context, users []model.User) error

	SyncCreatingUser(ctx context.Context, newUser *model.User) error
	SyncUpdatingUser(ctx context.Context, updatedUser *model.User) error
	SyncDeletingUserById(ctx context.Context, id int64) error

	GetUsers(ctx context.Context, offset int, limit int, sortFields []utils.SortField,
		fullName string,
		email string,
		username string,
		address string,
		roleName string,
		createdAtGTE string,
		createdAtLTE string,
	) ([]dto.UserView, error)
	StatisticsNumberOfUsersCreated(ctx context.Context,
		calendarInterval string,
		createdAtGTE string,
		createdAtLTE string,
	) (*dto.NumberOfUsersCreatedReport, error)
}

func NewUserElasticsearchRepository() UserElasticsearchRepository {
	return &userElasticsearchRepository{}
}

func (userElasticsearchRepository *userElasticsearchRepository) SyncAllAvailableUsers(ctx context.Context, users []model.User) error {
	// Check if index already exists on Elasticsearch
	existsRes, err := infrastructure.ElasticsearchClient.Indices.Exists([]string{"users"})
	if err != nil {
		return err
	}
	defer existsRes.Body.Close()

	// If index does not exists on Elasticsearch
	if existsRes.StatusCode == 404 {
		// Create index on Elasticsearch using custom schema
		createRes, err := infrastructure.ElasticsearchClient.Indices.Create("users",
			infrastructure.ElasticsearchClient.Indices.Create.WithBody(bytes.NewReader([]byte(model.UserSchemaMappingElasticsearch))))
		if err != nil {
			return err
		}
		defer createRes.Body.Close()

		if createRes.IsError() {
			return fmt.Errorf("%s", createRes.String())
		}

		// Create BulkIndexer for above index
		indexer, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
			Client: infrastructure.ElasticsearchClient,
			Index:  "users",
		})
		if err != nil {
			return err
		}
		defer indexer.Close(ctx)

		// Add all available data on PostgreSQL to BulkIndexer
		for _, user := range users {
			// Convert data to JSON data
			userJSON, err := json.Marshal(user)
			if err != nil {
				return fmt.Errorf("marshal user with id = %d failed: %s", user.Id, err.Error())
			}

			// Add data to BulkIndexer
			err = indexer.Add(ctx, esutil.BulkIndexerItem{
				Action:     "index",
				DocumentID: strconv.FormatInt(user.Id, 10),
				Body:       bytes.NewReader(userJSON),
				OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, resp esutil.BulkIndexerResponseItem, err error) {
					if err != nil {
						log.Printf("Bulk index failed: %s", err.Error())
					} else {
						log.Printf("Index user with id = %s failed: %s", item.DocumentID, resp.Error.Reason)
					}
				},
			})
			if err != nil {
				return err
			}
		}

		return nil
	}

	return fmt.Errorf("users index already exists after first sync all")
}

func (userElasticsearchRepository *userElasticsearchRepository) SyncCreatingUser(ctx context.Context, newUser *model.User) error {
	// Add data to Elasticsearch
	res, err := infrastructure.ElasticsearchClient.Index(
		"users",
		esutil.NewJSONReader(newUser),
		infrastructure.ElasticsearchClient.Index.WithDocumentID(strconv.FormatInt(newUser.Id, 10)),
		infrastructure.ElasticsearchClient.Index.WithRefresh("true"),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("%s", res.String())
	}

	return nil
}

func (userElasticsearchRepository *userElasticsearchRepository) SyncUpdatingUser(ctx context.Context, updatedUser *model.User) error {
	// Update data on Elasticsearch
	res, err := infrastructure.ElasticsearchClient.Index(
		"users",
		esutil.NewJSONReader(updatedUser),
		infrastructure.ElasticsearchClient.Index.WithDocumentID(strconv.FormatInt(updatedUser.Id, 10)),
		infrastructure.ElasticsearchClient.Index.WithRefresh("true"),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("%s", res.String())
	}

	return nil
}

func (userElasticsearchRepository *userElasticsearchRepository) SyncDeletingUserById(ctx context.Context, id int64) error {
	// Delete data from Elasticsearch
	res, err := infrastructure.ElasticsearchClient.Delete(
		"users",
		strconv.FormatInt(id, 10),
		infrastructure.ElasticsearchClient.Delete.WithRefresh("true"),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("%s", res.String())
	}

	return nil
}

func (userElasticsearchRepository *userElasticsearchRepository) GetUsers(ctx context.Context, offset int, limit int, sortFields []utils.SortField,
	fullName string,
	email string,
	username string,
	address string,
	roleName string,
	createdAtGTE string,
	createdAtLTE string,
) ([]dto.UserView, error) {
	mustConditions := []map[string]interface{}{}

	// If filtering by full_name
	if fullName != "" {
		mustConditions = append(mustConditions, map[string]interface{}{
			"match": map[string]interface{}{
				"full_name": fullName,
			},
		})
	}

	// If filtering by email
	if email != "" {
		mustConditions = append(mustConditions, map[string]interface{}{
			"match": map[string]interface{}{
				"email": email,
			},
		})
	}

	// If filtering by username
	if username != "" {
		mustConditions = append(mustConditions, map[string]interface{}{
			"match": map[string]interface{}{
				"username": username,
			},
		})
	}

	// If filtering by address
	if address != "" {
		mustConditions = append(mustConditions, map[string]interface{}{
			"match": map[string]interface{}{
				"address": address,
			},
		})
	}

	// If filtering by role_name
	if roleName != "" {
		mustConditions = append(mustConditions, map[string]interface{}{
			"match": map[string]interface{}{
				"role_name": roleName,
			},
		})
	}

	// If filtering by created_at in range or partial range
	createdAtRange := map[string]interface{}{}
	if createdAtGTE != "" {
		createdAtRange["gte"] = createdAtGTE
	}
	if createdAtLTE != "" {
		createdAtRange["lte"] = createdAtLTE
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
		"from": offset,
		"size": limit,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": mustConditions,
			},
		},
	}

	// Apply sorting to query
	_sortFields := []map[string]interface{}{}
	for _, sortField := range sortFields {
		_sortFields = append(_sortFields, map[string]interface{}{
			model.UserSchemaMappingElasticsearchSortFieldMap[sortField.Field]: strings.ToLower(sortField.Direction),
		})
	}
	query["sort"] = _sortFields

	// Convert query to JSON query
	queryJSON, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("marshal query failed: %s", err.Error())
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
		return nil, fmt.Errorf("%s", res.String())
	}

	// Declare Elasticsearch response
	var elasticsearchResponse struct {
		Hits struct {
			Hits []struct {
				Source dto.UserView `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	// Unmarshal Elasticsearch response body to Elasticsearch response
	elasticsearchResponseBody := json.NewDecoder(res.Body)
	if err := elasticsearchResponseBody.Decode(&elasticsearchResponse); err != nil {
		return nil, fmt.Errorf("unmarshal elasticsearch response failed: %s", err.Error())
	}

	// Extract data from Elasticsearch response
	users := make([]dto.UserView, len(elasticsearchResponse.Hits.Hits))
	for i, hit := range elasticsearchResponse.Hits.Hits {
		users[i] = hit.Source
	}

	return users, nil
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
			"users_per_interval": map[string]interface{}{
				"date_histogram": map[string]interface{}{
					"field":             "created_at",
					"calendar_interval": calendarInterval,
					"format":            "yyyy-MM-dd'T'HH:mm:ss",
					// "time_zone":         "Asia/Ho_Chi_Minh",
				},
			},
			"total_users": map[string]interface{}{
				"value_count": map[string]interface{}{
					"field": "id",
				},
			},
			"avg_users_per_interval": map[string]interface{}{
				"avg_bucket": map[string]interface{}{
					"buckets_path": "users_per_interval>_count",
				},
			},
		},
	}

	// Convert query to JSON query
	queryJSON, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("marshal query failed: %s", err.Error())
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
		return nil, fmt.Errorf("%s", res.String())
	}

	// Declare Elasticsearch response
	var elasticsearchResponse struct {
		Aggregations struct {
			UsersPerInterval struct {
				Buckets []struct {
					KeyAsString string `json:"key_as_string"`
					DocCount    int64  `json:"doc_count"`
				} `json:"buckets"`
			} `json:"users_per_interval"`
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
		return nil, fmt.Errorf("unmarshal elasticsearch response failed: %s", err.Error())
	}

	// Extract data from Elasticsearch response
	report.Total = elasticsearchResponse.Aggregations.TotalUsers.Value
	report.Average = elasticsearchResponse.Aggregations.AvgUsersPerInterval.Value
	for _, bucket := range elasticsearchResponse.Aggregations.UsersPerInterval.Buckets {
		startTime := bucket.KeyAsString
		endTime, err := utils.AddInterval(startTime, calendarInterval)
		if err != nil {
			return nil, fmt.Errorf("standardize end time failed: %s", err.Error())
		}

		report.Details = append(report.Details, struct {
			StartTime string  `json:"start_time"`
			EndTime   string  `json:"end_time"`
			Total     float64 `json:"total"`
		}{
			StartTime: startTime,
			EndTime:   *endTime,
			Total:     float64(bucket.DocCount),
		})
	}

	return report, nil
}
