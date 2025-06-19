package models

import "time"

type ControlLibraryTitle struct {
	ID                  uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	ControlLibraryID    uint       `json:"control_library_id" gorm:"not null"`
	ComplianceLibraryID uint       `json:"compliance_library_id" gorm:"not null"`
	ControlCode         string     `json:"control_code" gorm:"type:nvarchar(255)"`
	ControlTitleEn      string     `json:"control_title_en" gorm:"type:nvarchar(max)"`
	ControlTitleAr      string     `json:"control_title_ar" gorm:"type:nvarchar(max)"`
	CreatedAt           time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt           *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relations (Optional)
	ControlLibrary    *ControlLibrary    `json:"control" gorm:"foreignKey:ControlLibraryID;references:ID"`
	ComplianceLibrary *ComplianceLibrary `json:"compliance" gorm:"foreignKey:ComplianceLibraryID;references:ID"`
}

func (ControlLibraryTitle) TableName() string {
	return "control_library_titles"
}
