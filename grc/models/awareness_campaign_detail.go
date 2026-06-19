package models

import (
	"time"
)

type AwarenessCampaignDetail struct {
	ID                  uint    `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	AwarenessCampaignID uint    `json:"awareness_campaign_id" binding:"required"`
	AwarenessPackName   *string `json:"awareness_pack_name" gorm:"type:nvarchar(255)"`
	RecipientCount      *int    `json:"recipient_count"`
	SentCount           *int    `json:"sent_count"`
	EmailOpenRate       *int    `json:"email_open_rate"`
	ClickRate           *int    `json:"click_rate"`
	DataSubmissionRate  *int    `json:"data_submission_rate"`
	ReportRate          *int    `json:"report_rate"`

	RiskClassificationID       *uint   `json:"risk_classification_id"`
	CampaignResult             *string `json:"campaign_result" gorm:"type:nvarchar(max)"`
	ImprovementRecommendations *string `json:"improvement_recommendations" gorm:"type:nvarchar(max)"`

	SubtypeRole *string `json:"subtype_role" gorm:"type:varchar(100)"`
	Presenter   *string `json:"presenter" gorm:"type:nvarchar(255)"`
	Agenda      *string `json:"agenda" gorm:"type:nvarchar(max)"`

	SessionStartTime *time.Time `json:"session_start_time" gorm:"type:time"`
	SessionEndTime   *time.Time `json:"session_end_time" gorm:"type:time"`

	Venue         *string `json:"venue" gorm:"type:nvarchar(255)"`
	AttendeeCount *int    `json:"attendee_count"`

	AwarenessEvaluationRatingID *uint `json:"awareness_evaluation_rating_id"`

	CreatedByID uint       `json:"created_by_id"`
	CreatedAt   time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedByID *uint      `json:"updated_by_id"`
	UpdatedAt   *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	AwarenessCampaign         *AwarenessCampaign         `gorm:"foreignKey:AwarenessCampaignID"`
	RiskClassification        *AwarenessPriority         `gorm:"foreignKey:RiskClassificationID"`
	AwarenessEvaluationRating *AwarenessEvaluationRating `gorm:"foreignKey:AwarenessEvaluationRatingID"`

	CreatedByUser *User `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser *User `gorm:"foreignKey:UpdatedByID"`

	Channels  []*AwarenessCampaignChannel  `gorm:"foreignKey:AwarenessCampaignDetailsID"`
	Documents []*AwarenessCampaignDocument `gorm:"foreignKey:AwarenessCampaignDetailID"`
}

func (AwarenessCampaignDetail) TableName() string {
	return "awareness_campaign_details"
}
