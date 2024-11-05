package models

import (
	"time"
)

type Threat struct {
	ID                uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	TitleEn           string     `json:"title_en" gorm:"type:nvarchar(500);not null" binding:"required"`
	TitleAr           string     `json:"title_ar" gorm:"type:nvarchar(500);not null" binding:"required"`
	ThreatTypeID      uint       `json:"threat_type_id" gorm:"not null" binding:"required"`
	SeverityLevelID   uint       `json:"severity_level_id" gorm:"not null" binding:"required"`
	ThreatStatusID    uint       `json:"threat_status_id" gorm:"not null" binding:"required"`
	ResponsibleUserID uint       `json:"responsible_user_id" gorm:"not null" binding:"required"`
	ResponsibleUnitID uint       `json:"responsible_unit_id" gorm:"not null" binding:"required"`
	DescriptionEn     *string    `json:"description_en,omitempty" gorm:"type:nvarchar(MAX);default:null"`
	DescriptionAr     *string    `json:"description_ar,omitempty" gorm:"type:nvarchar(MAX);default:null"`
	DetectionDate     *time.Time `json:"detection_date,omitempty" gorm:"type:date;default:null"`
	ResolutionDate    *time.Time `json:"resolution_date,omitempty" gorm:"type:date;default:null"`
	CreatedBYID       uint       `json:"created_by_id" gorm:"not null" binding:"required"`
	CreatedAt         time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedBYID       *uint      `json:"updated_by_id,omitempty" gorm:"default:null"`
	UpdatedAt         *time.Time `json:"updated_at,omitempty" gorm:"type:datetime;default:null"`

	// Relationships
	ThreatType      *ThreatType       `gorm:"foreignKey:ThreatTypeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	SeverityLevel   *SeverityLevel    `gorm:"foreignKey:SeverityLevelID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ThreatStatus    *ThreatStatus     `gorm:"foreignKey:ThreatStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResponsibleUser *User             `gorm:"foreignKey:ResponsibleUserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResponsibleUnit *OrganizationUnit `gorm:"foreignKey:ResponsibleUnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser   *User             `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	UpdatedByUser   *User             `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// TableName overrides the default table name.
func (Threat) TableName() string {
	return "threats"
}
