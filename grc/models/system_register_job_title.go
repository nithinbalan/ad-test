package models

import (
	"time"
)

type SystemRegisterJobTitle struct {
	ID                uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	SystemRegisterID  uint      `json:"system_register_id" binding:"required" gorm:"not null"`
	DesignationID     uint      `json:"designation_id" binding:"required" gorm:"not null"`
	CreatedBYID       uint      `json:"created_by_id" binding:"required" gorm:"not null"`
	CreatedAt         time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	SystemRegister    *SystemRegister `gorm:"foreignKey:SystemRegisterID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Designation       *Designation    `gorm:"foreignKey:DesignationID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser     *User           `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (SystemRegisterJobTitle) TableName() string {
	return "system_register_job_titles"
}
