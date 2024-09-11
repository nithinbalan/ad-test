package models

import (
	"time"
)

type AuditPlanAuditableItemRiskControl struct {
	ID                            uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	AuditPlanAuditableItemRisksID uint       `json:"audit_plan_auditable_item_risks_id" gorm:"not null;index" binding:"required"`
	ControlsID                    uint       `json:"controls_id" gorm:"not null;index" binding:"required"`
	CreatedByID                   uint       `json:"created_by_id" gorm:"not null"`
	UpdatedByID                   *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt                     time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                     *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	AuditPlanAuditableItemRisk *AuditPlanAuditableItemRisk `gorm:"foreignKey:AuditPlanAuditableItemRisksID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	Control                    *Control                    `gorm:"foreignKey:ControlsID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	CreatedByUser              *User                       `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	UpdatedByUser              *User                       `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
}

// TableName overrides the table name used by GORM to `audit_plan_auditable_item_risk_controls`
func (AuditPlanAuditableItemRiskControl) TableName() string {
	return "audit_plan_auditable_item_risk_controls"
}
