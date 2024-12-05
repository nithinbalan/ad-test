package models

import (
	"time"
)

type Audit struct {
	ID            uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	AuditPlanID   uint       `json:"audit_plan_id" binding:"required" gorm:"not null"`
	StartDate     time.Time  `json:"start_date" binding:"required" gorm:"type:date;not null"`
	EndDate       time.Time  `json:"end_date" binding:"required" gorm:"type:date;not null"`
	CreatedAt     time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt     *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedBYID   uint       `json:"created_by_id" binding:"required" gorm:"not null"`
	UpdatedBYID   *uint      `json:"updated_by_id,omitempty" gorm:"null"`
	AuditStatusID uint       `json:"audit_status_id" binding:"required" gorm:"not null"`

	// Relationships
	AuditPlan          *AuditPlan           `gorm:"foreignKey:AuditPlanID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	AuditChecklists    []*AuditChecklist    `gorm:"foreignKey:AuditID"`
	AuditDraftFindings []*AuditDraftFinding `gorm:"foreignKey:AuditID"`
	AuditFindings      []*AuditFinding      `gorm:"foreignKey:AuditID"`
	AuditStatus        *AuditStatus         `gorm:"foreignKey:AuditStatusID;;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	CreatedByUser      *User                `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser      *User                `gorm:"foreignKey:UpdatedBYID"`
}

// TableName overrides the default table name used by GORM
func (Audit) TableName() string {
	return "audits"
}
