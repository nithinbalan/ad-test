package models

type MeetingMinuteStatus struct {
	ID        int    `json:"id" gorm:"primaryKey;autoIncrement"`
	TitleEn   string `json:"title_en" binding:"required" gorm:"type:nvarchar(100);not null"`
	TitleAr   string `json:"title_ar" binding:"required" gorm:"type:nvarchar(100);not null"`
	Type      string `json:"type" binding:"required,oneof=scheduled cancelled completed" gorm:"type:varchar(10);not null;check:type IN ('scheduled','cancelled','completed')"`
	ColorCode string `json:"color_code" binding:"required" gorm:"type:varchar(10);not null"`
	Status    string `json:"status" binding:"required,oneof=active inactive" gorm:"type:varchar(10);not null;default:'active';check:status IN ('active', 'inactive')"`
}

func (MeetingMinuteStatus) TableName() string {
	return "meeting_minute_statuses"
}
