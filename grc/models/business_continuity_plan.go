package models

import (
	"time"
)

type BusinessContinuityPlan struct {
	ID                 uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ReferenceNumber    string     `json:"reference_number" binding:"omitempty,max=45" gorm:"size:45"`
	BusinessFunctionID uint       `json:"business_function_id" binding:"required"`
	Version            *int       `json:"version"`
	Purpose            *string    `json:"purpose"`
	Scope              *string    `json:"scope"`
	Objectives         *string    `json:"objectives"`
	Assumptions        *string    `json:"assumptions"`
	CreatedBYID        uint       `json:"created_by_id" binding:"required"`
	UpdatedBYID        *uint      `json:"updated_by_id"`
	CreatedAt          time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt          *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	BusinessFunction *BusinessFunction `gorm:"foreignKey:BusinessFunctionID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser    *User             `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser    *User             `gorm:"foreignKey:UpdatedBYID"`
}

// TableName overrides the table name used by BusinessContinuityPlan to `business_continuity_plans`
func (BusinessContinuityPlan) TableName() string {
	return "business_continuity_plans"
}
