package utility

import (
	"fmt"
	"strconv"
	"time"
)

func ParseDate(dateStr string) (*time.Time, error) {
	if dateStr == "" {
		return nil, nil
	}

	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, err
	}

	return &parsedDate, nil
}

func ParseMandatoryDate(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, nil
	}

	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}, err
	}

	return parsedDate, nil
}

func ParseDateTime(dateTimeStr string) (*time.Time, error) {
	if dateTimeStr == "" {
		return nil, nil
	}

	parsedDateTime, err := time.Parse("2006-01-02T15:04:05", dateTimeStr)
	if err != nil {
		return nil, err
	}

	return &parsedDateTime, nil
}

func ParseMandatoryTime(timeStr string) (time.Time, error) {
	if timeStr == "" {
		return time.Time{}, nil
	}

	parsedTime, err := time.Parse("15:04:05", timeStr)
	if err != nil {
		return time.Time{}, err
	}

	return parsedTime, nil
}

func FormatDate(date *time.Time) string {
	if date != nil {
		return date.Format("2006-01-02") // Formats the date as YYYY-MM-DD
	}
	return ""
}

func ConvertToUTC(localTime string, timezoneName string) (*time.Time, error) {
	// Load the location using the timezone name
	location, err := time.LoadLocation(timezoneName)
	if err != nil {
		return nil, fmt.Errorf("failed to load timezone: %w", err)
	}

	// Parse the local time
	parsedTime, err := time.ParseInLocation("2006-01-02T15:04:05", localTime, location)
	if err != nil {
		return nil, fmt.Errorf("invalid time format: %w", err)
	}

	// Convert to UTC
	utcTime := parsedTime.UTC()
	return &utcTime, nil
}

func ConvertToUTCWithLayouts(localTime string, timezoneName string) (*time.Time, error) {
	// Load the location using the timezone name
	location, err := time.LoadLocation(timezoneName)
	if err != nil {
		return nil, fmt.Errorf("failed to load timezone: %w", err)
	}

	// Try common datetime layouts in order of specificity
	layouts := []string{
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04",
		"2006-01-02 15:04",
		"2006-01-02",
	}

	for _, layout := range layouts {
		parsedTime, err := time.ParseInLocation(layout, localTime, location)
		if err == nil {
			utcTime := parsedTime.UTC()
			return &utcTime, nil
		}
	}

	return nil, fmt.Errorf("invalid time format")
}

func ParseUnixTimestamp(value interface{}) time.Time {
	switch v := value.(type) {
	case float64:
		return time.Unix(int64(v), 0).UTC()
	case string:
		if ts, err := strconv.ParseInt(v, 10, 64); err == nil {
			return time.Unix(ts, 0).UTC()
		}
	}
	return time.Time{}
}
