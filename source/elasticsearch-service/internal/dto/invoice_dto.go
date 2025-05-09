package dto

type RevenueReport struct {
	Status       string  `json:"status,omitempty"`
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
