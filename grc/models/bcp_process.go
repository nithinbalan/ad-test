package models

import (
	"time"
)

type BcpProcess struct {
	ID                       uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	BusinessContinuityPlanID uint       `json:"incident_register_id" binding:"required"`
	ProcessID                uint       `json:"process_id" binding:"required"`
	RiskLevelID              *uint      `json:"risk_level_id"`
	Rto                      *uint      `json:"rto"`
	RtoUnitId                *uint      `json:"rto_unit_id"`
	Mtpd                     *uint      `json:"mtpd"`
	MtpdUnitId               *uint      `json:"mtpd_unit_id"`
	Rpo                      *uint      `json:"rpo"`
	RpoUnitId                *uint      `json:"rpo_unit_id"`
	CreatedBYID              uint       `json:"created_by_id" binding:"required"`
	UpdatedBYID              *uint      `json:"updated_by_id"`
	CreatedAt                time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	BusinessContinuityPlan *BusinessContinuityPlan `gorm:"foreignKey:BusinessContinuityPlanID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Process                *Process                `gorm:"foreignKey:ProcessID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RiskLevel              *RiskLevel              `gorm:"foreignKey:RiskLevelID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RtoUnit                *TimeUnit               `gorm:"foreignKey:RtoUnitId;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	MtpdUnit               *TimeUnit               `gorm:"foreignKey:MtpdUnitId;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RpoUnit                *TimeUnit               `gorm:"foreignKey:RpoUnitId;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser          *User                   `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser          *User                   `gorm:"foreignKey:UpdatedBYID"`
}

func (BcpProcess) TableName() string {
	return "bcp_processes"
}
