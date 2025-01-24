package models

type TestResultOption struct {
	ID        uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	TitleEn   string `json:"title_en" binding:"required" gorm:"type:nvarchar(100);not null"`
	TitleAr   string `json:"title_ar" binding:"required" gorm:"type:nvarchar(100);not null"`
	ColorCode string `json:"color_code" binding:"required" gorm:"type:varchar(25);not null"`
	Type      string `json:"type" binding:"required,oneof=successful partially-successful failure" gorm:"type:varchar(10);not null;default:'successful';check:type IN ('successful', 'partially-successful', 'failure')"`
	Status    string `json:"status" binding:"required,oneof=active inactive" gorm:"type:varchar(10);not null;default:'active';check:status IN ('active', 'inactive')"`
}

func (TestResultOption) TableName() string {
	return "test_result_options"
}
