package models

import (
	"time"
)

type KPIValueUpdate struct {
	ID            uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	KPIRegisterID uint       `json:"kpi_register_id" binding:"required" gorm:"not null"`
	ReviewUserID  uint       `json:"review_user_id"`
	ReviewedDate  time.Time  `json:"reviewed_date" binding:"required" gorm:"type:date;not null"`
	UnitID        uint       `json:"unit_id" binding:"required" gorm:"not null"`
	Value         *float64   `json:"value" gorm:"type:int;not null"`
	Remarks       *string    `json:"remarks" gorm:"type:nvarchar(max)"`
	CreatedAt     time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt     *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedByID   uint       `json:"created_by_id" gorm:"not null"`
	UpdatedByID   *uint      `json:"updated_by_id" gorm:"default:null"`

	// Relationships
	KPIRegister   *KPIRegister               `gorm:"foreignKey:KPIRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Unit          *Unit                      `gorm:"foreignKey:UnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Documents     []KPIValueUpdateDocument   `gorm:"foreignKey:KPIValueUpdateID"`
	CreatedByUser *User                      `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser *User                      `gorm:"foreignKey:UpdatedByID"`
	ReviewUser    *User                      `gorm:"foreignKey:ReviewUserID"`
	DataInputs    []*KPIValueUpdateDataInput `gorm:"foreignKey:KPIValueUpdateID"`
}

// TableName overrides the table name used by KPIValueUpdate to `kpi_value_updates`
func (KPIValueUpdate) TableName() string {
	return "kpi_value_updates"
}
