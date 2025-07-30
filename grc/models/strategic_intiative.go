package models

import "time"

type StrategicIntiative struct {
	ID                uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	TitleEn           string     `json:"title_en" gorm:"type:nvarchar(MAX);not null" binding:"required"`
	TitleAr           string     `json:"title_ar" gorm:"type:nvarchar(MAX);not null" binding:"required"`
	DescriptionEn     *string    `json:"description_en" gorm:"type:nvarchar(MAX);default:null"`
	DescriptionAr     *string    `json:"description_ar" gorm:"type:nvarchar(MAX);default:null"`
	StartDate         *time.Time `json:"start_date" gorm:"type:date;default:null"`
	EndDate           *time.Time `json:"end_date" gorm:"type:date;default:null"`
	ResponsibleUserID uint       `json:"responsible_user_id" gorm:"default:null"`
	ResponsibleUnitID uint       `json:"responsible_unit_id" gorm:"default:null"`
	CreatedAt         time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt         *time.Time `json:"updated_at" gorm:"type:datetime" binding:"omitempty"`
	CreatedBYID       uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID       *uint      `json:"updated_by_id" gorm:"default:null" binding:"omitempty"`

	// Relationships
	ResponsibleUser User             `gorm:"foreignKey:ResponsibleUserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	ResponsibleUnit OrganizationUnit `gorm:"foreignKey:ResponsibleUnitID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	CreatedByUser   User             `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser   *User            `gorm:"foreignKey:UpdatedBYID"`
}

func (StrategicIntiative) TableName() string {
	return "strategic_intiatives"
}
