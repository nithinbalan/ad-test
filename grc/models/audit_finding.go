package models

import (
	"time"
)

type AuditFinding struct {
	ID                     uint       `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	AuditID                uint       `json:"audit_id" gorm:"not null"`
	AuditableItemID        uint       `json:"auditable_item_id" gorm:"not null"`
	AuditFindingCategoryID uint       `json:"audit_finding_category_id" gorm:"not null"`
	SeverityLevelID        uint       `json:"severity_level_id" gorm:"not null"`
	Recommendations        *string    `json:"recommendations" gorm:"type:nvarchar(MAX);default:null"`
	AuditFindingStatusID   uint       `json:"audit_finding_status_id" gorm:"not null"`
	CreatedAt              time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt              *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedByID            uint       `json:"created_by_id" gorm:"not null"`
	UpdatedByID            *uint      `json:"updated_by_id" gorm:"default:null"`

	// Relationships
	Audit                *Audit                `gorm:"foreignKey:AuditID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	AuditableItem        *AuditableItem        `gorm:"foreignKey:AuditableItemID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	AuditFindingCategory *AuditFindingCategory `gorm:"foreignKey:AuditFindingCategoryID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	SeverityLevel        *SeverityLevel        `gorm:"foreignKey:SeverityLevelID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	AuditFindingStatus   *AuditFindingStatus   `gorm:"foreignKey:AuditFindingStatusID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	CreatedByUser        *User                 `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	UpdatedByUser        *User                 `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
}

// TableName overrides the table name used by AuditFinding to `audit_findings`
func (AuditFinding) TableName() string {
	return "audit_findings"
}
