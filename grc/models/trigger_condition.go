package models

type TriggerCondition struct {
	ID      uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn string `json:"title_en" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	TitleAr string `json:"title_ar" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	Type    string `json:"type" binding:"required,max=25,oneof='before-due-date on-due-date after-due-date'" gorm:"type:varchar(25);not null"`
	Status  string `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null"`
}

// TableName overrides the table name used by TriggerCondition to `trigger_conditions`
func (TriggerCondition) TableName() string {
	return "trigger_conditions"
}
