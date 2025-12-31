package models

import (
	"time"
)

type TPRMAssessmentSection struct {
	ID                           uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ThirdPartyAssessmentID       uint      `json:"third_party_assessment_id" binding:"required" gorm:"not null"`
	ThirdPartyQuestionCategoryID uint      `json:"third_party_question_category_id" binding:"required" gorm:"not null"`
	TitleEn                      string    `json:"title_en" binding:"required,max=1000" gorm:"type:nvarchar(1000);not null"`
	TitleAr                      string    `json:"title_ar" binding:"required,max=1000" gorm:"type:nvarchar(1000);not null"`
	Order                        int       `json:"order"`
	CreatedBYID                  uint      `json:"created_by_id" gorm:"not null"`
	CreatedAt                    time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	TPRMAssessment             *ThirdPartyAssessment           `gorm:"foreignKey:ThirdPartyAssessmentID;constraint:OnUpdate:NO ACTION,OnDelete:CASCADE"`
	ThirdPartyQuestionCategory *ThirdPartyQuestionCategory     `gorm:"foreignKey:ThirdPartyQuestionCategoryID;constraint:OnUpdate:NO ACTION,OnDelete:CASCADE"`
	Questions                  []TPRMAssessmentSectionQuestion `gorm:"foreignKey:TPRMAssessmentSectionID;constraint:OnUpdate:NO ACTION,OnDelete:CASCADE"`
	CreatedByUser              *User                           `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

// TableName overrides the table name used by TPRMAssessmentQuestion to `tprm_assessment_sections`
func (TPRMAssessmentSection) TableName() string {
	return "tprm_assessment_sections"
}
