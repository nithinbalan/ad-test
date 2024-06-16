package models

import (
	"time"
)

type SLADocument struct {
	ID            uint       `json:"id" gorm:"primary_key;auto_increment"`
	SLARegisterID uint       `json:"sla_register_id"`
	Version       string     `json:"version" binding:"required"`
	IssueDate     *time.Time `json:"issue_date" binding:"omitempty" gorm:"type:date"`
	ExpiryDate    *time.Time `json:"expiry_date" binding:"omitempty" gorm:"type:date"`
	Comment       *string    `json:"comment"`
	URL           string     `json:"url"`
	Name          string     `json:"name" binding:"required"`
	OriginalName  string     `json:"original_name"`
	Size          int        `json:"size"`
	Token         string     `json:"token"`
	Type          string     `json:"type"`
	IsLatest      bool       `json:"is_latest" gorm:"default:1"`
	SA1Date       *time.Time `json:"sa1_date" binding:"omitempty" gorm:"type:date"`
	SA2Date       *time.Time `json:"sa2_date" binding:"omitempty" gorm:"type:date"`
	CreatedBy     uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBy     uint       `json:"updated_by_id" gorm:"default:null"`
	CreatedAt     time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt     time.Time  `json:"updated_at" gorm:"type:datetime;not null"`

	// Relationships
	SLARegister   *SLARegister `gorm:"foreignKey:SLARegisterID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	CreatedByUser *User        `gorm:"foreignKey:CreatedBy"`
	UpdatedByUser *User        `gorm:"foreignKey:UpdatedBy"`
}

// TableName overrides the table name used by SLARegister to `sla_register`
func (SLADocument) TableName() string {
	return "sla_documents"
}
