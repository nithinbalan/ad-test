package models

import (
	"time"
)

type BcpRole struct {
	ID                       uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	BusinessContinuityPlanID uint       `json:"business_continuity_plan_id" binding:"required"`
	ResponsibleUserID        uint       `json:"responsible_user_id" binding:"required"`
	ResponsibleUnitID        uint       `json:"responsible_unit_id" binding:"required"`
	Responsibility           string     `json:"responsibility" binding:"required" gorm:"type:nvarchar(max);not null"`
	BackupUserID             *uint      `json:"backup_user_id"`
	CreatedByID              uint       `json:"created_by_id" binding:"required"`
	CreatedAt                time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	BusinessContinuityPlan *BusinessContinuityPlan `gorm:"foreignKey:BusinessContinuityPlanID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResponsibleUser        *User                   `gorm:"foreignKey:ResponsibleUserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResponsibleUnit        *Unit                   `gorm:"foreignKey:ResponsibleUnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	BackupUser             *User                   `gorm:"foreignKey:BackupUserID"`
	CreatedByUser          *User                   `gorm:"foreignKey:CreatedByID"`
}


func (BcpRole) TableName() string {
	return "bcp_roles"
}
