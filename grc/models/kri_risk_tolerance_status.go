package models

type KRIRiskToleranceStatus struct {
	ID        uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn   string `json:"title_en" binding:"required,max=45" gorm:"type:varchar(45);not null"`
	TitleAr   string `json:"title_ar" binding:"required,max=45" gorm:"type:nvarchar(45);not null"`
	ColorCode string `json:"color_code" binding:"required,max=25" gorm:"type:varchar(25);not null"`
	Type      string `json:"type" binding:"required,max=100,oneof=within_minimum_susceptibility within_maximum_susceptibility at_maximum_tolerable exceeded_maximum_tolerable" gorm:"type:varchar(100);not null"`
	Status    string `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null"`
}

// TableName overrides the table name used by KpiStatus to `kpi_statuses`
func (KRIRiskToleranceStatus) TableName() string {
	return "kri_risk_tolerance_statuses"
}
