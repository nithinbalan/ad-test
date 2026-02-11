package models

import "time"

type ComplianceAssessmentResidualControl struct {
	ID                                               uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	ComplianceAssessmentRequirementResponsibleUnitID uint      `json:"compliance_assessment_requirement_responsible_unit_id" binding:"required" gorm:"not null"`
	ComplianceAssessmentImpactTypeID                 uint      `json:"compliance_assessment_impact_type_id" binding:"required" gorm:"not null"`
	ControlID                                        *uint     `json:"control_id"`
	LinkControl                                      bool      `json:"link_control" gorm:"default:false"`
	IsNotApplicable                                  bool      `json:"is_not_applicable" gorm:"default:false"`
	CreatedBYID                                      uint      `json:"created_by_id" gorm:"not null"`
	CreatedAt                                        time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	ComplianceAssessmentRequirementResponsibleUnit *ComplianceAssessmentRequirementResponsibleUnit `gorm:"foreignKey:ComplianceAssessmentRequirementResponsibleUnitID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	ComplianceAssessmentImpactType                 *ComplianceAssessmentImpactType                 `gorm:"foreignKey:ComplianceAssessmentImpactTypeID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Control                                        *Control                                        `gorm:"foreignKey:ControlID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser                                  *User                                           `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

// TableName overrides the table name
func (ComplianceAssessmentResidualControl) TableName() string {
	return "compliance_assessment_residual_controls"
}
