package models

import (
	"time"
)

type RiskResponseWorkflowHistory struct {
	ID               uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	RiskResponseID   uint       `json:"risk_response_id" binding:"required" gorm:"not null;index"`
	Level            int        `json:"level" binding:"required" gorm:"not null"`
	WorkflowStatusID uint       `json:"workflow_status_id" binding:"required" gorm:"not null"`
	Comment          string     `json:"comment" gorm:"type:text"`
	CreatedBYID      uint       `json:"created_by_id" binding:"required" gorm:"not null"`
	UpdatedBYID      *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt        time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	UpdatedAt        *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	RiskResponse   *RiskResponse   `gorm:"foreignKey:RiskResponseID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	WorkflowStatus *WorkflowStatus `gorm:"foreignKey:WorkflowStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser  *User           `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	UpdatedByUser  *User           `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (RiskResponseWorkflowHistory) TableName() string {
	return "risk_response_workflow_histories"
}
