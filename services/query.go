package services

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseQueryFilters(raw string) (map[string]interface{}, error) {
	filters := map[string]interface{}{}
	if strings.TrimSpace(raw) == "" {
		return filters, nil
	}

	parts := strings.Split(raw, ",")
	for _, part := range parts {
		chunks := strings.SplitN(strings.TrimSpace(part), ":", 2)
		if len(chunks) != 2 {
			return nil, fmt.Errorf("invalid query filter %q", part)
		}

		key := strings.TrimSpace(chunks[0])
		value := strings.TrimSpace(chunks[1])
		if key == "" || value == "" {
			return nil, fmt.Errorf("invalid query filter %q", part)
		}

		filters[key] = coerceValue(value)
	}

	return filters, nil
}

func coerceValue(value string) interface{} {
	if parsedInt, err := strconv.ParseInt(value, 10, 64); err == nil {
		return parsedInt
	}

	if parsedBool, err := strconv.ParseBool(value); err == nil {
		return parsedBool
	}

	return value
}
