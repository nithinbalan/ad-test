package models

import (
	"time"
)

type BcpRecoveryPlan struct {
	ID                       uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	BusinessContinuityPlanID uint       `json:"business_continuity_id" binding:"required"`
	BcpPhaseID               uint       `json:"bcp_phase_id" binding:"required"`
	BcpProcessID             uint       `json:"bcp_process_id" binding:"required"`
	SeverityLevelID          *uint      `json:"severity_level_id"`
	Title                    string     `json:"title"`
	ResponsibleUserID        uint       `json:"responsible_user_id" gorm:"not null" binding:"required"`
	ResponsibleUnitID        uint       `json:"responsible_unit_id" gorm:"not null" binding:"required"`
	TargetCompletion         *string    `json:"target_completion"`
	TargetCompletionUnitId   *uint      `json:"target_completion_unit_id"`
	CreatedBYID              uint       `json:"created_by_id" binding:"required"`
	UpdatedBYID              *uint      `json:"updated_by_id"`
	CreatedAt                time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	BusinessContinuityPlan *BusinessContinuityPlan `gorm:"foreignKey:BusinessContinuityPlanID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	BcpPhase               *BcpPhase               `gorm:"foreignKey:BcpPhaseID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	BcpProcess             *BcpProcess             `gorm:"foreignKey:BcpProcessID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	SeverityLevel          *SeverityLevel          `gorm:"foreignKey:SeverityLevelID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	TargetCompletionUnit   *TimeUnit               `gorm:"foreignKey:TargetCompletionUnitId;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResponsibleUser        *User                   `gorm:"foreignKey:ResponsibleUserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResponsibleUnit        *OrganizationUnit       `gorm:"foreignKey:ResponsibleUnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser          *User                   `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser          *User                   `gorm:"foreignKey:UpdatedBYID"`
}

func (BcpRecoveryPlan) TableName() string {
	return "bcp_recovery_plans"
}
