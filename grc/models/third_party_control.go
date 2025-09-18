package models

import (
	"time"
)

type ThirdPartyControl struct {
	ID                   uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ControlID            uint      `json:"control_id" binding:"required"`
	ThirdPartyRegisterID uint      `json:"third_party_register_id" binding:"required"`
	CreatedByID          uint      `json:"created_by_id" binding:"required"`
	CreatedAt            time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	Control            *Control            `gorm:"foreignKey:ControlID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ThirdPartyRegister *ThirdPartyRegister `gorm:"foreignKey:ThirdPartyRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser      *User               `gorm:"foreignKey:CreatedByID"`
}

func (ThirdPartyControl) TableName() string {
	return "third_party_controls"
}
