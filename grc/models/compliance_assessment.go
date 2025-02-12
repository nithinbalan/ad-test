package models

import "time"

type ComplianceAssessment struct {
	ID                           uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ComplianceRegisterID         uint       `json:"compliance_register_id" binding:"required"`
	PrimaryUnitID                uint       `json:"organization_unit_id" binding:"required" gorm:"column:organization_unit_id;not null"`
	ReviewUserID                 uint       `json:"review_user_id" binding:"required"`
	ComplianceAssessmentStatusID uint       `json:"compliance_assessment_status_id" gorm:"not null;default:1"`
	SubmittedBy                  *uint      `json:"submitted_by"`
	SubmittedAt                  *time.Time `json:"submitted_at"`
	ReviewedBy                   *uint      `json:"reviewed_by"`
	Description                  *string    `json:"description"`
	CreatedAt                    time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                    *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedBYID                  uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID                  *uint      `json:"updated_by_id" gorm:"default:null"`

	// Relationships
	ComplianceRegister         *ComplianceRegister                `gorm:"foreignKey:ComplianceRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	OrganizationUnit           *OrganizationUnit                  `gorm:"foreignKey:PrimaryUnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ReviewUser                 *User                              `gorm:"foreignKey:ReviewUserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ComplianceAssessmentStatus *ComplianceAssessmentStatus        `gorm:"foreignKey:ComplianceAssessmentStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	SubmittedByUser            *User                              `gorm:"foreignKey:SubmittedBy"`
	ReviewedByUser             *User                              `gorm:"foreignKey:ReviewedBy"`
	CreatedByUser              *User                              `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser              *User                              `gorm:"foreignKey:UpdatedBYID"`
	AssessmentRequirements     []*ComplianceAssessmentRequirement `gorm:"foreignKey:ComplianceAssessmentID"`
}

// TableName overrides the table name used by ComplianceAssessment to `compliance_assessments`
func (ComplianceAssessment) TableName() string {
	return "compliance_assessments"
}
