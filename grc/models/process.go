package models

import (
	"time"
)

type Process struct {
	ID                uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	ProcessGroupID    uint       `json:"process_group_id" gorm:"not null" binding:"required"`
	OwnerUnitID       uint       `json:"owner_unit_id" gorm:"not null" binding:"required"`
	OwnerID           uint       `json:"owner_id" gorm:"not null" binding:"required"`
	LifecycleStatusID uint       `json:"lifecycle_status_id" gorm:"not null" binding:"required"`
	ReferenceNumber   string     `json:"reference_number" gorm:"type:varchar(45);not null;unique" binding:"required"`
	TitleEn           string     `json:"title_en" gorm:"type:nvarchar(500);not null" binding:"required"`
	TitleAr           string     `json:"title_ar" gorm:"type:nvarchar(500);not null" binding:"required"`
	Description       *string    `json:"description" gorm:"type:nvarchar(MAX);default:null" binding:"omitempty"`
	CreatedByID       uint       `json:"created_by_id" gorm:"not null;column:created_by_id" binding:"required"`
	UpdatedByID       *uint      `json:"updated_by_id" gorm:"default:null;column:updated_by_id" binding:"omitempty"`
	CreatedAt         time.Time  `json:"created_at" gorm:"type:datetime;not null" binding:"-"`
	UpdatedAt         *time.Time `json:"updated_at" gorm:"type:datetime;default:null" binding:"omitempty"`

	// Relationships
	ProcessGroup    *ProcessGroup     `gorm:"foreignKey:ProcessGroupID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	OwnerUnit       *OrganizationUnit `gorm:"foreignKey:OwnerUnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Owner           *User             `gorm:"foreignKey:OwnerID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	LifecycleStatus *LifecycleStatus  `gorm:"foreignKey:LifecycleStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser   *User             `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser   *User             `gorm:"foreignKey:UpdatedByID"`
	Controls        []*ProcessControl `gorm:"foreignKey:ProcessID"`
}

func (Process) TableName() string {
	return "processes"
}
