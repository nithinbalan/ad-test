package models

import (
	"time"
)

type KnowledgeCenterResource struct {
	ID                              uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ModuleID                        uint       `json:"module_id" binding:"required"`
	KnowledgeCenterResourceStatusID uint       `json:"knowledge_center_resource_status_id" binding:"required"`
	TitleEn                         string     `json:"title_en" binding:"required" gorm:"type:nvarchar(100);not null"`
	TitleAr                         string     `json:"title_ar" binding:"required" gorm:"type:nvarchar(100);not null"`
	Icon                            string     `json:"icon" gorm:"type:nvarchar(255)"`
	Type                            string     `json:"type" gorm:"type:varchar(50)"`
	ItemCount                       int        `json:"item_count" gorm:"default:0;not null"`
	CreatedAt                       time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                       *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	Module                        *Module                        `gorm:"foreignKey:ModuleID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	KnowledgeCenterResourceStatus *KnowledgeCenterResourceStatus `gorm:"foreignKey:KnowledgeCenterResourceStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (KnowledgeCenterResource) TableName() string {
	return "knowledge_center_resources"
}
