package utils

import "strings"

type SortField struct {
	Field     string
	Direction string
}

func ParseSorter(sortBy string) []SortField {
	sortFields := []SortField{}

	items := strings.Split(sortBy, ",")
	for _, item := range items {
		var field string
		var direction string

		if strings.Contains(item, ":") {
			parts := strings.Split(item, ":")
			field = parts[0]
			direction = strings.ToUpper(parts[1])
		} else {
			field = item
			direction = "ASC"
		}

		sortFields = append(sortFields, SortField{
			Field:     field,
			Direction: direction,
		})
	}

	return sortFields
}
