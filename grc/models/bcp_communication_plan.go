package models

import (
	"time"
)

type BcpCommunicationPlan struct {
	ID                       uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	BusinessContinuityPlanID uint       `json:"business_continuity_id" binding:"required"`
	BcpPhaseID               uint       `json:"bcp_phase_id" binding:"required"`
	Title                    string     `json:"title"`
	ResponsibleUserID        uint       `json:"responsible_user_id" gorm:"not null" binding:"required"`
	ResponsibleUnitID        uint       `json:"responsible_unit_id" gorm:"not null" binding:"required"`
	Recipient                *string    `json:"recipient"`
	Frequency                *string    `json:"frequency"`
	Channel                  *string    `json:"channel"`
	CreatedBYID              uint       `json:"created_by_id" binding:"required"`
	UpdatedBYID              *uint      `json:"updated_by_id"`
	CreatedAt                time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	BusinessContinuityPlan *BusinessContinuityPlan `gorm:"foreignKey:BusinessContinuityPlanID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	BcpPhase               *BcpPhase               `gorm:"foreignKey:BcpPhaseID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResponsibleUser        *User                   `gorm:"foreignKey:ResponsibleUserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResponsibleUnit        *OrganizationUnit       `gorm:"foreignKey:ResponsibleUnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser          *User                   `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser          *User                   `gorm:"foreignKey:UpdatedBYID"`
}

func (BcpCommunicationPlan) TableName() string {
	return "bcp_communication_plans"
}
