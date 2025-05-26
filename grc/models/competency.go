package models

import (
	"time"
)

type Competency struct {
	ID                uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	CompetencyGroupID uint      `json:"competency_group_id" binding:"required"`
	TitleEn           string    `json:"title_en" binding:"required,max=255" gorm:"type:nvarchar(100);not null"`
	TitleAr           string    `json:"title_ar" binding:"required,max=255" gorm:"type:nvarchar(100);not null"`
	Status            string    `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null"`
	CreatedBYID       uint      `json:"created_by_id" gorm:"not null"`
	CreatedAt         time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	CompetencyGroup *CompetencyGroup `gorm:"foreignKey:CompetencyGroupID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser   *User            `gorm:"foreignKey:CreatedBYID"`
}

func (Competency) TableName() string {
	return "competencies"
}
