package utils

import "time"

func AddInterval(startTimeStr string, interval string) (*string, error) {
	startTime, err := time.Parse("2006-01-02T15:04:05", startTimeStr)
	if err != nil {
		return nil, err
	}

	var endTime string
	switch interval {
	case "hour":
		endTime = startTime.Add(time.Hour).Format("2006-01-02T15:04:05")
	case "day":
		endTime = startTime.AddDate(0, 0, 1).Format("2006-01-02T15:04:05")
	case "week":
		endTime = startTime.AddDate(0, 0, 7).Format("2006-01-02T15:04:05")
	default:
		endTime = startTime.Format("2006-01-02T15:04:05")
	}

	return &endTime, nil
}
