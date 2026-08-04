package models

import "time"

type KPIValueUpdateDataInput struct {
	ID               uint     `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	KPIValueUpdateID uint     `json:"kpi_value_update_id"`
	Variable         string   `json:"variable"`
	Input            string   `json:"input"`
	Value            *float64 `json:"value" gorm:"type:decimal(18,2)"`

	CreatedAt   time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt   *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedBYID uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID *uint      `json:"updated_by_id" gorm:"default:null"`

	// Relationships
	KPIValueUpdate *KPIValueUpdate `gorm:"foreignKey:KPIValueUpdateID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser  *User           `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser  *User           `gorm:"foreignKey:UpdatedBYID"`
}

// TableName overrides the table name used by KPIValueUpdateDataInput to `kpi_value_update_input_data`
func (KPIValueUpdateDataInput) TableName() string {
	return "kpi_value_update_data_inputs"
}
