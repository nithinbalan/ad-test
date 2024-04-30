package models

import (
	"time"

	"jethur/grc/models"
)

type ComplianceRequirementStatusUpdate struct {
	ID                      uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ComplianceRequirementID uint       `json:"compliance_requirement_id" gorm:"not null"`
	ComplianceStatusID      uint       `json:"compliance_status_id" gorm:"not null"`
	Comment                 string     `json:"comment" binding:"required" gorm:"type:nvarchar(max);not null"`
	CreatedAt               time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt               *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedBy               uint       `json:"created_by" gorm:"not null"`
	UpdatedBy               *uint      `json:"updated_by" gorm:"default:null"`

	//Relationships
	ComplianceRequirement *ComplianceRequirement `gorm:"foreignkey:ComplianceRequirementID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ComplianceStatus      *ComplianceStatus      `gorm:"foreignkey:ComplianceStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser         *models.User           `gorm:"foreignKey:CreatedBy"`
	UpdatedByUser         *models.User           `gorm:"foreignKey:UpdatedBy"`
}

// TableName overrides the table name used by ComplianceRequirementStatusUpdate to `compliance_requirement_status_updates`
func (ComplianceRequirementStatusUpdate) TableName() string {
	return "compliance_requirement_status_updates"
}
