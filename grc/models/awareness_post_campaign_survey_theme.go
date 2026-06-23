package models

import (
	"time"
)

type AwarenessPostCampaignSurveyTheme struct {
	ID uint `json:"id" gorm:"primary_key;AUTO_INCREMENT"`

	AwarenessPostCampaignSurveyID uint `json:"awareness_post_campaign_survey_id"`
	AwarenessAssessmentThemeID    uint `json:"awareness_assessment_theme_id"`

	CreatedByID uint      `json:"created_by_id"`
	CreatedAt   time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	AwarenessPostCampaignSurvey *AwarenessPostCampaignSurvey `gorm:"foreignKey:AwarenessPostCampaignSurveyID"`
	AwarenessAssessmentTheme    *AwarenessAssessmentTheme    `gorm:"foreignKey:AwarenessAssessmentThemeID"`

	CreatedByUser *User `gorm:"foreignKey:CreatedByID"`
}

func (AwarenessPostCampaignSurveyTheme) TableName() string {
	return "awareness_post_campaign_survey_themes"
}
