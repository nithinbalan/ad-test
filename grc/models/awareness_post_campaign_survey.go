package models

import (
	"time"
)

type AwarenessPostCampaignSurvey struct {
	ID              uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ReferenceNumber string `json:"reference_number" gorm:"type:varchar(50);unique;not null"`

	LinkedCampaignID *uint      `json:"linked_campaign_id"`
	SendDate         *time.Time `json:"send_date" gorm:"type:datetime"`

	TargetAudience *string `json:"target_audience" gorm:"type:nvarchar(max)"`
	Recipients     *int    `json:"recipients"`
	ResponseRate   *int    `json:"response_rate"`
	QuestionCount  *int    `json:"question_count"`
	AverageRating  *int    `json:"average_rating"`

	MaturityLevelID *uint   `json:"maturity_level_id"`
	Recommendations *string `json:"recommendations" gorm:"type:nvarchar(max)"`
	ExecutingParty  *string `json:"executing_party" gorm:"type:nvarchar(max)"`
	ResultsLink     *string `json:"results_link" gorm:"type:nvarchar(max)"`

	CreatedBYID uint       `json:"created_by_id"`
	CreatedAt   time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedBYID *uint      `json:"updated_by_id"`
	UpdatedAt   *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	LinkedCampaign *AwarenessCampaign      `gorm:"foreignKey:LinkedCampaignID"`
	MaturityLevel  *AwarenessMaturityLevel `gorm:"foreignKey:MaturityLevelID"`

	CreatedByUser *User `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser *User `gorm:"foreignKey:UpdatedBYID"`

	SurveyThemes []*AwarenessPostCampaignSurveyTheme `gorm:"foreignKey:AwarenessPostCampaignSurveyID"`
}

func (AwarenessPostCampaignSurvey) TableName() string {
	return "awareness_post_campaign_surveys"
}
