package service

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
	"thanhldt060802/internal/grpc/pb"
	"thanhldt060802/internal/schema"
	"thanhldt060802/utils"

	"github.com/elastic/go-elasticsearch/v8/esutil"
)

type userService struct {
	userServiceClient pb.UserServiceClient
}

type UserService interface {
	SyncAllAvailableUsers(ctx context.Context) error

	SyncCreatingUser(ctx context.Context, newUser *dto.UserView) error
	SyncUpdatingUserById(ctx context.Context, updatedUser *dto.UserView) error
	SyncDeletingUserById(ctx context.Context, id int64) error

	GetUsers(ctx context.Context, reqDTO *dto.GetUsersRequest) ([]dto.UserView, error)

	// StatisticsNumberOfUsersCreated(ctx context.Context, reqDTO *dto.StatisticsNumberOfUsersCreatedRequest) (*dto.NumberOfUsersCreatedReport, error)
}

func NewUserService(userServiceClient pb.UserServiceClient) UserService {
	return &userService{
		userServiceClient: userServiceClient,
	}
}

func (userService *userService) SyncAllAvailableUsers(ctx context.Context) error {
	// Check if index already exists on Elasticsearch
	existsRes, err := infrastructure.ElasticsearchClient.Indices.Exists([]string{"users"})
	if err != nil {
		return err
	}
	defer existsRes.Body.Close()

	// If index does not exists on Elasticsearch
	if existsRes.StatusCode == 404 {
		// // Get all available users from user-service
		// req, err := http.NewRequest("GET", "asdasd", nil)
		// if err != nil {
		// 	return err
		// }
		// req.Header.Set("Authorization", "Bearer "+ctx.Value("access_token").(string))

		// resp, err := http.DefaultClient.Do(req)
		// if err != nil {
		// 	return err
		// }
		// defer resp.Body.Close()

		// if resp.StatusCode != http.StatusOK {
		// 	return fmt.Errorf("get all available users from user-service failed")
		// }

		// body, err := io.ReadAll(resp.Body)
		// if err != nil {
		// 	return err
		// }

		// var respBody dto.BodyResponse[[]dto.UserView]
		// if err := json.Unmarshal(body, &respBody); err != nil {
		// 	return err
		// }

		// users := respBody.Body.Data

		grpcRes, err := userService.userServiceClient.GetAllUsers(ctx, &pb.GetAllUsersRequest{})
		if err != nil {
			return fmt.Errorf("some thing wrong when loading data from user-service")
		}
		users := grpcRes.Users

		// Create index on Elasticsearch using custom schema
		res, err := infrastructure.ElasticsearchClient.Indices.Create("users",
			infrastructure.ElasticsearchClient.Indices.Create.WithBody(bytes.NewReader([]byte(schema.User))))
		if err != nil {
			return err
		}
		defer res.Body.Close()

		if res.IsError() {
			return fmt.Errorf("some thing wrong when creating users index on elasticsearch")
		}

		hasFailure := false
		var closeBulkIndexer error

		// Create BulkIndexer for above index to index to Elasticsearch
		indexer, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
			Client: infrastructure.ElasticsearchClient,
			Index:  "users",
		})
		if err != nil {
			return err
		}
		defer func() {
			if err := indexer.Close(ctx); err != nil {
				closeBulkIndexer = fmt.Errorf("close bulk indexer failed: %s", err.Error())
			}
		}()

		// Add all available data on PostgreSQL to BulkIndexer
		for _, user := range users {
			// Convert data to JSON data
			userJSON, err := json.Marshal(dto.ToUserViewFromProto(user))
			if err != nil {
				return err
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
					hasFailure = true
				},
			})
			if err != nil {
				return err
			}
		}

		if hasFailure {
			return fmt.Errorf("some thing wrong when syncing all available users to elasticsearch")
		}
		if closeBulkIndexer != nil {
			return fmt.Errorf("some thing wrong when syncing all available users (%s) to elasticsearch", closeBulkIndexer.Error())
		}

		return nil
	}

	return fmt.Errorf("users index on elasticsearch already exists after first syncing all")
}

func (userService *userService) SyncCreatingUser(ctx context.Context, newUser *dto.UserView) error {
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
		return fmt.Errorf("some thing wrong when syncing creating user on elasticsearch")
	}

	return nil
}

func (userService *userService) SyncUpdatingUserById(ctx context.Context, updatedUser *dto.UserView) error {
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
		return fmt.Errorf("some thing wrong when syncing updating user on elasticsearch")
	}

	return nil
}

func (userService *userService) SyncDeletingUserById(ctx context.Context, id int64) error {
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
		return fmt.Errorf("some thing wrong when syncing deleting user from elasticsearch")
	}

	return nil
}

func (userService *userService) GetUsers(ctx context.Context, reqDTO *dto.GetUsersRequest) ([]dto.UserView, error) {
	mustConditions := []map[string]interface{}{}

	// If filtering by full_name
	if reqDTO.FullName != "" {
		mustConditions = append(mustConditions, map[string]interface{}{
			"match": map[string]interface{}{
				"full_name": reqDTO.FullName,
			},
		})
	}

	// If filtering by email
	if reqDTO.Email != "" {
		mustConditions = append(mustConditions, map[string]interface{}{
			"match": map[string]interface{}{
				"email": reqDTO.Email,
			},
		})
	}

	// If filtering by username
	if reqDTO.Username != "" {
		mustConditions = append(mustConditions, map[string]interface{}{
			"match": map[string]interface{}{
				"username": reqDTO.Username,
			},
		})
	}

	// If filtering by address
	if reqDTO.Address != "" {
		mustConditions = append(mustConditions, map[string]interface{}{
			"match": map[string]interface{}{
				"address": reqDTO.Address,
			},
		})
	}

	// If filtering by role_name
	if reqDTO.RoleName != "" {
		mustConditions = append(mustConditions, map[string]interface{}{
			"match": map[string]interface{}{
				"role_name": reqDTO.RoleName,
			},
		})
	}

	// If filtering by created_at in range or partial range
	createdAtRange := map[string]interface{}{}
	if reqDTO.CreatedAtGTE != "" {
		createdAtRange["gte"] = reqDTO.CreatedAtGTE
	}
	if reqDTO.CreatedAtLTE != "" {
		createdAtRange["lte"] = reqDTO.CreatedAtLTE
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
		"from": reqDTO.Offset,
		"size": reqDTO.Limit,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": mustConditions,
			},
		},
	}

	// Apply sorting to query
	sortFields := utils.ParseSorter(reqDTO.SortBy)
	_sortFields := []map[string]interface{}{}
	for _, sortField := range sortFields {
		_sortFields = append(_sortFields, map[string]interface{}{
			schema.UserStandardizeSortFieldMap[sortField.Field]: strings.ToLower(sortField.Direction),
		})
	}
	query["sort"] = _sortFields

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
		return nil, fmt.Errorf("some thing wrong when querying users on elasticsearch")
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
		return nil, err
	}

	// Extract data from Elasticsearch response
	users := make([]dto.UserView, len(elasticsearchResponse.Hits.Hits))
	for i, hit := range elasticsearchResponse.Hits.Hits {
		users[i] = hit.Source
	}

	return users, nil
}

// func (userService *userService) StatisticsNumberOfUsersCreated(ctx context.Context, reqDTO *dto.StatisticsNumberOfUsersCreatedRequest) (*dto.NumberOfUsersCreatedReport, error) {
// 	report := &dto.NumberOfUsersCreatedReport{}
// 	report.TimeInterval = reqDTO.TimeInterval

// 	mustConditions := []map[string]interface{}{}

// 	// If filtering by created_at in range or partial range
// 	createdAtRange := map[string]interface{}{}
// 	if reqDTO.CreatedAtGTE != "" {
// 		createdAtRange["gte"] = reqDTO.CreatedAtGTE
// 		report.StartTime = reqDTO.CreatedAtGTE
// 	}
// 	if reqDTO.CreatedAtLTE != "" {
// 		createdAtRange["lte"] = reqDTO.CreatedAtLTE
// 		report.EndTime = reqDTO.CreatedAtLTE
// 	}
// 	if len(createdAtRange) > 0 {
// 		createdAtRange["format"] = "strict_date_optional_time" // For format YYYY-MM-ddTHH:mm:ss
// 		mustConditions = append(mustConditions, map[string]interface{}{
// 			"range": map[string]interface{}{
// 				"created_at": createdAtRange,
// 			},
// 		})
// 	}

// 	// If not filtering -> get all
// 	if len(mustConditions) == 0 {
// 		mustConditions = append(mustConditions, map[string]interface{}{
// 			"match_all": map[string]interface{}{},
// 		})
// 	}

// 	// Setup query
// 	query := map[string]interface{}{
// 		"size": 0,
// 		"query": map[string]interface{}{
// 			"bool": map[string]interface{}{
// 				"must": mustConditions,
// 			},
// 		},
// 		"aggs": map[string]interface{}{
// 			"total_users_per_interval": map[string]interface{}{
// 				"date_histogram": map[string]interface{}{
// 					"field":             "created_at",
// 					"calendar_interval": report.TimeInterval,
// 					"format":            "yyyy-MM-dd'T'HH:mm:ss",
// 				},
// 			},
// 			"total_users": map[string]interface{}{
// 				"value_count": map[string]interface{}{
// 					"field": "id",
// 				},
// 			},
// 			"avg_users_per_interval": map[string]interface{}{
// 				"avg_bucket": map[string]interface{}{
// 					"buckets_path": "total_users_per_interval>_count",
// 				},
// 			},
// 		},
// 	}

// 	// Convert query to JSON query
// 	queryJSON, err := json.Marshal(query)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Send request to Elasticsearch
// 	res, err := infrastructure.ElasticsearchClient.Search(
// 		infrastructure.ElasticsearchClient.Search.WithContext(ctx),
// 		infrastructure.ElasticsearchClient.Search.WithIndex("users"),
// 		infrastructure.ElasticsearchClient.Search.WithBody(bytes.NewReader(queryJSON)),
// 	)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer res.Body.Close()

// 	// Parse Elasticsearch response
// 	if res.IsError() {
// 		return nil, fmt.Errorf("some thing wrong when querying users")
// 	}

// 	// Declare Elasticsearch response
// 	var elasticsearchResponse struct {
// 		Aggregations struct {
// 			TotalUsersPerInterval struct {
// 				Buckets []struct {
// 					KeyAsString string `json:"key_as_string"`
// 					DocCount    int64  `json:"doc_count"`
// 				} `json:"buckets"`
// 			} `json:"total_users_per_interval"`
// 			TotalUsers struct {
// 				Value float64 `json:"value"`
// 			} `json:"total_users"`
// 			AvgUsersPerInterval struct {
// 				Value float64 `json:"value"`
// 			} `json:"avg_users_per_interval"`
// 		} `json:"aggregations"`
// 	}

// 	// Unmarshal Elasticsearch response body to Elasticsearch response
// 	elasticsearchResponseBody := json.NewDecoder(res.Body)
// 	if err := elasticsearchResponseBody.Decode(&elasticsearchResponse); err != nil {
// 		return nil, err
// 	}

// 	// Extract data from Elasticsearch response
// 	report.Total = elasticsearchResponse.Aggregations.TotalUsers.Value
// 	report.Average = elasticsearchResponse.Aggregations.AvgUsersPerInterval.Value
// 	for _, bucket := range elasticsearchResponse.Aggregations.TotalUsersPerInterval.Buckets {
// 		startTime := bucket.KeyAsString
// 		endTime, err := utils.GenerateEndTimeString(startTime, reqDTO.TimeInterval)
// 		if err != nil {
// 			return nil, err
// 		}

// 		report.Details = append(report.Details, struct {
// 			StartTime string  `json:"start_time"`
// 			EndTime   string  `json:"end_time"`
// 			Total     float64 `json:"total"`
// 		}{
// 			StartTime: startTime,
// 			EndTime:   endTime,
// 			Total:     float64(bucket.DocCount),
// 		})
// 	}

// 	return report, nil
// }
