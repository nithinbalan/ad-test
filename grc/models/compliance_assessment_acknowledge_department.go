package models

import (
	"time"
)

type ComplianceAssessmentAcknowledgeDepartment struct {
	ID                     uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ComplianceAssessmentID uint      `gorm:"not null" json:"compliance_assessment_id"`
	UserID                 uint      `gorm:"not null" json:"user_id"`
	DepartmentID           uint      `gorm:"not null" json:"department_id"`
	Comment                *string   `json:"comment"`
	CreatedBYID            uint      `gorm:"not null" json:"created_by_id"`
	CreatedAt              time.Time `gorm:"type:datetime;not null" json:"created_at"`

	// Relationships
	ComplianceAssessment *ComplianceAssessment `gorm:"foreignKey:ComplianceAssessmentID;constraint:OnDelete:NO ACTION,OnUpdate:NO ACTION"`
	User                 *User                 `gorm:"foreignKey:UserID;constraint:OnDelete:NO ACTION,OnUpdate:NO ACTION"`
	Department           *OrganizationUnit     `gorm:"foreignKey:DepartmentID;constraint:OnDelete:NO ACTION,OnUpdate:NO ACTION"`
	CreatedByUser        *User                 `gorm:"foreignKey:CreatedBYID;"`
}

func (ComplianceAssessmentAcknowledgeDepartment) TableName() string {
	return "compliance_assessment_acknowledge_departments"
}
