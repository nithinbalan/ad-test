package models

import (
	"time"
)

type AuditDraftFinding struct {
	ID                        uint       `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	AuditID                   uint       `json:"audit_id" gorm:"not null"`
	Title                     string     `json:"title" binding:"required,max=255" gorm:"type:nvarchar(255);not null"`
	Description               *string    `json:"description" gorm:"type:nvarchar(max)"`
	AuditChecklistID          *uint      `json:"audit_checklist_id" gorm:"default:null"`
	AuditableItemID           *uint      `json:"auditable_item_id" gorm:"default:null"`
	AuditFindingCategoryID    uint       `json:"audit_finding_category_id" gorm:"not null"`
	SeverityLevelID           *uint      `json:"severity_level_id" gorm:"default:null"`
	Recommendations           *string    `json:"recommendations" gorm:"type:nvarchar(MAX);default:null"`
	AuditDraftFindingStatusID uint       `json:"audit_draft_finding_status_id" gorm:"not null"`
	CreatedAt                 time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                 *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedBYID               uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID               *uint      `json:"updated_by_id" gorm:"default:null"`

	// Relationships
	Audit                   *Audit                       `gorm:"foreignKey:AuditID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	AuditChecklist          *AuditChecklist              `gorm:"foreignKey:AuditChecklistID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	AuditableItem           *AuditableItem               `gorm:"foreignKey:AuditableItemID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	AuditFindingCategory    *AuditFindingCategory        `gorm:"foreignKey:AuditFindingCategoryID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	SeverityLevel           *SeverityLevel               `gorm:"foreignKey:SeverityLevelID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	AuditDraftFindingStatus *AuditDraftFindingStatus     `gorm:"foreignKey:AuditDraftFindingStatusID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	Documents               []*AuditDraftFindingDocument `gorm:"foreignKey:AuditDraftFindingID"`
	CreatedByUser           *User                        `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	UpdatedByUser           *User                        `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
}

// TableName overrides the table name used by AuditDraftFinding to `audit_draft_findings`
func (AuditDraftFinding) TableName() string {
	return "audit_draft_findings"
}
