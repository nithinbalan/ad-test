package utility

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func SetDefault(value string, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

func RemoveLangCodeSuffix(text string) string {
	text = strings.TrimSuffix(text, "En")
	text = strings.TrimSuffix(text, "Ar")
	return strings.ToLower(text)
}

func GetCurrentUserID(c *gin.Context) (uint, error) {
	currentUserID, exists := c.Get("currentUserID")
	if !exists {
		return 0, fmt.Errorf("current user id not found")
	}

	id, ok := currentUserID.(uint)
	if !ok {
		return 0, fmt.Errorf("current user id is of incorrect type")
	}

	return id, nil
}

type ChangeDetail struct {
	Field string `json:"field"`
	From  string `json:"from"`
	To    string `json:"to"`
}

func CompareFields(existing, updated interface{}, fieldsToCompare []string) ([]ChangeDetail, error) {
	var changes []ChangeDetail
	valExisting := reflect.ValueOf(existing)
	valUpdated := reflect.ValueOf(updated)

	// Ensure we are working with structs or pointers to structs
	if valExisting.Kind() == reflect.Ptr {
		valExisting = valExisting.Elem()
	}
	if valUpdated.Kind() == reflect.Ptr {
		valUpdated = valUpdated.Elem()
	}

	// Ensure the values are structs
	if valExisting.Kind() != reflect.Struct || valUpdated.Kind() != reflect.Struct {
		return nil, errors.New("CompareFields: expected struct or pointer to struct")
	}

	for _, field := range fieldsToCompare {
		fieldValueExisting := valExisting.FieldByName(field)
		fieldValueUpdated := valUpdated.FieldByName(field)

		// Check if fields are valid before calling Interface()
		if !fieldValueExisting.IsValid() || !fieldValueUpdated.IsValid() {
			return nil, fmt.Errorf("CompareFields: field '%s' not found in struct", field)
		}

		fromValue := fmt.Sprintf("%v", fieldValueExisting.Interface())
		toValue := fmt.Sprintf("%v", fieldValueUpdated.Interface())

		if fromValue != toValue {
			changes = append(changes, ChangeDetail{
				Field: field,
				From:  fromValue,
				To:    toValue,
			})
		}
	}

	return changes, nil
}

// This function allows adding manual changes, including those for foreign keys or complex fields.
func AddManualChanges(changes []ChangeDetail, manualChanges []ChangeDetail) ([]ChangeDetail, error) {
	changes = append(changes, manualChanges...)
	return changes, nil
}

// This function converts the slice of ChangeDetail to a JSON string.
func ChangesToJSON(changes []ChangeDetail) (string, error) {
	changesJSON, err := json.Marshal(changes)
	if err != nil {
		return "", err
	}
	return string(changesJSON), nil
}

// GenerateOTP generates a random OTP of specified length
func GenerateOTP(length int) string {
	const charset = "0123456789" // Use digits for OTP generation
	otp := make([]byte, length)
	for i := range otp {
		otp[i] = charset[rand.Intn(len(charset))]
	}
	return string(otp)
}

// @todo to be removed
func GeneratePassword() {
	password := "12345678" // Your plaintext password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error hashing password:", err)
		return
	}
	fmt.Println("Hashed Password:", string(hashedPassword))
}

// CheckUniqueFields checks if the specified fields in a slice of struct objects are unique.
// It returns an error if duplicate values are found for any of the specified fields.
func CheckUniqueFields(slice interface{}, fields ...string) error {
	seen := make(map[string]struct{})
	val := reflect.ValueOf(slice)

	if val.Kind() != reflect.Slice {
		return fmt.Errorf("provided value is not a slice")
	}

	for i := 0; i < val.Len(); i++ {
		item := val.Index(i)

		// If the item is a pointer to a struct, get the struct it points to
		if item.Kind() == reflect.Ptr {
			item = item.Elem()
		}

		for _, field := range fields {
			fieldValue := item.FieldByName(field)

			// Ensure the field exists
			if !fieldValue.IsValid() {
				return fmt.Errorf("field %s does not exist in the item", field)
			}

			key := fmt.Sprintf("%s_%v", field, fieldValue.Interface())
			if _, found := seen[key]; found {
				return fmt.Errorf("duplicate value found for field %s", field)
			}
			seen[key] = struct{}{}
		}
	}

	return nil
}

func CalculateNextReviewDate(frequency string, lastDate ...time.Time) time.Time {
	var currentDate time.Time
	if len(lastDate) > 0 {
		currentDate = lastDate[0]
	} else {
		currentDate = time.Now()
	}

	switch frequency {
	case "daily":
		return currentDate.AddDate(0, 0, 1)
	case "weekly":
		return currentDate.AddDate(0, 0, 7)
	case "monthly":
		return currentDate.AddDate(0, 1, 0)
	case "quarterly":
		return currentDate.AddDate(0, 3, 0)
	case "half-yearly":
		return currentDate.AddDate(0, 6, 0)
	case "annually":
		return currentDate.AddDate(1, 0, 0)
	default:
		return currentDate
	}
}

var (
	commonColors = []string{
		"#1976D2", // Dark Blue
		"#388E3C", // Dark Green
		"#F57C00", // Dark Orange
		"#7B1FA2", // Dark Purple
		"#E64A19", // Dark Red
		"#5D4037", // Dark Brown
		"#455A64", // Dark Gray
		"#00838F", // Dark Cyan
		"#7B1FA2", // Dark Indigo
		"#AFB42B", // Dark Lime
		"#03A9F4", // Light Blue
		"#4CAF50", // Light Green
		"#FF9800", // Light Orange
		"#9C27B0", // Light Purple
		"#F44336", // Light Red
		"#795548", // Light Brown
		"#607D8B", // Light Gray
		"#00BCD4", // Light Cyan
		"#3F51B5", // Indigo
		"#CDDC39", // Lime
		"#009688", // Teal
		"#FFEB3B", // Yellow
		"#9E9E9E", // Grey
		"#FF5722", // Deep Orange
		"#2196F3", // Blue
	}

	userColorIndex int
	deptColorIndex int
)

func GenerateColorCode(entityType string) string {
	var colorIndex *int
	switch entityType {
	case "user":
		colorIndex = &userColorIndex
	case "dept":
		colorIndex = &deptColorIndex
	default:
		// Default to user if invalid type is provided
		colorIndex = &userColorIndex
	}

	// Define usedColors map
	usedColors := make(map[string]bool)

	// If all colors are used, reset the cache
	if len(usedColors) == len(commonColors) {
		usedColors = make(map[string]bool)
		*colorIndex = 0
	}

	// Find the next available color
	var color string
	for {
		color = commonColors[*colorIndex]
		if !usedColors[color] {
			usedColors[color] = true
			*colorIndex = (*colorIndex + 1) % len(commonColors)
			break
		}
		*colorIndex = (*colorIndex + 1) % len(commonColors)
	}

	return color
}

func GetAppURL() string {
	appURL := os.Getenv("APP_URL")
	if appURL == "" {
		appURL = "http://grc-dev.ervarioo.com"
	}
	return appURL
}

func GetIntValueAsString(pointer *int) string {
	if pointer != nil {
		return fmt.Sprintf("%d", *pointer)
	}
	return ""
}

func GetStringValue(pointer *string) string {
	if pointer != nil {
		return *pointer
	}
	return ""
}
