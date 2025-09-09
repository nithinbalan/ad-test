package models

import (
	"time"
)

type BcpAsset struct {
	ID                       uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	BusinessContinuityPlanID uint      `json:"business_continuity_plan_id" binding:"required"`
	AssetRegisterID          uint      `json:"asset_register_id" binding:"required"`
	CreatedBYID              uint      `json:"created_by_id" binding:"required"`
	CreatedAt                time.Time `json:"created_at" gorm:"type:datetime;not null"`

	BusinessContinuityPlan *BusinessContinuityPlan `gorm:"foreignKey:BusinessContinuityPlanID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	AssetRegister          *AssetRegister          `gorm:"foreignKey:AssetRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser          *User                   `gorm:"foreignKey:CreatedBYID"`
}

func (BcpAsset) TableName() string {
	return "bcp_assets"
}
