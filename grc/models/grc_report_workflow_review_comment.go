package models

import "time"

type GRCReportWorkflowReviewComment struct {
	ID                        uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	GRCReportWorkflowID       uint       `json:"grc_report_workflow_id"`
	Comment                   string     `json:"comment"`
	GRCReportCommentStatusID  uint       `json:"grc_report_comment_status_id"`
	GRCReportWorkflowStatusID *uint      `json:"grc_report_workflow_status_id"`
	ResolvedBYID              *uint      `json:"resolved_by_id"`
	CreatedBYID               uint       `json:"created_by_id" binding:"required"`
	UpdatedBYID               *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt                 time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                 *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	CreatedByUser           *User                    `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedByUser           *User                    `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ResolvedUser            *User                    `gorm:"foreignKey:ResolvedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	GRCReportWorkflow       *GRCReportWorkflow       `gorm:"foreignKey:GRCReportWorkflowID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	GRCReportCommentStatus  *GRCReportCommentStatus  `gorm:"foreignKey:GRCReportCommentStatusID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	GRCReportWorkflowStatus *GRCReportWorkflowStatus `gorm:"foreignKey:GRCReportWorkflowStatusID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (GRCReportWorkflowReviewComment) TableName() string {
	return "grc_report_workflow_review_comments"
}
