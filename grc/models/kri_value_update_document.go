package models

import (
	"time"
)

type KRIValueUpdateDocument struct {
	ID               uint      `json:"id" gorm:"primary_key;auto_increment"`
	KRIValueUpdateID uint      `json:"kri_value_update_id"`
	URL              string    `json:"url"`
	Name             string    `json:"name" binding:"required"`
	Size             int       `json:"size"`
	Token            string    `json:"token"`
	OriginalName     string    `json:"original_name"`
	Type             string    `json:"type"`
	CreatedAt        time.Time `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	KRIValueUpdate *KRIValueUpdate `gorm:"foreignKey:KRIValueUpdateID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
}

// TableName overrides the table name used by KRIValueUpdateDocument to `kri_value_update_documents`
func (KRIValueUpdateDocument) TableName() string {
	return "kri_value_update_documents"
}
