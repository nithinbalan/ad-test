package models

import (
	"time"
)

type AuditPlanAuditableItemRisk struct {
	ID                       uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	AuditPlanAuditableItemID uint       `json:"audit_plan_auditable_item_id" gorm:"not null;index" binding:"required"`
	RiskRegisterID           uint       `json:"risk_register_id" gorm:"not null;index" binding:"required"`
	CreatedByID              uint       `json:"created_by_id" gorm:"not null"`
	UpdatedByID              *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt                time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	AuditPlanAuditableItem *AuditPlanAuditableItem `gorm:"foreignKey:AuditPlanAuditableItemID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	RiskRegister           *RiskRegister           `gorm:"foreignKey:RiskRegisterID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	CreatedByUser          *User                   `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	UpdatedByUser          *User                   `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
}

// TableName overrides the table name used by GORM to `audit_plan_auditable_item_risks`
func (AuditPlanAuditableItemRisk) TableName() string {
	return "audit_plan_auditable_item_risks"
}
