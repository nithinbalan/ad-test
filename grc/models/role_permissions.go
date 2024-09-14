package models

import "time"

type RolesPermissions struct {
	ID                   uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	RoleID               uint       `json:"role_id" gorm:"not null"`
	ModuleID             uint       `json:"module_id" gorm:"not null"`
	RolePermissionTypeID uint       `json:"role_permission_type_id" gorm:"not null"`
	IsGranted            bool       `json:"is_granted" gorm:"not null"`
	CreatedAt            time.Time  `json:"created_at" gorm:"autoCreateTime"`
	CreatedBy            int        `json:"created_by" gorm:"not null"`
	UpdatedAt            *time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	UpdatedBy            *int       `json:"updated_by"`

	RolePermissionType RolePermissionType `gorm:"foreignKey:RolePermissionTypeID;references:ID"`
	Role               Role               `gorm:"foreignKey:RoleID;references:ID"`
	Module             Module             `gorm:"foreignKey:ModuleID;references:ID"`
	CreatedByUser      User               `gorm:"foreignKey:CreatedBy;references:ID"`
	UpdatedByUser      User               `gorm:"foreignKey:UpdatedBy;references:ID"`
}
