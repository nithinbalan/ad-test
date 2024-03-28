package utility

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"strings"

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
