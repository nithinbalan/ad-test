package models

import (
	"time"
)

type KPIValueUpdateDocument struct {
	ID               uint      `json:"id" gorm:"primary_key;auto_increment"`
	KPIValueUpdateID uint      `json:"kpi_value_update_id"`
	URL              string    `json:"url"`
	Name             string    `json:"name" binding:"required"`
	Size             int       `json:"size"`
	Token            string    `json:"token"`
	OriginalName     string    `json:"original_name"`
	Type             string    `json:"type"`
	CreatedAt        time.Time `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	KPIValueUpdate *KPIValueUpdate `gorm:"foreignKey:KPIValueUpdateID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
}

// TableName overrides the table name used by KPIValueUpdateDocument to `kpi_value_update_documents`
func (KPIValueUpdateDocument) TableName() string {
	return "kpi_value_update_documents"
}
