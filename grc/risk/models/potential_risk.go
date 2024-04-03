package models

import (
	"time"

	masterModels "jethur/grc/masters/models"
	organizationModels "jethur/grc/organization/models"
)

type PotentialRisk struct {
	ID                    uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	RiskCategoryID        uint       `json:"risk_category_id" binding:"required" gorm:"not null"`
	SeverityLevelID       uint       `json:"severity_level_id" binding:"required" gorm:"not null"`
	PotentialRiskStatusID uint       `json:"potential_risk_status_id" binding:"required" gorm:"not null"`
	Title                 string     `json:"title" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	Description           *string    `json:"description" gorm:"type:varchar(max)"`
	RiskSource            string     `json:"risk_source" binding:"required,max=45" gorm:"type:varchar(45);not null"`
	CreatedAt             time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt             *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedBy             uint       `json:"created_by" gorm:"not null"`
	UpdatedBy             *uint      `json:"updated_by" gorm:"default:null"`

	// Relationships
	RiskCategory        *RiskCategory               `gorm:"foreignkey:RiskCategoryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	SeverityLevel       *masterModels.SeverityLevel `gorm:"foreignkey:SeverityLevelID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	PotentialRiskStatus *PotentialRiskStatus        `gorm:"foreignkey:PotentialRiskStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser       *organizationModels.User    `gorm:"foreignKey:CreatedBy"`
	UpdatedByUser       *organizationModels.User    `gorm:"foreignKey:UpdatedBy"`
}

// TableName overrides the table name used by PotentialRisk to `potential_risks`
func (PotentialRisk) TableName() string {
	return "potential_risks"
}
