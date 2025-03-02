package models

import (
	"time"
)

type BusinessImpactAssessmentRecoveryObjective struct {
	ID                uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ProcessID         *uint      `json:"process_id"`
	AssetRegisterID   *uint      `json:"asset_register_id"`
	DocumentID        *uint      `json:"document_id"`
	RoleInProcess     *string    `json:"role_in_process" gorm:"type:nvarchar(max)"`
	Rto               *int       `json:"rto"`
	RtoUnitId         *uint      `json:"rto_unit_id"`
	Mtpd              *int       `json:"mtpd"`
	MtpdUnitId        *uint      `json:"mtpd_unit_id"`
	Rpo               *int       `json:"rpo"`
	RpoUnitId         *uint      `json:"rpo_unit_id"`
	ImpactDescription *string    `json:"impact_description"`
	CreatedByID       uint       `json:"created_by_id" binding:"required"`
	UpdatedByID       *uint      `json:"updated_by_id"`
	CreatedAt         time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt         *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	Process       *Process       `gorm:"foreignKey:ProcessID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	AssetRegister *AssetRegister `gorm:"foreignKey:AssetRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RtoUnit       *TimeUnit      `gorm:"foreignKey:RtoUnitId;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	MtpdUnit      *TimeUnit      `gorm:"foreignKey:MtpdUnitId;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RpoUnit       *TimeUnit      `gorm:"foreignKey:RpoUnitId;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser *User          `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser *User          `gorm:"foreignKey:UpdatedByID"`
}

// TableName overrides the table name used by BusinessFunction to `business_functions`
func (BusinessImpactAssessmentRecoveryObjective) TableName() string {
	return "business_impact_assessment_recovery_objectives"
}
