package models

import (
	"time"
)

type AuditPlan struct {
	ID              uint       `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	ReferenceNumber string     `json:"reference_number" gorm:"type:varchar(45);not null"`
	AuditCategoryID uint       `json:"audit_category_id" gorm:"not null"`
	AuditTeamID     *uint      `json:"audit_team_id" gorm:"default:null"`
	TeamLeadID      *uint      `json:"team_lead_id" gorm:"default:null"`
	Title           string     `json:"title" gorm:"type:nvarchar(255);not null"`
	Description     *string    `json:"description" gorm:"type:nvarchar(MAX);default:null"`
	StartDate       *time.Time `json:"start_date" gorm:"type:date;default:null"`
	TargetDate      *time.Time `json:"target_date" gorm:"type:date;default:null"`
	Purpose         *string    `json:"purpose" gorm:"type:nvarchar(MAX);default:null"`
	Objective       *string    `json:"objective" gorm:"type:nvarchar(MAX);default:null"`
	AuditStatusID   uint       `json:"audit_status_id" gorm:"not null"`
	CreatedAt       time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt       *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedBYID     uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID     *uint      `json:"updated_by_id" gorm:"default:null"`

	// Relationships
	AuditCategory       *AuditCategory            `gorm:"foreignKey:AuditCategoryID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	AuditTeam           *AuditTeam                `gorm:"foreignKey:AuditTeamID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	TeamLead            *User                     `gorm:"foreignKey:TeamLeadID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	AuditStatus         *AuditStatus              `gorm:"foreignKey:AuditStatusID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	Members             []*AuditPlanAuditMember   `gorm:"foreignKey:AuditPlanID"`
	AuditableItems      []*AuditPlanAuditableItem `gorm:"foreignKey:AuditPlanID"`
	AuditPlanChecklists []*AuditPlanChecklist     `gorm:"foreignKey:AuditPlanID"`
	CreatedByUser       *User                     `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser       *User                     `gorm:"foreignKey:UpdatedBYID"`
}

// TableName overrides the table name used by AuditPlan to `audit_plans`
func (AuditPlan) TableName() string {
	return "audit_plans"
}
