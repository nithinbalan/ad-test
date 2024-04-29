package models

import (
	"time"
)

type RiskAppetiteStatement struct {
	ID              uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ReferenceNumber string     `json:"reference_number" binding:"omitempty,max=45" gorm:"size:45"`
	RiskCategoryID  uint       `json:"risk_category_id" binding:"required" gorm:"not null"`
	StatementEn     string     `json:"statement_en" gorm:"type:varchar(max);not null"`
	StatementAr     string     `json:"statement_ar" gorm:"type:varchar(max);not null"`
	Status          string     `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null"`
	CreatedAt       time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt       *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedBy       uint       `json:"created_by" gorm:"not null"`
	UpdatedBy       *uint      `json:"updated_by" gorm:"default:null"`

	// Relationships
	RiskCategory  *RiskCategory `gorm:"foreignkey:RiskCategoryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser *User         `gorm:"foreignKey:CreatedBy"`
	UpdatedByUser *User         `gorm:"foreignKey:UpdatedBy"`
}

// TableName overrides the table name used by RiskAppetiteStatement to `risk_appetite_statements`
func (RiskAppetiteStatement) TableName() string {
	return "risk_appetite_statements"
}
