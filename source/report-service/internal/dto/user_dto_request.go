package dto

type StatisticsNumberOfUsersCreatedRequest struct {
	CreatedAtGTE string `query:"created_at_gte" default:"2024-01-01T00:00:00" example:"2024-01-15T00:00:00" doc:"Filter by created_at greater than or equal, with format is YYYY-MM-ddTHH:mm:ss."`
	CreatedAtLTE string `query:"created_at_lte" default:"2025-01-01T00:00:00" example:"2024-03-15T23:59:59" doc:"Filter by created_at less than or equal, with format is YYYY-MM-ddTHH:mm:ss."`
	TimeInterval string `query:"time_interval" default:"day" enum:"hour,day,week,month" doc:"Statistics by time interval. Available values: hour,day,week,month"`
}
