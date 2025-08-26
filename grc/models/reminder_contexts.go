package models

type ReminderContext struct {
	ID      uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Type    string `json:"type" binding:"required,max=100" gorm:"type:varchar(255);not null;unique"`
	TitleEn string `json:"title_en" binding:"required,max=500" gorm:"type:nvarchar(500);not null"`
	TitleAr string `json:"title_ar" binding:"required,max=500" gorm:"type:nvarchar(500);not null"`
}

// TableName overrides the default table name
func (ReminderContext) TableName() string {
	return "reminder_contexts"
}
