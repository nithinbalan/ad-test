package models

import "time"

type ComplianceRequirementInherentRiskAnalysis struct {
	ID                               uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	ComplianceRequirementID          uint       `json:"compliance_requirement_id" binding:"required" gorm:"not null"`
	ComplianceAssessmentImpactTypeID uint       `json:"compliance_assessment_impact_type_id" binding:"required" gorm:"not null"`
	Score                            int        `json:"score" binding:"required" gorm:"not null"`
	Comment                          *string    `json:"comment" gorm:"type:nvarchar(max);default:null"`
	CreatedByID                      uint       `json:"created_by_id" gorm:"not null"`
	CreatedAt                        time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedByID                      *uint      `json:"updated_by_id" gorm:"default:null"`
	UpdatedAt                        *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	ComplianceRequirement          *ComplianceRequirement          `gorm:"foreignKey:ComplianceRequirementID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	ComplianceAssessmentImpactType *ComplianceAssessmentImpactType `gorm:"foreignKey:ComplianceAssessmentImpactTypeID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser                  *User                           `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedByUser                  *User                           `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

// TableName overrides the table name
func (ComplianceRequirementInherentRiskAnalysis) TableName() string {
	return "compliance_requirement_inherent_risk_analysis"
}
