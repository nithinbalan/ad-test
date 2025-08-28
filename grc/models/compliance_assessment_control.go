package models

import "time"

type ComplianceAssessmentControl struct {
	ID                                      uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ComplianceAssessmentID                  uint       `json:"compliance_assessment_id" binding:"required" gorm:"column:compliance_assessment_id;not null"`
	ComplianceControlID                     uint       `json:"compliance_control_id" binding:"required" gorm:"column:compliance_control_id;not null"`
	CompletedByUserID                       *uint      `json:"completed_by_user_id" gorm:"column:completed_by_user_id"`
	CompletedAt                             *time.Time `json:"completed_at" gorm:"column:completed_at"`
	ComplianceAssessmentFrameworkValueID    *uint      `json:"compliance_assessment_framework_value_id" gorm:"column:compliance_assessment_framework_value_id"`
	Comments                                *string    `json:"comments" gorm:"type:nvarchar(max)"`
	IsCompleted                             bool       `json:"is_completed" gorm:"column:is_completed;not null;default:false"`
	Score                                   *float64   `json:"score" gorm:"type:decimal(10,2)"`
	ComplianceAssessmentRequirementStatusID *uint      `json:"compliance_assessment_requirement_status_id" gorm:"column:compliance_assessment_requirement_status_id;default:null"`
	CreatedAt                               time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                               *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedBYID                             uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID                             *uint      `json:"updated_by_id" gorm:"default:null"`

	// Relationships
	ComplianceAssessment                  *ComplianceAssessment                  `gorm:"foreignKey:ComplianceAssessmentID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ComplianceControl                     *ComplianceControl                     `gorm:"foreignKey:ComplianceControlID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CompletedByUser                       *User                                  `gorm:"foreignKey:CompletedByUserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ComplianceAssessmentFrameworkValue    *ComplianceAssessmentFrameworkValue    `gorm:"foreignKey:ComplianceAssessmentFrameworkValueID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ComplianceAssessmentRequirementStatus *ComplianceAssessmentRequirementStatus `gorm:"foreignKey:ComplianceAssessmentRequirementStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser                         *User                                  `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser                         *User                                  `gorm:"foreignKey:UpdatedBYID"`
}

// TableName overrides the table name used by ComplianceAssessmentControl to `compliance_assessment_controls`
func (ComplianceAssessmentControl) TableName() string {
	return "compliance_assessment_controls"
}
