package models

import "time"

type ComplianceAssessment struct {
	ID                   uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ComplianceRegisterID uint       `json:"compliance_register_id" binding:"required" gorm:"column:compliance_register_id;not null"`
	OrganizationUnitID   uint       `json:"organization_unit_id" binding:"required" gorm:"not null"`
	ReviewUserID         uint       `json:"review_user_id" binding:"required" gorm:"column:review_user_id;not null"`
	Description          *string    `json:"description"`
	CreatedAt            time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt            *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedBYID          uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID          *uint      `json:"updated_by_id" gorm:"default:null"`

	// Relationships

	ComplianceRegister     *ComplianceRegister                `gorm:"foreignKey:ComplianceRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	OrganizationUnit       *OrganizationUnit                  `gorm:"foreignKey:OrganizationUnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ReviewUser             *User                              `gorm:"foreignKey:ReviewUserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser          *User                              `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser          *User                              `gorm:"foreignKey:UpdatedBYID"`
	AssessmentRequirements []*ComplianceAssessmentRequirement `gorm:"foreignKey:ComplianceAssessmentID"`
}

// TableName overrides the table name used by Compliance Assessment to `compliance_assessments`
func (ComplianceAssessment) TableName() string {
	return "compliance_assessments"
}
