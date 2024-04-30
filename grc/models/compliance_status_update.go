package models

import (
	"time"
)

type ComplianceStatusUpdate struct {
	ID                   uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ComplianceRegisterID uint       `json:"compliance_register_id" gorm:"not null"`
	ComplianceStatusID   uint       `json:"compliance_status_id" gorm:"not null"`
	Comment              string     `json:"comment" binding:"required,max=100" gorm:"type:nvarchar(max);not null"`
	CreatedAt            time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt            *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedBy            uint       `json:"created_by" gorm:"not null"`
	UpdatedBy            *uint      `json:"updated_by" gorm:"default:null"`

	//Relationships
	ComplianceRegister *ComplianceRegister `gorm:"foreignkey:ComplianceRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ComplianceStatus   *ComplianceStatus   `gorm:"foreignkey:ComplianceStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser      *User               `gorm:"foreignKey:CreatedBy"`
	UpdatedByUser      *User               `gorm:"foreignKey:UpdatedBy"`
}

// TableName overrides the table name used by ComplianceStatusUpdate to `compliance_status_updates`
func (ComplianceStatusUpdate) TableName() string {
	return "compliance_status_updates"
}
