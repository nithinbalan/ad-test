package models

type Escalation struct {
	ID            uint    `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn       string  `json:"title_en" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	TitleAr       string  `json:"title_ar" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	DescriptionEn *string `json:"description_en"`
	DescriptionAr *string `json:"description_ar"`
	ModuleID      uint    `json:"module_id"`
	Status        string  `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null"`

	// Relationship
	EscalationLevels []EscalationLevel `gorm:"foreignKey:EscalationID;references:ID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

// TableName overrides the table name used by Escalation to `escalations`
func (Escalation) TableName() string {
	return "escalations"
}
