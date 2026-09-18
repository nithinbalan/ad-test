package utility

import (
	"fmt"
	"strconv"
	"strings"
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

// DateTimeFormatter combines a user's timezone with the configured system date/time format.
type DateTimeFormatter struct {
	Timezone   *time.Location
	DateFormat string
	TimeFormat string
}

// NewDateTimeFormatter creates a formatter from the user's timezone name and the system date/time format.
func NewDateTimeFormatter(tzName, dateFormat, timeFormat string) *DateTimeFormatter {
	f := &DateTimeFormatter{
		DateFormat: dateFormat,
		TimeFormat: timeFormat,
	}
	if tzName != "" {
		if loc, err := time.LoadLocation(tzName); err == nil {
			f.Timezone = loc
		}
	}
	return f
}

// dayjsToGoLayout converts a subset of dayjs-style date format tokens to Go's reference layout.
func dayjsToGoLayout(format string) string {
	if format == "" {
		return "2006-01-02"
	}
	r := strings.NewReplacer(
		"YYYY", "2006",
		"YY", "06",
		"MMMM", "January",
		"MMM", "Jan",
		"MM", "01",
		"M", "1",
		"DD", "02",
		"D", "2",
		"dddd", "Monday",
		"ddd", "Mon",
	)
	return r.Replace(format)
}

func (f *DateTimeFormatter) dateLayout() string {
	return dayjsToGoLayout(f.DateFormat)
}

func (f *DateTimeFormatter) timeLayout() string {
	if f.TimeFormat == "24" {
		return "15:04"
	}
	return "3:04 PM"
}

func (f *DateTimeFormatter) dateTimeLayout() string {
	return f.dateLayout() + " " + f.timeLayout()
}

// FormatDate formats the value using only the date part, without timezone conversion.
func (f *DateTimeFormatter) FormatDate(t *time.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}
	return t.Format(f.dateLayout())
}

// FormatDateTime converts the value to the user's timezone and then formats it with date and time.
func (f *DateTimeFormatter) FormatDateTime(t *time.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}
	value := *t
	if f.Timezone != nil {
		value = value.In(f.Timezone)
	}
	return value.Format(f.dateTimeLayout())
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
