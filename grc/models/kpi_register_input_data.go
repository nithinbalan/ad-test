package models

import "time"

type KPIRegisterDataInput struct {
	ID            uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	KPIRegisterID uint   `json:"kpi_register_id"`
	Variable      string `json:"variable"`
	Input         string `json:"input"`

	CreatedAt   time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt   *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedBYID uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID *uint      `json:"updated_by_id" gorm:"default:null"`

	// Relationships
	KPIRegister   *KPIRegister `gorm:"foreignKey:KPIRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser *User        `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser *User        `gorm:"foreignKey:UpdatedBYID"`
}

// TableName overrides the table name used by KPIRegisterDataInput to `kpi_register_input_data`
func (KPIRegisterDataInput) TableName() string {
	return "kpi_register_data_inputs"
}
