package models

import (
	"time"
)

type ProcessActivityUser struct {
	ID                uint  `json:"id" gorm:"primaryKey;autoIncrement"`
	ProcessActivityID uint  `json:"process_activity_id" gorm:"not null"`
	ResponsibleUnitID *uint `json:"responsible_unit_id" gorm:"default:null"`

	UserID    uint      `json:"user_id" gorm:"not null"`
	RoleType  string    `json:"role_type" gorm:"type:varchar(50);not null;check:role_type IN ('responsible','accountable','consulted','informed')"`
	CreatedAt time.Time `json:"created_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP"`

	// Relationships
	ProcessActivity *ProcessActivity  `gorm:"foreignKey:ProcessActivityID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	User            *User             `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ResponsibleUnit *OrganizationUnit `gorm:"foreignKey:ResponsibleUnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (ProcessActivityUser) TableName() string {
	return "process_activity_users"
}
