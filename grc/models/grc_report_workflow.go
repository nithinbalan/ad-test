package models

import "time"

type GRCReportWorkflow struct {
	ID                        uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	GRCReportID               uint       `json:"grc_report_id"`
	Type                      string     `json:"type" gorm:"type:varchar(20);check:type IN ('risk', 'compliance', 'governance');not null"`
	GRCReportWorkflowStatusID uint       `json:"grc_report_workflow_status_id"`
	ReviewUserID              *uint      `json:"review_user_id"`
	ApproveUserID             *uint      `json:"approve_user_id"`
	ReviewedAt                *time.Time `json:"reviewed_at" gorm:"type:datetime"`
	ApprovedAt                *time.Time `json:"approved_at" gorm:"type:datetime"`
	CreatedBYID               uint       `json:"created_by_id" binding:"required"`
	UpdatedBYID               *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt                 time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                 *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	CreatedByUser                   *User                            `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedByUser                   *User                            `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ReviewUser                      *User                            `gorm:"foreignKey:ReviewUserID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	ApproveUser                     *User                            `gorm:"foreignKey:ApproveUserID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	GRCReport                       *GRCReport                       `gorm:"foreignKey:GRCReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	GRCReportWorkflowStatus         *GRCReportWorkflowStatus         `gorm:"foreignKey:GRCReportWorkflowStatusID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	GRCReportWorkflowReviewComments []GRCReportWorkflowReviewComment `gorm:"foreignKey:GRCReportWorkflowID;references:ID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (GRCReportWorkflow) TableName() string {
	return "grc_report_workflows"
}
