package models

import "time"

type IdentityReviewCampaignDocument struct {
	ID                       uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	IdentityReviewCampaignID uint       `json:"identity_review_campaign_id" gorm:"not null"`
	URL                      string     `json:"url" gorm:"type:varchar(1000)"`
	Name                     string     `json:"name" gorm:"type:nvarchar(255)"`
	Size                     int        `json:"size"`
	Token                    string     `json:"token" gorm:"type:varchar(255)"`
	OriginalName             string     `json:"original_name" gorm:"type:nvarchar(255)"`
	Type                     string     `json:"type" gorm:"type:varchar(255)"`
	CreatedByID              uint       `json:"created_by_id"`
	UpdatedByID              *uint      `json:"updated_by_id"`
	CreatedAt                time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                *time.Time `json:"updated_at" gorm:"type:datetime"`

	IdentityReviewCampaign *IdentityReviewCampaign `json:"-" gorm:"foreignKey:IdentityReviewCampaignID"`
	CreatedByUser          *User                   `json:"created_by" gorm:"foreignKey:CreatedByID"`
	UpdatedByUser          *User                   `json:"updated_by" gorm:"foreignKey:UpdatedByID"`
}

func (IdentityReviewCampaignDocument) TableName() string {
	return "identity_review_campaign_documents"
}
