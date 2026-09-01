package models

type ProjectSecurityGateStage struct {
	ID        uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn   string `json:"title_en" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	TitleAr   string `json:"title_ar" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	Type      string `json:"type" binding:"required,max=30,oneof=pre_implementation during_implementation post_implementation" gorm:"type:varchar(30);not null"`
	ColorCode string `json:"color_code" binding:"required,max=25" gorm:"type:varchar(25);not null"`
	Status    string `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null"`
}

func (ProjectSecurityGateStage) TableName() string {
	return "project_security_gate_stages"
}
