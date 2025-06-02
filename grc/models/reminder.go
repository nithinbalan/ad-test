package models

type Reminder struct {
	ID       uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn  string `json:"title_en" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	TitleAr  string `json:"title_ar" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	ModuleID uint   `json:"module_id"`
	Type     string `json:"type" binding:"required,max=25,oneof='user-task'" gorm:"type:varchar(25);not null"`
	Status   string `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null"`

	// Relationship
	ReminderRules []ReminderRule `gorm:"foreignKey:ReminderID;references:ID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

// TableName overrides the table name used by Reminder to `reminders`
func (Reminder) TableName() string {
	return "reminders"
}
