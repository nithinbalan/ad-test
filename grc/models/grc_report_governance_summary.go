package models

import "time"

type GRCReportGovernanceSummary struct {
	ID                    uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	GRCReportID           uint       `json:"grc_report_id" gorm:"not null"`
	ReviewUserID          *uint      `json:"review_user_id"`
	ApproveUserID         *uint      `json:"approve_user_id"`
	PoliciesInSystem      string     `json:"policies_in_system" gorm:"not null"`
	RolesDefined          int        `json:"roles_defined" gorm:"not null;default:0"`
	TotalPolicies         int        `json:"total_policies" gorm:"not null;default:0"`
	PolicyApprovedRate    *float64   `json:"policy_approved_rate" gorm:"type:decimal(5,2)"`
	Highlights            *string    `json:"highlights" gorm:"type:nvarchar(max)"`
	Challenges            *string    `json:"challenges" gorm:"type:nvarchar(max)"`
	DocumentOverviewChart *string    `json:"document_overview_chart" gorm:"type:nvarchar(max)"`
	DocumentTypeChart     *string    `json:"document_type_chart" gorm:"type:nvarchar(max)"`
	TopPublishedDocuments *string    `json:"top_published_documents" gorm:"type:nvarchar(max)"`
	CreatedBYID           uint       `json:"created_by_id" gorm:"not null"`
	CreatedAt             time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedBYID           *uint      `json:"updated_by_id"`
	UpdatedAt             *time.Time `json:"updated_at"`

	GRCReport     *GRCReport `gorm:"foreignKey:GRCReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser *User      `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser *User      `gorm:"foreignKey:UpdatedBYID"`
	ReviewUser    *User      `gorm:"foreignKey:ReviewUserID"`
	ApproveUser   *User      `gorm:"foreignKey:ApproveUserID"`
}

func (GRCReportGovernanceSummary) TableName() string {
	return "grc_report_governance_summary"
}
