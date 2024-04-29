package models

type RiskLevel struct {
	ID        uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Type      string `json:"type" binding:"required,max=10,oneof=high medium low" gorm:"type:varchar(10);not null"`
	TitleEn   string `json:"title_en" binding:"required,max=45" gorm:"type:nvarchar(45);not null"`
	TitleAr   string `json:"title_ar" binding:"required,max=45" gorm:"type:nvarchar(45);not null"`
	ColorCode string `json:"color_code" binding:"required,max=25" gorm:"type:varchar(25);not null"`
	Status    string `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null"`
}

// TableName overrides the table name used by RiskLevel to `risk_levels`
func (RiskLevel) TableName() string {
	return "risk_levels"
}
