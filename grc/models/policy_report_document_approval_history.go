package models

import "time"

type PolicyReportDocumentApprovalHistory struct {
	ID                 uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	PolicyReportID     uint      `json:"policy_report_id" gorm:"not null;index"`
	ApprovedUserNameEn string    `json:"approved_user_name_en" gorm:"type:nvarchar(255)"`
	ApprovedUserNameAr string    `json:"approved_user_name_ar" gorm:"type:nvarchar(255)"`
	DesignationTitleEn string    `json:"designation_title_en" gorm:"type:nvarchar(255)"`
	DesignationTitleAr string    `json:"designation_title_ar" gorm:"type:nvarchar(255)"`
	ActionTitleEn      string    `json:"action_title_en" gorm:"type:nvarchar(100)"`
	ActionTitleAr      string    `json:"action_title_ar" gorm:"type:nvarchar(100)"`
	ActionColorCode    string    `json:"action_color_code" gorm:"type:varchar(9)"`
	ActionDate         *string   `json:"action_date" gorm:"type:varchar(50)"`
	Comment            *string   `json:"comment" gorm:"type:nvarchar(max)"`
	CreatedByID        uint      `json:"created_by_id" gorm:"not null;index"`
	CreatedAt          time.Time `json:"created_at" gorm:"not null"`

	// Relations
	PolicyReport  *PolicyReport `gorm:"foreignKey:PolicyReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser *User         `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (PolicyReportDocumentApprovalHistory) TableName() string {
	return "policy_report_document_approval_histories"
}
