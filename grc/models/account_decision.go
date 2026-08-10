package models

type AccountDecision struct {
	ID        uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn   string `json:"title_en" gorm:"type:nvarchar(100);not null"`
	TitleAr   string `json:"title_ar" gorm:"type:nvarchar(100);not null"`
	ColorCode string `json:"color_code" gorm:"type:varchar(25);not null"`
	Type      string `json:"type" binding:"required,oneof=keep revoke modify investigate" gorm:"type:varchar(20);not null"`
	Status    string `json:"status" gorm:"type:varchar(10);not null"`
}

func (AccountDecision) TableName() string {
	return "account_decisions"
}
