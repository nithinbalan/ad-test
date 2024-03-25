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
