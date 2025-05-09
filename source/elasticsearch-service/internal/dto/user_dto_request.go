package dto

type SyncAllAvailableUsersRequest struct {
	Body struct {
		URL string `json:"url" required:"true" minLength:"1" doc:"URL to user-service for getting all users."`
	}
}

type GetUsersRequest struct {
	Offset       int    `query:"offset" default:"0" minimum:"0" example:"0" doc:"Skip item by offset."`
	Limit        int    `query:"limit" default:"5" minimum:"1" maximum:"10" example:"10" doc:"Limit item from offset."`
	SortBy       string `query:"sort_by" default:"created_at:asc" example:"full_name:desc,created_at" doc:"Sort by one or more fields separated by commas. For example: sort_by=full_name:desc,created_at will sort by full_name in descending order, then by created_at in ascending order."`
	FullName     string `query:"full_name" example:"Thành Lê" doc:"Filter by full name."`
	Email        string `query:"email" example:"thanhle" doc:"Filter by email."`
	Username     string `query:"username" example:"thanhle" doc:"Filter by username."`
	Address      string `query:"address" example:"Quận 7, Hồ Chí Minh" doc:"Filter by address."`
	RoleName     string `query:"role_name" enum:"ADMIN,CUSTOMER" example:"CUSTOMER" doc:"Filter by role name."`
	CreatedAtGTE string `query:"created_at_gte" example:"2024-01-15T00:00:00" doc:"Filter by created_at greater than or equal, with format is YYYY-MM-ddTHH:mm:ss."`
	CreatedAtLTE string `query:"created_at_lte" example:"2024-02-05T23:59:59" doc:"Filter by created_at less than or equal, with format is YYYY-MM-ddTHH:mm:ss."`
}

type StatisticsNumberOfUsersCreatedRequest struct {
	CreatedAtGTE string `query:"created_at_gte" default:"2024-01-01T00:00:00" example:"2024-01-15T00:00:00" doc:"Filter by created_at greater than or equal, with format is YYYY-MM-ddTHH:mm:ss."`
	CreatedAtLTE string `query:"created_at_lte" default:"2025-01-01T00:00:00" example:"2024-03-15T23:59:59" doc:"Filter by created_at less than or equal, with format is YYYY-MM-ddTHH:mm:ss."`
	TimeInterval string `query:"time_interval" default:"day" enum:"hour,day,week,month" doc:"Statistics by time interval. Available values: hour,day,week,month"`
}
