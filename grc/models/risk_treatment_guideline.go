package models

type RiskTreatmentGuideline struct {
	ID          uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	RiskLevelID uint      `json:"risk_level_id" gorm:"not null"`
	TimeframeEn string    `json:"timeframe_en" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	TimeframeAr string    `json:"timeframe_ar" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	Status      string    `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null;default:'active'"`
	RiskLevel   RiskLevel `json:"risk_level" gorm:"foreignKey:RiskLevelID;references:ID"`
}

// TableName overrides the table name used by RiskTreatmentGuideline to `risk_treatment_guidelines`
func (RiskTreatmentGuideline) TableName() string {
	return "risk_treatment_guidelines"
}
