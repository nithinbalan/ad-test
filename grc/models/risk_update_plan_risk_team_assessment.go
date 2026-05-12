package models

import (
	"time"
)

type RiskUpdatePlanRiskTeamAssessment struct {
	ID                       uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	RiskUpdatePlanTimelineID uint       `json:"risk_update_plan_timeline_id" binding:"required" gorm:"not null"`
	RiskChampionID           uint       `json:"risk_champion_id" binding:"required" gorm:"not null"`
	CreatedBYID              *uint      `json:"created_by_id" gorm:"default:null"`
	UpdatedBYID              *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt                time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	RiskUpdatePlanTimeline *RiskUpdatePlanTimeline `gorm:"foreignKey:RiskUpdatePlanTimelineID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	RiskChampion           *User                   `gorm:"foreignKey:RiskChampionID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser          *User                   `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedByUser          *User                   `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (RiskUpdatePlanRiskTeamAssessment) TableName() string {
	return "risk_update_plan_risk_team_assessments"
}
