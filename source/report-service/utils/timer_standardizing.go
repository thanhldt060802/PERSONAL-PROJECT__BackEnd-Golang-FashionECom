package utils

import (
	"fmt"
	"time"
)

func GenerateEndTimeString(startTimeStr string, interval string) (string, error) {
	sampleLayout := "2006-01-02T15:04:05"

	startTime, err := time.Parse(sampleLayout, startTimeStr)
	if err != nil {
		return "", err
	}

	var endTimeStr string
	switch interval {
	case "hour":
		endTimeStr = startTime.Add(time.Hour).Format(sampleLayout)
	case "day":
		endTimeStr = startTime.AddDate(0, 0, 1).Format(sampleLayout)
	case "week":
		endTimeStr = startTime.AddDate(0, 0, 7).Format(sampleLayout)
	case "month":
		endTimeStr = startTime.AddDate(0, 1, 0).Format(sampleLayout)
	default:
		return "", fmt.Errorf("interval unit is not valid")
	}

	return endTimeStr, nil
}
