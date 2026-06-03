package models

type BCMBIAImpactLevel struct {
	ID                    uint    `json:"id" gorm:"primaryKey;autoIncrement"`
	ImpactAnalysisLevelID uint    `json:"impact_analysis_level_id" binding:"required" gorm:"not null"`
	ClassificationEn      string  `json:"classification_en" binding:"required" gorm:"type:varchar(50);not null"`
	ClassificationAr      string  `json:"classification_ar" binding:"required" gorm:"type:varchar(50);not null"`
	IsCritical            bool    `json:"is_critical" gorm:"type:bit;not null;default:0"`
	RtoFromValue          *int    `json:"rto_from_value" gorm:"type:int;null"`
	RtoToValue            *int    `json:"rto_to_value" gorm:"type:int;null"`
	RtoCondition          *string `json:"rto_condition" gorm:"type:varchar(20);null"`
	RtoUnitID             *uint   `json:"rto_unit_id" gorm:"null"`
	MtpdFromValue         *int    `json:"mtpd_from_value" gorm:"type:int;null"`
	MtpdToValue           *int    `json:"mtpd_to_value" gorm:"type:int;null"`
	MtpdCondition         *string `json:"mtpd_condition" gorm:"type:varchar(20);null"`
	MtpdUnitID            *uint   `json:"mtpd_unit_id" gorm:"null"`
	RpoFromValue          *int    `json:"rpo_from_value" gorm:"type:int;null"`
	RpoToValue            *int    `json:"rpo_to_value" gorm:"type:int;null"`
	RpoCondition          *string `json:"rpo_condition" gorm:"type:varchar(20);null"`
	RpoUnitID             *uint   `json:"rpo_unit_id" gorm:"null"`
	RspFromValue          *int    `json:"rsp_from_value" gorm:"type:int;null"`
	RspToValue            *int    `json:"rsp_to_value" gorm:"type:int;null"`
	RspCondition          *string `json:"rsp_condition" gorm:"type:varchar(20);null"`
	RspUnitID             *uint   `json:"rsp_unit_id" gorm:"null"`
	Order                 int     `json:"order" gorm:"type:int;not null"`
	Status                string  `json:"status" binding:"required,oneof=active inactive" gorm:"type:varchar(20);not null;default:'active'"`

	// Relationships
	ImpactAnalysisLevel *BCMImpactAnalysisLevel         `gorm:"foreignKey:ImpactAnalysisLevelID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RtoUnit             *BusinessContinuityDurationUnit `gorm:"foreignKey:RtoUnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	MtpdUnit            *BusinessContinuityDurationUnit `gorm:"foreignKey:MtpdUnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RpoUnit             *BusinessContinuityDurationUnit `gorm:"foreignKey:RpoUnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RspUnit             *BusinessContinuityDurationUnit `gorm:"foreignKey:RspUnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (BCMBIAImpactLevel) TableName() string {
	return "bcm_bia_impact_levels"
}
