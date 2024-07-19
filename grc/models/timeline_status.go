package models

type TimelineStatus struct {
	ID        uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn   string `json:"title_en" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	TitleAr   string `json:"title_ar" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	ColorCode string `json:"color_code" binding:"required,max=25" gorm:"type:varchar(25);not null"`
	Type      string `json:"type" binding:"required,max=25,oneof='pending completed'" gorm:"type:varchar(25);not null;default:'pending'"`
}

// TableName sets the insert table name for this struct type
func (TimelineStatus) TableName() string {
	return "timeline_statuses"
}
