package models

import "time"

type PolicyReportAcknowledgementUser struct {
	ID                      uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	PolicyReportID          uint      `json:"policy_report_id" gorm:"not null;index"`
	UserNameEn              string    `json:"user_name_en" gorm:"type:nvarchar(255)"`
	UserNameAr              string    `json:"user_name_ar" gorm:"type:nvarchar(255)"`
	OrganizationUnitTitleEn string    `json:"organization_unit_title_en" gorm:"type:nvarchar(255)"`
	OrganizationUnitTitleAr string    `json:"organization_unit_title_ar" gorm:"type:nvarchar(255)"`
	Acknowledged            bool      `json:"acknowledged" gorm:"not null;default:false"`
	AcknowledgedOn          *string   `json:"acknowledged_on" gorm:"type:varchar(50)"`
	Comment                 *string   `json:"comment" gorm:"type:nvarchar(max)"`
	CreatedByID             uint      `json:"created_by_id" gorm:"not null;index"`
	CreatedAt               time.Time `json:"created_at" gorm:"not null"`

	// Relations
	PolicyReport  *PolicyReport `gorm:"foreignKey:PolicyReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser *User         `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (PolicyReportAcknowledgementUser) TableName() string {
	return "policy_report_acknowledgement_users"
}
