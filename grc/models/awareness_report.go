package models

import (
	"time"
)

type AwarenessReport struct {
	ID                       uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ReferenceNumber          string     `json:"reference_number" gorm:"size:50;not null;unique"`
	AwarenessReportTypeID    uint       `json:"awareness_report_type_id" gorm:"not null"`
	SubCategoryID            *uint      `json:"sub_category_id" gorm:"default:null"`
	ReportSourceChannel      *string    `json:"report_source_channel" gorm:"type:varchar(255)"`
	SeverityID               *uint      `json:"severity_id" gorm:"default:null"`
	ReportDate               *time.Time `json:"report_date" gorm:"type:datetime"`
	ResultedFromCampaign     bool       `json:"resulted_from_campaign" gorm:"not null;default:0"`
	AwarenessHandlingStatusID *uint     `json:"awareness_handling_status_id" gorm:"default:null"`
	HandlingOwnerID          *uint      `json:"handling_owner_id" gorm:"default:null"`
	ClosureDate              *time.Time `json:"closure_date" gorm:"type:datetime"`
	Result                   *string    `json:"result" gorm:"type:nvarchar(max)"`

	CreatedBYID uint       `json:"created_by_id" gorm:"not null"`
	CreatedAt   time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedBYID *uint      `json:"updated_by_id" gorm:"default:null"`
	UpdatedAt   *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	AwarenessReportType      *AwarenessCampaignType      `gorm:"foreignKey:AwarenessReportTypeID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION;"`
	SubCategory              *AwarenessReportSubCategory `gorm:"foreignKey:SubCategoryID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION;"`
	SeverityLevel            *SeverityLevel              `gorm:"foreignKey:SeverityID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION;"`
	AwarenessHandlingStatus  *AwarenessHandlingStatus    `gorm:"foreignKey:AwarenessHandlingStatusID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION;"`
	HandlingOwner            *User                       `gorm:"foreignKey:HandlingOwnerID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION;"`

	CreatedBYUser *User `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION;"`
	UpdatedBYUser *User `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION;"`
}

func (AwarenessReport) TableName() string {
	return "awareness_reports"
}