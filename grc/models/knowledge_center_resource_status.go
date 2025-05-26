package models

type KnowledgeCenterResourceStatus struct {
	ID        uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn   string `json:"title_en" binding:"required" gorm:"type:nvarchar(100);not null"`
	TitleAr   string `json:"title_ar" binding:"required" gorm:"type:nvarchar(100);not null"`
	Status    string `json:"status" binding:"required" gorm:"type:nvarchar(50);not null"`
	ColorCode string `json:"color_code" gorm:"type:varchar(25);not null"`
}

func (KnowledgeCenterResourceStatus) TableName() string {
	return "knowledge_center_resource_statuses"
}
