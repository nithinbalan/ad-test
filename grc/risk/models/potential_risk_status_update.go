package models

import "time"

type PotentialRiskStatusUpdate struct {
	ID                    uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	PotentialRiskID       uint      `json:"potential_risk_id" gorm:"not null"`
	PotentialRiskStatusID uint      `json:"potential_risk_status_id" gorm:"not null"`
	Comment               string    `json:"comment" binding:"required,max=100" gorm:"type:nvarchar(max);not null"`
	CreatedAt             time.Time `json:"created_at" gorm:"type:datetime;not null"`
	CreatedBy             uint      `json:"created_by" gorm:"not null"`

	//Relationships
	PotentialRisk       *PotentialRisk       `gorm:"foreignkey:PotentialRiskID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	PotentialRiskStatus *PotentialRiskStatus `gorm:"foreignkey:PotentialRiskStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// TableName overrides the table name used by PotentialRiskStatusUpdate to `potential_risk_status_updates`
func (PotentialRiskStatusUpdate) TableName() string {
	return "potential_risk_status_updates"
}
