package models

import (
	"time"
)

type AwarenessCampaignChannel struct {
	ID uint `json:"id" gorm:"primary_key;AUTO_INCREMENT"`

	AwarenessCampaignDetailsID uint `json:"awareness_campaign_details_id" gorm:"not null"`
	AwarenessChannelID         uint `json:"awareness_channel_id" gorm:"not null"`

	CreatedByID uint      `json:"created_by_id" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	AwarenessCampaignDetail *AwarenessCampaignDetail `gorm:"foreignKey:AwarenessCampaignDetailsID;constraint:OnUpdate:NO ACTION,OnDelete:CASCADE;"`
	AwarenessChannel        *AwarenessChannel        `gorm:"foreignKey:AwarenessChannelID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION;"`
	CreatedByUser           *User                    `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION;"`
}

func (AwarenessCampaignChannel) TableName() string {
	return "awareness_campaign_channels"
}
