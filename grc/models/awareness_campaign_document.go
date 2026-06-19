package models

import (
	"time"
)

type AwarenessCampaignDocument struct {
	ID                       uint       `json:"id" gorm:"primary_key;auto_increment"`
	AwarenessCampaignTypeID  uint       `json:"awareness_campaign_type_id" gorm:"not null"`
	AwarenessCampaignDetailID *uint     `json:"awareness_campaign_detail_id" gorm:"default:null"`
	
	URL          string `json:"url" gorm:"type:varchar(1000)"`
	Name         string `json:"name" binding:"required" gorm:"type:nvarchar(255)"`
	Size         int    `json:"size"`
	Token        string `json:"token" gorm:"type:varchar(255)"`
	OriginalName string `json:"original_name" gorm:"type:nvarchar(255)"`
	Type         string `json:"type" gorm:"type:varchar(255)"`

	CreatedByID uint       `json:"created_by_id" gorm:"not null"`
	UpdatedByID *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt   time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt   *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	AwarenessCampaignType   *AwarenessCampaignType   `gorm:"foreignKey:AwarenessCampaignTypeID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION;"`
	AwarenessCampaignDetail *AwarenessCampaignDetail `gorm:"foreignKey:AwarenessCampaignDetailID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION;"`

	CreatedByUser *User `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION;"`
	UpdatedByUser *User `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION;"`
}

func (AwarenessCampaignDocument) TableName() string {
	return "awareness_campaign_documents"
}