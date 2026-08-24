package models

type KPIStatus struct {
	ID            uint    `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn       string  `json:"title_en" binding:"required,max=45" gorm:"type:varchar(45);not null"`
	TitleAr       string  `json:"title_ar" binding:"required,max=45" gorm:"type:nvarchar(45);not null"`
	ColorCode     string  `json:"color_code" binding:"required,max=25" gorm:"type:varchar(25);not null"`
	Type          string  `json:"type" binding:"required,max=30,oneof=exceeded near-target needs-improvement low active inactive archived on-track at-risk achieved new" gorm:"type:varchar(30);not null"`
	Status        string  `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null"`
	Label         *string `json:"label" gorm:"type:varchar(255)"`
	IsNotAchieved bool    `json:"is_not_achieved" gorm:"type:bit;not null;default:0"`
	Order         *int    `json:"order" gorm:"column:order;type:int"`
}

// TableName overrides the table name used by KpiStatus to `kpi_statuses`
func (KPIStatus) TableName() string {
	return "kpi_statuses"
}
