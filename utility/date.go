package utility

import "time"

func ParseDate(dateStr string) (*time.Time, error) {
	if dateStr == "" {
		return nil, nil // No date provided, return nil
	}

	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, err // Return the parsing error
	}

	return &parsedDate, nil // Return a pointer to the parsed date
}

func ParseDateTime(dateTimeStr string) (*time.Time, error) {
	if dateTimeStr == "" {
		return nil, nil // No date provided, return nil
	}

	parsedDateTime, err := time.Parse("2006-01-02T15:04:05", dateTimeStr)
	if err != nil {
		return nil, err // Return the parsing error
	}

	return &parsedDateTime, nil // Return a pointer to the parsed date and time
}
