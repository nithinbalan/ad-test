package models

import (
	"time"
)

type KPIThirdParty struct {
	ID                   uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	KPIRegisterID        uint      `json:"kpi_register_id" binding:"required"`
	ThirdPartyRegisterID uint      `json:"third_party_register_id" binding:"required"`
	CreatedByID          uint      `json:"created_by_id" binding:"required"`
	CreatedAt            time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	KPIRegister        *KPIRegister        `gorm:"foreignKey:KPIRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ThirdPartyRegister *ThirdPartyRegister `gorm:"foreignKey:ThirdPartyRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser      *User               `gorm:"foreignKey:CreatedByID"`
}

func (KPIThirdParty) TableName() string {
	return "kpi_third_parties"
}
