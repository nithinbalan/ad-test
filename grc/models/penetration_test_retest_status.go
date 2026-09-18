package models

type PenetrationTestRetestStatus struct {
	ID        uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn   string `json:"title_en" gorm:"type:nvarchar(255);not null"`
	TitleAr   string `json:"title_ar" gorm:"type:nvarchar(255);not null"`
	ColorCode string `json:"color_code" gorm:"type:varchar(25);not null"`
	Type      string `json:"type" gorm:"type:varchar(40);not null"`
	Status    string `json:"status" gorm:"type:varchar(10);not null;default:'active'"`
}

func (PenetrationTestRetestStatus) TableName() string {
	return "penetration_test_retest_statuses"
}
