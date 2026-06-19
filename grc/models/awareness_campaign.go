package models

import (
	"time"
)

type AwarenessCampaign struct {
	ID                        uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ReferenceNumber           string     `json:"reference_number" gorm:"size:50;not null;unique"`
	AwarenessCampaignTypeID   uint       `json:"awareness_campaign_type_id" binding:"required"`
	CampaignTitle             string     `json:"campaign_title" binding:"required,max=255" gorm:"type:nvarchar(255);not null"`
	Objective                 string     `json:"objective" binding:"required" gorm:"type:nvarchar(max);not null"`
	AwarenessPriorityID       *uint      `json:"awareness_priority_id"`
	TargetAudience            string     `json:"target_audience" binding:"required" gorm:"type:nvarchar(max);not null"`
	OwnerID                   uint       `json:"owner_id" binding:"required"`
	OrganizationUnitID        uint       `json:"organization_unit_id" binding:"required"`
	StartDate                 time.Time  `json:"start_date" binding:"required" gorm:"type:datetime;not null"`
	EndDate                   *time.Time `json:"end_date" gorm:"type:datetime"`
	AwarenessCampaignStatusID uint       `json:"awareness_campaign_status_id" binding:"required"`
	LinkedSecurityControl     *string    `json:"linked_security_control" gorm:"type:nvarchar(255)"`

	CreatedByID uint       `json:"created_by_id"`
	CreatedAt   time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedByID *uint      `json:"updated_by_id"`
	UpdatedAt   *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	AwarenessCampaignType   *AwarenessCampaignType   `gorm:"foreignKey:AwarenessCampaignTypeID"`
	AwarenessPriority       *AwarenessPriority       `gorm:"foreignKey:AwarenessPriorityID"`
	AwarenessCampaignStatus *AwarenessCampaignStatus `gorm:"foreignKey:AwarenessCampaignStatusID"`
	Owner                   *User                    `gorm:"foreignKey:OwnerID"`
	OrganizationUnit        *OrganizationUnit        `gorm:"foreignKey:OrganizationUnitID"`

	CreatedByUser *User `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser *User `gorm:"foreignKey:UpdatedByID"`

	AwarenessCampaignDetails *AwarenessCampaignDetail `gorm:"foreignKey:AwarenessCampaignID"`
}

func (AwarenessCampaign) TableName() string {
	return "awareness_campaigns"
}
