package models

import "time"

type ComplianceAssessment struct {
	ID                   uint  `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ComplianceRegisterID *uint `json:"compliance_register_id"`
	// PrimaryUnitID                uint       `json:"organization_unit_id" binding:"required" gorm:"column:organization_unit_id;not null"`
	ReviewUserID                 uint       `json:"review_user_id" binding:"required"`
	ComplianceAssessmentStatusID uint       `json:"compliance_assessment_status_id" gorm:"not null;default:1"`
	AverageResidualRiskScore     *float64   `json:"average_residual_risk_score" gorm:"type:decimal(10,2)"`
	AverageResidualRiskRatingID  *uint      `json:"average_residual_risk_rating_id" gorm:"default:null"`
	SubmittedBy                  *uint      `json:"submitted_by"`
	SubmittedAt                  *time.Time `json:"submitted_at"`
	ReviewedBy                   *uint      `json:"reviewed_by"`
	Description                  *string    `json:"description"`
	Type                         *string    `json:"type" gorm:"default:framework"`
	ResponsibleUnitID            *uint      `json:"responsible_unit_id"`
	CreatedAt                    time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                    *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedBYID                  uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID                  *uint      `json:"updated_by_id" gorm:"default:null"`

	// Relationships
	ComplianceRegister *ComplianceRegister `gorm:"foreignKey:ComplianceRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	// OrganizationUnit           *OrganizationUnit                  `gorm:"foreignKey:PrimaryUnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ReviewUser                      *User                              `gorm:"foreignKey:ReviewUserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ComplianceAssessmentStatus      *ComplianceAssessmentStatus        `gorm:"foreignKey:ComplianceAssessmentStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	AverageResidualRiskRating       *ComplianceAssessmentRating        `gorm:"foreignKey:AverageResidualRiskRatingID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	SubmittedByUser                 *User                              `gorm:"foreignKey:SubmittedBy"`
	ReviewedByUser                  *User                              `gorm:"foreignKey:ReviewedBy"`
	CreatedByUser                   *User                              `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser                   *User                              `gorm:"foreignKey:UpdatedBYID"`
	AssessmentRequirements          []*ComplianceAssessmentRequirement `gorm:"foreignKey:ComplianceAssessmentID"`
	AssessmentControls              []*ComplianceAssessmentControl     `gorm:"foreignKey:ComplianceAssessmentID"` // ✅ New relationship
	ComplianceAssessmentCompliances []*ComplianceAssessmentCompliance  `gorm:"foreignKey:ComplianceAssessmentID"`
	ResponsibleOrganizationUnit     *OrganizationUnit                  `gorm:"foreignKey:ResponsibleUnitID"`
}

// TableName overrides the table name used by ComplianceAssessment to `compliance_assessments`
func (ComplianceAssessment) TableName() string {
	return "compliance_assessments"
}
