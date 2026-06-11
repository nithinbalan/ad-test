package models

import "time"

type ComplianceAssessmentCompliance struct {
	ID                     uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ComplianceRegisterID   uint      `json:"compliance_register_id" binding:"required" gorm:"not null"`
	ComplianceAssessmentID uint      `json:"compliance_assessment_id" binding:"required" gorm:"not null"`
	CreatedAt              time.Time `json:"created_at" gorm:"type:datetime;not null"`
	CreatedByID            uint      `json:"created_by_id" gorm:"not null"`

	ComplianceRegister   *ComplianceRegister   `gorm:"foreignKey:ComplianceRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ComplianceAssessment *ComplianceAssessment `gorm:"foreignKey:ComplianceAssessmentID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser        *User                 `gorm:"foreignKey:CreatedByID"`
}

func (ComplianceAssessmentCompliance) TableName() string {
	return "compliance_assessment_compliances"
}
