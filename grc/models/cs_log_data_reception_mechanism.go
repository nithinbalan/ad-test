package models

type CSLogDataReceptionMechanism struct {
	ID      uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn string `json:"title_en" binding:"required" gorm:"type:nvarchar(100);not null"`
	TitleAr string `json:"title_ar" binding:"required" gorm:"type:nvarchar(100);not null"`
	Status  string `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null"`
}

func (CSLogDataReceptionMechanism) TableName() string {
	return "cs_log_data_reception_mechanisms"
}
