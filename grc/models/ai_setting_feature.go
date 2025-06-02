package models

type AIModuleFeature struct {
	ID                              uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	AIModuleID                      uint   `json:"ai_module_id"`
	TitleEn                         string `json:"title_en" gorm:"type:nvarchar(255);not null"`
	TitleAr                         string `json:"title_ar,omitempty" gorm:"type:nvarchar(255)"`
	DescriptionEn                   string `json:"description_en,omitempty" gorm:"type:nvarchar(max)"`
	DescriptionAr                   string `json:"description_ar,omitempty" gorm:"type:nvarchar(max)"`
	KnowledgeCenterResourceStatusID *int   `json:"knowledge_center_resource_status_id"`
	Order                           int    `json:"order" gorm:"default:0"`

	Module *AIModule                      `json:"-" gorm:"foreignKey:AIModuleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Status *KnowledgeCenterResourceStatus `json:"status" gorm:"foreignKey:KnowledgeCenterResourceStatusID;references:ID"`
}

func (AIModuleFeature) TableName() string {
	return "ai_module_features"
}
