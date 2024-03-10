package models

import (
	"github.com/go-playground/validator/v10"
)

type ReviewFrequency struct {
	ID      uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn string `json:"title_en" binding:"required,max=100" gorm:"type:varchar(100);not null"`
	TitleAr string `json:"title_ar" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	Type    string `json:"type" binding:"required,oneof=daily weekly monthly quarterly half-yearly annually" gorm:"type:varchar(45);not null"`
	Status  string `json:"status" binding:"required,oneof=active inactive" gorm:"type:varchar(10);not null"`
}

// TableName overrides the table name used by ReviewFrequency to `review_frequencies`
func (ReviewFrequency) TableName() string {
	return "review_frequencies"
}

// Validation for the 'type' field to check if the value is within the allowed types
func typeValidation(fl validator.FieldLevel) bool {
	switch fl.Field().String() {
	case "daily", "weekly", "monthly", "quarterly", "half-yearly", "annually":
		return true
	default:
		return false
	}
}

func init() {
	validate := validator.New()
	validate.RegisterValidation("reviewfrequencytype", typeValidation)
}
