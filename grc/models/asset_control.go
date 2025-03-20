package models

import (
	"time"
)

type AssetControl struct {
	ID                            uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	AssetRegisterID               uint       `json:"asset_register_id" binding:"required"`
	ControlID                     uint       `json:"control_id" binding:"required"`
	ControlImplementationStatusID *uint      `json:"control_implementation_status_id" gorm:"default:null"`
	ControlEfficiencyStatusID     *uint      `json:"control_efficiency_status_id" gorm:"default:null"`
	Comments                      *string    `json:"comments" gorm:"type:varchar(max);default: null"`
	CreatedBYID                   uint       `json:"created_by_id" binding:"required"`
	CreatedAt                     time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedByID                   *uint      `json:"updated_by_id,omitempty" gorm:"default:null"`
	UpdatedAt                     *time.Time `json:"updated_at,omitempty" gorm:"type:datetime;default:null"`

	// Relationships
	AssetRegister               *AssetRegister               `gorm:"foreignKey:AssetRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Control                     *Control                     `gorm:"foreignKey:ControlID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ControlImplementationStatus *ControlImplementationStatus `gorm:"foreignKey:ControlImplementationStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ControlEfficiencyStatus     *ControlEfficiencyStatus     `gorm:"foreignKey:ControlEfficiencyStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser               *User                        `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser               *User                        `gorm:"foreignKey:UpdatedByID"`
}

func (AssetControl) TableName() string {
	return "asset_controls"
}
