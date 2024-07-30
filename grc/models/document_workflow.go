package models

import (
	"time"
)

// DocumentWorkflow represents the structure of our resource
type DocumentWorkflow struct {
	ID                uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"` // Primary key, auto-incremented
	DocumentVersionID uint       `json:"document_version_id" binding:"required" gorm:"not null;index;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Level             int        `json:"level" binding:"required" gorm:"not null"`
	ApprovalTypes     string     `json:"approval_types" binding:"required,oneof='anyone' 'all'" gorm:"type:varchar(255);not null;default:'anyone';check:approval_types IN ('anyone', 'all')"`
	CreatedBYID       uint       `json:"created_by_id" binding:"required" gorm:"not null"`
	UpdatedBYID       *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt         time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	UpdatedAt         *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relations
	DocumentVersion *DocumentVersion `gorm:"foreignKey:DocumentVersionID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedBy       *User            `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	UpdatedBy       *User            `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// TableName overrides the table name used by User to `document_workflows`
func (DocumentWorkflow) TableName() string {
	return "document_workflows"
}
