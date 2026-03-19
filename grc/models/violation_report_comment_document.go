package models

import "time"

type ViolationReportCommentDocument struct {
	ID                       uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	ViolationReportCommentID uint       `json:"violation_report_comment_id" gorm:"not null"`
	URL                      string     `json:"url"`
	Name                     string     `json:"name"`
	Size                     int        `json:"size"`
	Token                    string     `json:"token"`
	OriginalName             string     `json:"original_name"`
	Type                     string     `json:"type"`
	CreatedByID              uint       `json:"created_by_id"`
	UpdatedByID              *uint      `json:"updated_by_id"`
	CreatedAt                time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	ViolationReportComment *ViolationReportComment `gorm:"foreignKey:ViolationReportCommentID;constraint:OnUpdate:RESTRICT,OnDelete:CASCADE;"`
	CreatedByUser          *User                   `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser          *User                   `gorm:"foreignKey:UpdatedByID"`
}

func (ViolationReportCommentDocument) TableName() string {
	return "violation_report_comment_documents"
}
