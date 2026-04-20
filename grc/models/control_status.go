package models

type ControlStatus struct {
	ID        uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	TitleEn   string `json:"title_en" binding:"required,max=500" gorm:"type:nvarchar(500);not null"`
	TitleAr   string `json:"title_ar" binding:"required,max=500" gorm:"type:nvarchar(500);not null"`
	ColorCode string `json:"color_code" binding:"required,max=25" gorm:"type:varchar(25);not null"`
	Type      string `json:"type" binding:"required,max=100,oneof=needs_improvement doesnt_need_improvement" gorm:"type:varchar(100);not null"`
	Status    string `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);default:active;not null"`
}

// TableName overrides the table name used by ControlStatus
func (ControlStatus) TableName() string {
	return "control_statuses"
}
