package models

import (
	"time"
)

type AuditPlanCommencementLetter struct {
	ID                uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	AuditPlanID       uint       `json:"audit_plan_id" binding:"required" gorm:"not null"`
	ResponsibleUnitID uint       `json:"responsible_unit_id" binding:"required" gorm:"not null"`
	ToEmailIds        *string    `json:"to_email_ids" gorm:"type:text;null"`
	CcEmailIds        *string    `json:"cc_email_ids" gorm:"type:text;null"`
	Subject           string     `json:"subject" gorm:"type:nvarchar(255);not null"`
	Body              string     `json:"body" gorm:"type:text;not null"`
	FromUserID        uint       `json:"from_user_id" binding:"required" gorm:"not null"`
	CreatedByID       uint       `json:"created_by_id" binding:"required" gorm:"not null"`
	CreatedAt         time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	UpdatedByID       *uint      `json:"updated_by_id,omitempty" gorm:"default:null"`
	UpdatedAt         *time.Time `json:"updated_at,omitempty" gorm:"type:datetime;default:null"`

	// Relationships
	AuditPlan       *AuditPlan `gorm:"foreignKey:AuditPlanID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResponsibleUnit *Unit      `gorm:"foreignKey:ResponsibleUnitID"`
	FromUser        *User      `gorm:"foreignKey:FromUserID"`
	CreatedByUser   *User      `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser   *User      `gorm:"foreignKey:UpdatedByID"`
}

// TableName overrides the table name used by AuditPlanCommencementLetter to `audit_plan_commencement_letters`
func (AuditPlanCommencementLetter) TableName() string {
	return "audit_plan_commencement_letters"
}
