package models

import "time"

type ComplianceFolder struct {
	ID          uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn     string     `json:"title_en" binding:"required,max=255" gorm:"type:nvarchar(255);not null"`
	TitleAr     string     `json:"title_ar" binding:"required,max=255" gorm:"type:nvarchar(255);not null"`
	Status      string     `json:"status" binding:"required,oneof=active inactive" gorm:"type:varchar(10);not null"`
	CreatedAt   time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt   *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedBYID uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID *uint      `json:"updated_by_id" gorm:"default:null"`
}

// TableName overrides the default table name
func (ComplianceFolder) TableName() string {
	return "compliance_folders"
}
