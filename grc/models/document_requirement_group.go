package models

import "time"

type DocumentRequirementGroup struct {
	ID          uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	DocumentID  uint       `gorm:"index" json:"document_id"`
	TitleEn     string     `gorm:"type:nvarchar(250);not null" json:"title_en"`
	TitleAr     string     `gorm:"type:nvarchar(250);not null" json:"title_ar"`
	Status      string     `gorm:"type:varchar(10);not null;default:'active';check:status IN ('active', 'inactive')" json:"status"`
	CreatedBYID uint       `json:"created_by_id" binding:"required" gorm:"not null"`
	CreatedAt   time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	UpdatedBYID *uint      `json:"updated_by_id" gorm:"default:null"`
	UpdatedAt   *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	// Relationships
	Document      Document `gorm:"foreignKey:DocumentID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser *User    `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	UpdatedByUser *User    `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// TableName overrides the default table name
func (DocumentRequirementGroup) TableName() string {
	return "document_requirement_groups"
}
