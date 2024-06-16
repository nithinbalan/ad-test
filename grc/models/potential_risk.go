package models

import (
	"time"
)

type PotentialRisk struct {
	ID                    uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	RiskCategoryID        uint       `json:"risk_category_id" binding:"required" gorm:"not null"`
	SeverityLevelID       uint       `json:"severity_level_id" binding:"required" gorm:"not null"`
	PotentialRiskStatusID uint       `json:"potential_risk_status_id" gorm:"not null"`
	Title                 string     `json:"title" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	Description           *string    `json:"description" gorm:"type:varchar(max)"`
	RiskSource            string     `json:"risk_source" binding:"required,max=45" gorm:"type:varchar(45);not null"`
	CreatedAt             time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt             *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedByID           uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBy             *uint      `json:"updated_by_id" gorm:"default:null"`

	// Relationships
	RiskCategory        *RiskCategory        `gorm:"foreignKey:RiskCategoryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	SeverityLevel       *SeverityLevel       `gorm:"foreignKey:SeverityLevelID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	PotentialRiskStatus *PotentialRiskStatus `gorm:"foreignKey:PotentialRiskStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser       *User                `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser       *User                `gorm:"foreignKey:UpdatedBy"`
}

// TableName overrides the table name used by PotentialRisk to `potential_risks`
func (PotentialRisk) TableName() string {
	return "potential_risks"
}
