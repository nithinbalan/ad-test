package models

import (
	"time"
)

type AuditableItemOrganizationUnit struct {
	ID                 uint      `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	AuditableItemID    uint      `json:"auditable_item_id" gorm:"not null"`
	OrganizationUnitID uint      `json:"organization_unit_id" gorm:"not null"`
	CreatedBYID        uint      `json:"created_by_id"`
	CreatedAt          time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	AuditableItem    *AuditableItem    `gorm:"foreignKey:AuditableItemID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	OrganizationUnit *OrganizationUnit `gorm:"foreignKey:OrganizationUnitID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	CreatedByUser    *User             `gorm:"foreignKey:CreatedBYID"`
}

// TableName overrides the table name used by AuditableItemOrganizationUnit to `auditable_item_organization_units`
func (AuditableItemOrganizationUnit) TableName() string {
	return "auditable_item_organization_units"
}
