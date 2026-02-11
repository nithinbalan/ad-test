package models

import "time"

type ComplianceAssessmentRequirementResidualRiskAnalysis struct {
	ID                                               uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	ComplianceAssessmentRequirementResponsibleUnitID uint       `json:"compliance_assessment_requirement_responsible_unit_id" binding:"required" gorm:"not null"`
	ComplianceAssessmentImpactTypeID                 uint       `json:"compliance_assessment_impact_type_id" binding:"required" gorm:"not null"`
	Score                                            int        `json:"score" binding:"required" gorm:"not null"`
	Comment                                          *string    `json:"comment" gorm:"type:nvarchar(max);default:null"`
	CreatedBYID                                      uint       `json:"created_by_id" gorm:"not null"`
	CreatedAt                                        time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedBYID                                      *uint      `json:"updated_by_id" gorm:"default:null"`
	UpdatedAt                                        *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	ComplianceAssessmentRequirementResponsibleUnit *ComplianceAssessmentRequirementResponsibleUnit `gorm:"foreignKey:ComplianceAssessmentRequirementResponsibleUnitID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	ComplianceAssessmentImpactType                 *ComplianceAssessmentImpactType                 `gorm:"foreignKey:ComplianceAssessmentImpactTypeID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser                                  *User                                           `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedByUser                                  *User                                           `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

// TableName overrides the table name
func (ComplianceAssessmentRequirementResidualRiskAnalysis) TableName() string {
	return "compliance_assessment_requirement_residual_risk_analysis"
}
