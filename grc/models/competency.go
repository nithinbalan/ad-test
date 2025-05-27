package models

type Competency struct {
	ID                uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	CompetencyGroupID uint   `json:"competency_group_id" binding:"required"`
	TitleEn           string `json:"title_en" binding:"required,max=255" gorm:"type:nvarchar(100);not null"`
	TitleAr           string `json:"title_ar" binding:"required,max=255" gorm:"type:nvarchar(100);not null"`
	Status            string `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null"`

	// Relationships
	CompetencyGroup *CompetencyGroup `gorm:"foreignKey:CompetencyGroupID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (Competency) TableName() string {
	return "competencies"
}
