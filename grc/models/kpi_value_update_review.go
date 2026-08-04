package models

import "time"

type KPIValueUpdateReview struct {
	ID                uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	KPIValueUpdateID  uint      `json:"kpi_value_update_id" binding:"required" gorm:"not null"`
	KPIReviewStatusID uint      `json:"kpi_review_status_id" binding:"required" gorm:"not null"`
	Comment           *string   `json:"comment" gorm:"type:nvarchar(max)"`
	CreatedBYID       uint      `json:"created_by_id" binding:"required" gorm:"not null"`
	CreatedAt         time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	KPIValueUpdate  *KPIValueUpdate  `gorm:"foreignKey:KPIValueUpdateID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
	KPIReviewStatus *KPIReviewStatus `gorm:"foreignKey:KPIReviewStatusID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
	CreatedByUser   *User            `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
}

func (KPIValueUpdateReview) TableName() string {
	return "kpi_value_update_reviews"
}
