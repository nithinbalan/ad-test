package models

import (
	"time"
)

type RiskWorkflowHistory struct {
	ID               uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	RiskRegisterID   uint       `json:"risk_register_id" binding:"required" gorm:"not null;index"`
	Level            int        `json:"level" binding:"required" gorm:"not null"`
	WorkflowStatusID uint       `json:"workflow_status_id" binding:"required" gorm:"not null"`
	Comment          string     `json:"comment" gorm:"type:text"`
	CreatedBYID      uint       `json:"created_by_id" binding:"required" gorm:"not null"`
	UpdatedBYID      *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt        time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	UpdatedAt        *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	RiskRegister   *RiskRegister   `gorm:"foreignKey:RiskRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	WorkflowStatus *WorkflowStatus `gorm:"foreignKey:WorkflowStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser  *User           `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	UpdatedByUser  *User           `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (RiskWorkflowHistory) TableName() string {
	return "risk_workflow_histories"
}
