package models

import "time"

type GRCReportGovernanceTask struct {
	ID                     uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	GRCReportID            uint       `json:"grc_report_id" gorm:"not null"`
	TaskTitle              string     `json:"task_title" gorm:"type:nvarchar(255);not null"`
	TaskDescription        *string    `json:"task_description" gorm:"type:nvarchar(max)"`
	ResponsibleUserNameEn  string     `json:"responsible_user_name_en" gorm:"type:nvarchar(255);not null"`
	ResponsibleUserNameAr  string     `json:"responsible_user_name_ar" gorm:"type:nvarchar(255);not null"`
	TaskStatus             string     `json:"task_status" gorm:"type:nvarchar(100);not null"`
	TaskStatusTitleEn      string     `json:"task_status_title_en" gorm:"type:nvarchar(255);not null"`
	TaskStatusTitleAr      string     `json:"task_status_title_ar" gorm:"type:nvarchar(255);not null"`
	TaskStatusColorCode    string     `json:"task_status_color_code" gorm:"type:nvarchar(50);not null"`
	TargetDate             *time.Time `json:"target_date" gorm:"type:date;not null"`
	SeverityLevelTitleEn   *string    `json:"severity_level_title_en" gorm:"type:varchar(255);null"`
	SeverityLevelTitleAr   *string    `json:"severity_level_title_ar" gorm:"type:varchar(255);null"`
	SeverityLevelColorCode *string    `json:"severity_level_color_code" gorm:"type:varchar(45);null"`
	CreatedBYID            uint       `json:"created_by_id" gorm:"not null"`
	CreatedAt              time.Time  `json:"created_at" gorm:"type:datetime;not null"`

	GRCReport     *GRCReport `gorm:"foreignKey:GRCReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser *User      `gorm:"foreignKey:CreatedBYID"`
}

func (GRCReportGovernanceTask) TableName() string {
	return "grc_report_governance_tasks"
}
