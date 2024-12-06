package models

import (
	"time"
)

type BusinessFunctionProcess struct {
	ID                 uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	BusinessFunctionID uint      `json:"business_function_id" binding:"required"`
	ProcessID          uint      `json:"process_id" binding:"required"`
	CreatedByID        uint      `json:"created_by_id" binding:"required"`
	CreatedAt          time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	BusinessFunction *BusinessFunction `gorm:"foreignKey:BusinessFunctionID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Process          *Process          `gorm:"foreignKey:ProcessID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser    *User             `gorm:"foreignKey:CreatedByID"`
}

func (BusinessFunctionProcess) TableName() string {
	return "business_function_processes"
}
