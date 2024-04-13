package models

type NotificationMaster struct {
	ID         uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	ModuleID   uint   `json:"modules_id" gorm:"not null"  binding:"required"`
	Identifier string `json:"identifier" gorm:"size:100;not null" binding:"required"`
	TitleEn    string `json:"title_en" gorm:"size:100;not null" binding:"required"`
	TitleAr    string `json:"title_ar" gorm:"size:100;not null" binding:"required"`

	Module Module `gorm:"foreignKey:ModulesID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}
