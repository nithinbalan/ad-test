package models

type BCMImpactAnalysisLevel struct {
	ID        uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn   string `json:"title_en" binding:"required,max=45" gorm:"type:varchar(45);not null"`
	TitleAr   string `json:"title_ar" binding:"required,max=45" gorm:"type:nvarchar(45);not null"`
	Score     int    `json:"score" binding:"required" gorm:"type:int;not null"`
	ColorCode string `json:"color_code" binding:"required,max=25" gorm:"type:varchar(25);not null"`
	Type      string `json:"type" binding:"required,max=20,oneof=non_applicable very_low low moderate high very_high" gorm:"type:varchar(20);not null"`
	Status    string `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null"`
}

// TableName overrides the table name used by BCMImpactAnalysisLevel to `bcm_impact_analysis_levels`
func (BCMImpactAnalysisLevel) TableName() string {
	return "bcm_impact_analysis_levels"
}
