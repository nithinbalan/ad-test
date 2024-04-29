package models

import (
	"time"
)

type ControlTestDocument struct {
	ID            uint      `json:"id" gorm:"primaryKey;autoIncrement" binding:"-"`
	ControlTestID uint      `json:"control_test_id" gorm:"not null" binding:"required"`
	URL           string    `json:"url" gorm:"type:varchar(1000);not null" binding:"required"`
	Name          string    `json:"name" gorm:"type:nvarchar(255);not null" binding:"required"`
	Size          uint      `json:"size" gorm:"not null" binding:"required"`
	Token         string    `json:"token" gorm:"type:varchar(255);not null" binding:"required"`
	Type          string    `json:"type" gorm:"type:varchar(255);not null" binding:"required"`
	Icon          string    `json:"icon" gorm:"type:varchar(10);not null" binding:"required,oneof=pdf excel ppt image text other"`
	CreatedAt     time.Time `json:"created_at" gorm:"type:datetime;not null" binding:"-"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"type:datetime" binding:"-"`

	//Relation ships
	ControlTest *ContolTest `gorm:"foreignKey:ControlTestID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (ControlTestDocument) TableName() string {
	return "control_test_documents"
}
