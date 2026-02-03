package models

import (
	"time"
)

type ComplianceRegisterSharedUnit struct {
	ID                   uint      `json:"id" binding:"required" gorm:"primaryKey;autoIncrement"`
	ComplianceRegisterID uint      `json:"compliance_register_id" binding:"required" gorm:"not null"`
	DepartmentID         uint      `json:"department_id" binding:"required" gorm:"not null"`
	CreatedBYID          uint      `json:"created_by_id" binding:"required" gorm:"not null"`
	CreatedAt            time.Time `json:"created_at" binding:"required" gorm:"type:datetime;not null"`

	// Relationships
	ComplianceRegister *ComplianceRegister `gorm:"foreignKey:ComplianceRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Department         *OrganizationUnit   `gorm:"foreignKey:DepartmentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser      *User               `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (ComplianceRegisterSharedUnit) TableName() string {
	return "compliance_register_shared_units"
}
