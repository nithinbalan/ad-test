package models

type RiskUpdatePlanTimelineType struct {
	ID        uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn   string `json:"title_en" binding:"required,max=500" gorm:"type:nvarchar(500);not null"`
	TitleAr   string `json:"title_ar" binding:"required,max=500" gorm:"type:nvarchar(500);not null"`
	ColorCode string `json:"color_code" binding:"required,max=25" gorm:"type:varchar(25);not null"`
	Type      string `json:"type" binding:"required,max=500,oneof=risk_register_update_plan erm_department_evaluation introductory_meeting business_context self_awareness_training leadership_assessment risk_team_assessment risk_identification control_register_updated risk_assessment_workshop risk_workshop_assessment risk_analysis risk_assessment_completed risk_mitigation" gorm:"type:varchar(500);not null"`
}

// TableName overrides the table name used by RiskUpdatePlanTimelineType to `risk_update_plan_timeline_types`
func (RiskUpdatePlanTimelineType) TableName() string {
	return "risk_update_plan_timeline_types"
}
