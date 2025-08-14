package models

import (
	"time"
)

type BCPAssociatedDocument struct {
	ID                       uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	BusinessContinuityPlanID uint      `json:"business_continuity_plan_id" binding:"required"`
	DocumentID               uint      `json:"document_id" binding:"required"`
	CreatedByID              uint      `json:"created_by_id" binding:"required"`
	CreatedAt                time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	BusinessContinuityPlan *BusinessContinuityPlan `gorm:"foreignKey:BusinessContinuityPlanID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Document               *Document               `gorm:"foreignKey:DocumentID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser          *User                   `gorm:"foreignKey:CreatedByID"`
}

func (BCPAssociatedDocument) TableName() string {
	return "bcp_associated_documents"
}
