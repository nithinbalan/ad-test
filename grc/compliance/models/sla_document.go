package models

import (
	organizationModel "jethur/grc/organization/models"
	"time"
)

type SLADocument struct {
	ID            uint      `json:"id" gorm:"primary_key;auto_increment"`
	SLARegisterID uint      `json:"sla_register_id" binding:"required" gorm:"column:sla_register_id"`
	Version       string    `json:"version" binding:"required"`
	IssueDate     string    `json:"issue_date" binding:"omitempty" gorm:"type:date"`
	ExpiryDate    string    `json:"expiry_date" binding:"omitempty" gorm:"type:date"`
	Comment       string    `json:"comment,omitempty"`
	Language      string    `json:"language" binding:"required"`
	URL           string    `json:"url,omitempty"`
	Name          string    `json:"name,omitempty"`
	Size          int       `json:"size,omitempty"`
	Token         string    `json:"token,omitempty"`
	Type          string    `json:"type,omitempty"`
	Icon          string    `json:"icon" binding:"required"`
	IsLatest      bool      `json:"is_latest" gorm:"default:1"`
	SA1Date       string    `json:"sa1_date" binding:"omitempty" gorm:"type:date"`
	SA2Date       string    `json:"sa2_date" binding:"omitempty" gorm:"type:date"`
	CreatedBy     uint      `json:"created_by" gorm:"not null"`
	UpdatedBy     uint      `json:"updated_by" gorm:"default:null"`
	CreatedAt     time.Time `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"type:datetime;not null"`

	// Relationships
	SLARegister   *SLARegister            `gorm:"foreignKey:SLARegisterID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	CreatedByUser *organizationModel.User `gorm:"foreignKey:CreatedBy"`
	UpdatedByUser *organizationModel.User `gorm:"foreignKey:UpdatedBy"`
}

// TableName overrides the table name used by SLARegister to `sla_registers`
func (SLADocument) TableName() string {
	return "sla_documents"
}
