package models

import "time"

type ComplianceAssessmentFrameworkValue struct {
	ID                              uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	ComplianceAssessmentFrameworkID uint       `json:"compliance_assessment_framework_id" binding:"required"`
	ComplianceStatusID              uint       `json:"compliance_status_id" binding:"required"`
	Value                           *int       `json:"value"`
	TitleEn                         *string    `json:"title_en" gorm:"type:nvarchar(255)"`
	TitleAr                         *string    `json:"title_ar" gorm:"type:nvarchar(255)"`
	Description                     *string    `json:"description" gorm:"type:nvarchar(max)"`
	ScoreMin                        *int       `json:"score_min"`
	ScoreMax                        *int       `json:"score_max"`
	IsApplicable                    bool       `json:"is_applicable" gorm:"not null;default:1"`
	CreatedBYID                     uint       `json:"created_by_id" binding:"required"`
	UpdatedBYID                     *uint      `json:"updated_by_id"`
	CreatedAt                       time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                       *time.Time `json:"updated_at"`

	ComplianceAssessmentFramework *ComplianceAssessmentFramework `gorm:"foreignKey:ComplianceAssessmentFrameworkID"`
	ComplianceStatus              *ComplianceStatus              `gorm:"foreignKey:ComplianceStatusID"`
	CreatedByUser                 *User                          `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser                 *User                          `gorm:"foreignKey:UpdatedBYID"`
}

func (ComplianceAssessmentFrameworkValue) TableName() string {
	return "compliance_assessment_framework_values"
}
