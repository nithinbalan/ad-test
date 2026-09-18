package models

import "time"

type ProjectSecurityGateReport struct {
	ID                uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Title             string    `json:"title" binding:"required,max=255" gorm:"type:nvarchar(255);not null"`
	SecurityGateIDs   string    `json:"security_gate_ids" gorm:"type:nvarchar(max);not null"`
	ResponseStatusIDs string    `json:"response_status_ids" gorm:"type:nvarchar(max);not null"`
	ReportData        string    `json:"-" gorm:"column:report_data;type:nvarchar(max);not null"`
	CreatedBYID       uint      `json:"created_by_id" gorm:"not null"`
	CreatedAt         time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	CreatedByUser *User `json:"created_by_user" gorm:"foreignKey:CreatedBYID"`
}

func (ProjectSecurityGateReport) TableName() string {
	return "project_security_gate_reports"
}
