package models

import "time"

type DocumentRequirement struct {
	ID                         uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	DocumentRequirementGroupID uint       `gorm:"index" json:"document_requirement_group_d"`
	TitleEn                    string     `gorm:"type:nvarchar(250);not null" json:"title_en"`
	TitleAr                    string     `gorm:"type:nvarchar(250);not null" json:"title_ar"`
	CreatedBYID                uint       `json:"created_by_id" binding:"required" gorm:"not null"`
	CreatedAt                  time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	UpdatedBYID                *uint      `json:"updated_by_id" gorm:"default:null"`
	UpdatedAt                  *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	// Relationships
	DocumentRequirementGroup DocumentRequirementGroup `gorm:"foreignKey:DocumentRequirementGroupID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser            *User                    `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	UpdatedByUser            *User                    `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// TableName overrides the default table name
func (DocumentRequirement) TableName() string {
	return "document_requirements"
}
