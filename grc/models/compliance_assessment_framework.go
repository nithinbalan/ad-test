package models

import "time"

type ComplianceAssessmentFramework struct {
	ID                        uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	TitleEn                   string     `json:"title_en" gorm:"type:nvarchar(255);not null"`
	TitleAr                   string     `json:"title_ar" gorm:"type:nvarchar(255);not null"`
	Description               *string    `json:"description" gorm:"type:nvarchar(max)"`
	AssessmentFrameworkTypeID uint       `json:"compliance_assessment_framework_type_id" gorm:"column:compliance_assessment_framework_type_id;not null"`
	Status                    string     `json:"status" gorm:"type:varchar(10);default:active"`
	CreatedBYID               uint       `json:"created_by_id" gorm:"column:created_by_id;not null"`
	UpdatedBYID               *uint      `json:"updated_by_id" gorm:"column:updated_by_id"`
	CreatedAt                 time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                 *time.Time `json:"updated_at" gorm:"type:datetime"`

	AssessmentFrameworkType *ComplianceAssessmentFrameworkType    `gorm:"foreignKey:AssessmentFrameworkTypeID;constraint:OnDelete:CASCADE"`
	CreatedByUser           *User                                 `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser           *User                                 `gorm:"foreignKey:UpdatedBYID"`
	Values                  []*ComplianceAssessmentFrameworkValue `gorm:"foreignKey:ComplianceAssessmentFrameworkID;constraint:OnDelete:CASCADE"`
}

func (ComplianceAssessmentFramework) TableName() string {
	return "compliance_assessment_frameworks"
}
