package dto

import "time"

type UserView struct {
	Id        int64     `json:"id"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Address   string    `json:"address"`
	RoleName  string    `json:"role_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type NumberOfUsersCreatedReport struct {
	StartTime    string  `json:"start_time"`
	EndTime      string  `json:"end_time"`
	TimeInterval string  `json:"time_interval"`
	Total        float64 `json:"total"`
	Average      float64 `json:"average"`
	Details      []struct {
		StartTime string  `json:"start_time"`
		EndTime   string  `json:"end_time"`
		Total     float64 `json:"total"`
	} `json:"detail"`
}
