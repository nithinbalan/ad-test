package models

import (
	"time"
)

type ThirdPartyProcess struct {
	ID                   uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ProcessID            uint      `json:"process_id" binding:"required"`
	ThirdPartyRegisterID uint      `json:"third_party_register_id" binding:"required"`
	CreatedByID          uint      `json:"created_by_id" binding:"required"`
	CreatedAt            time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	Process            *Process            `gorm:"foreignKey:ProcessID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ThirdPartyRegister *ThirdPartyRegister `gorm:"foreignKey:ThirdPartyRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser      *User               `gorm:"foreignKey:CreatedByID"`
}

func (ThirdPartyProcess) TableName() string {
	return "third_party_processes"
}
