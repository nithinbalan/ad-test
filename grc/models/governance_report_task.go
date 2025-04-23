// models/governance_report_task.go
package models

import "time"

type GovernanceReportTask struct {
	ID                    uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	GovernanceReportID    uint       `json:"governance_report_id" gorm:"not null"`
	TaskTitle             string     `json:"task_title" gorm:"type:nvarchar(255);not null"`
	TaskDescription       *string    `json:"task_description" gorm:"type:nvarchar(max);null"`
	ResponsibleUserNameEn string     `json:"responsible_user_name_en" gorm:"type:nvarchar(255);not null"`
	ResponsibleUserNameAr string     `json:"responsible_user_name_ar" gorm:"type:nvarchar(255);not null"`
	TaskStatus            string     `json:"task_status" gorm:"type:nvarchar(100);not null"`
	TaskStatusTitleEn     string     `json:"task_status_title_en" gorm:"type:nvarchar(255);not null"`
	TaskStatusTitleAr     string     `json:"task_status_title_ar" gorm:"type:nvarchar(255);not null"`
	TaskStatusColorCode   string     `json:"task_status_color_code" gorm:"type:nvarchar(50);not null"`
	TargetDate            *time.Time `json:"target_date" gorm:"type:date;null"`
	CreatedByID           uint       `json:"created_by_id" gorm:"not null"`
	CreatedAt             time.Time  `json:"created_at" gorm:"type:datetime;not null"`

	CreatedBy *User `gorm:"foreignKey:CreatedByID"`
}

func (GovernanceReportTask) TableName() string {
	return "governance_report_tasks"
}
