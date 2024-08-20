package utility

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func ParseQueryParamToInt(c *gin.Context, paramName string, defaultValue int) (int, error) {
	paramValue := c.Query(paramName)
	if paramValue == "" {
		return defaultValue, nil
	}
	return strconv.Atoi(paramValue)
}

func ParseQueryFilters(allowedParams []string, queryParams url.Values) map[string][]string {
	filters := make(map[string][]string)

	for key, values := range queryParams {
		for _, allowedParam := range allowedParams {
			if key == allowedParam {
				// Extract the column name by removing the trailing 's'
				columnName := strings.TrimSuffix(key, "s")
				if len(columnName) > 0 {
					// Split the values by comma
					valueList := strings.Split(values[0], ",")
					// Add the filtering condition if the column name is not empty
					filters[columnName] = valueList
				}
				break // Break the inner loop once a match is found
			}
		}
	}

	return filters
}

func ParseCommaSeparatedStringToArray(input string) []string {
	// Split the string by commas
	items := strings.Split(input, ",")

	// Create a slice to hold the trimmed and non-empty items
	var result []string
	for _, item := range items {
		// Trim whitespace and check if the item is non-empty
		trimmedItem := strings.TrimSpace(item)
		if trimmedItem != "" {
			result = append(result, trimmedItem)
		}
	}

	return result
}

func ParseCommaSeparatedStringToIntArray(input string) ([]int, error) {
	// Split the string by commas
	items := strings.Split(input, ",")

	// Create a slice to hold the converted integers
	var result []int
	for _, item := range items {
		// Trim whitespace and check if the item is non-empty
		trimmedItem := strings.TrimSpace(item)
		if trimmedItem != "" {
			intValue, err := strconv.Atoi(trimmedItem)
			if err != nil {
				return nil, err
			}
			result = append(result, intValue)
		}
	}

	return result, nil
}
