package models

import (
	"time"
)

type ThirdPartyRisk struct {
	ID                   uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ThirdPartyRegisterID uint      `json:"third_party_register_id" binding:"required"`
	RiskRegisterID       uint      `json:"risk_register_id" binding:"required"`
	CreatedBYID          uint      `json:"created_by_id" binding:"required"`
	CreatedAt            time.Time `json:"created_at" gorm:"type:datetime;not null"`

	ThirdPartyRegister *ThirdPartyRegister `gorm:"foreignKey:ThirdPartyRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RiskRegister       *RiskRegister       `gorm:"foreignKey:RiskRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser      *User               `gorm:"foreignKey:CreatedBYID"`
}

func (ThirdPartyRisk) TableName() string {
	return "third_party_risks"
}
