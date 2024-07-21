package models

import (
	"time"
)

type DocumentSharedUnit struct {
	ID                 uint      `json:"id" binding:"required" gorm:"primaryKey;autoIncrement"`
	DocumentID         uint      `json:"document_id" binding:"required" gorm:"not null"`
	OrganizationUnitID uint      `json:"organization_unit_id" binding:"required" gorm:"not null"`
	CreatedBYID        uint      `json:"created_by_id" binding:"required" gorm:"not null"`
	CreatedAt          time.Time `json:"created_at" binding:"required" gorm:"type:datetime;not null"`

	// Relationships
	Document         *Document         `gorm:"foreignKey:DocumentID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	OrganizationUnit *OrganizationUnit `gorm:"foreignKey:OrganizationUnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser    *User             `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (DocumentSharedUnit) TableName() string {
	return "document_shared_units"
}
