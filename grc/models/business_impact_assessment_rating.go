package models

import (
	"time"
)

type BusinessImpactAssessmentRating struct {
	ID                   uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ProcessID            uint       `json:"process_id" binding:"required" gorm:"column:process_id;not null"`
	OverallImpactLevelID uint       `json:"overall_impact_level_id" binding:"required" gorm:"column:overall_impact_level_id;not null"`
	BCMBIAImpactLevelID  uint       `json:"bcm_bia_impact_level_id" binding:"required" gorm:"column:bcm_bia_impact_level_id;not null"`
	IsEdited             bool       `json:"is_edited" gorm:"type:bit;not null;default:0"`
	RtoFromValue         *int       `json:"rto_from_value" gorm:"type:int;null"`
	RtoToValue           *int       `json:"rto_to_value" gorm:"type:int;null"`
	RtoCondition         *string    `json:"rto_condition" gorm:"type:varchar(20);null"`
	RtoUnitID            *uint      `json:"rto_unit_id" gorm:"null"`
	MtpdFromValue        *int       `json:"mtpd_from_value" gorm:"type:int;null"`
	MtpdToValue          *int       `json:"mtpd_to_value" gorm:"type:int;null"`
	MtpdCondition        *string    `json:"mtpd_condition" gorm:"type:varchar(20);null"`
	MtpdUnitID           *uint      `json:"mtpd_unit_id" gorm:"null"`
	RpoFromValue         *int       `json:"rpo_from_value" gorm:"type:int;null"`
	RpoToValue           *int       `json:"rpo_to_value" gorm:"type:int;null"`
	RpoCondition         *string    `json:"rpo_condition" gorm:"type:varchar(20);null"`
	RpoUnitID            *uint      `json:"rpo_unit_id" gorm:"null"`
	RspFromValue         *int       `json:"rsp_from_value" gorm:"type:int;null"`
	RspToValue           *int       `json:"rsp_to_value" gorm:"type:int;null"`
	RspCondition         *string    `json:"rsp_condition" gorm:"type:varchar(20);null"`
	RspUnitID            *uint      `json:"rsp_unit_id" gorm:"null"`
	OverrideReason       *string    `json:"override_reason" gorm:"type:nvarchar(max);null"`
	CreatedBYID          *uint      `json:"created_by_id" gorm:"column:created_by_id"`
	UpdatedBYID          *uint      `json:"updated_by_id"`
	CreatedAt            time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt            *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	Process            *Process                        `gorm:"foreignKey:ProcessID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	OverallImpactLevel *BCMImpactAnalysisLevel         `gorm:"foreignKey:OverallImpactLevelID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	BCMBIAImpactLevel  *BCMBIAImpactLevel              `gorm:"foreignKey:BCMBIAImpactLevelID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	RtoUnit            *BusinessContinuityDurationUnit `gorm:"foreignKey:RtoUnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	MtpdUnit           *BusinessContinuityDurationUnit `gorm:"foreignKey:MtpdUnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RpoUnit            *BusinessContinuityDurationUnit `gorm:"foreignKey:RpoUnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RspUnit            *BusinessContinuityDurationUnit `gorm:"foreignKey:RspUnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser      *User                           `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedByUser      *User                           `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

// TableName overrides the table name used by BusinessImpactAssessmentRating to `business_impact_assessment_ratings`
func (BusinessImpactAssessmentRating) TableName() string {
	return "business_impact_assessment_ratings"
}
