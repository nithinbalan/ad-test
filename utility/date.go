package utility

import "time"

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

	// Parse the time string with milliseconds (".0000000" for time(7))
	parsedTime, err := time.Parse("15:04:05.9999999", timeStr)
	if err != nil {
		return time.Time{}, err
	}

	return parsedTime, nil
}
