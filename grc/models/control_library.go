package models

import "time"

type ControlLibrary struct {
	ID        uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	TitleEn   string     `json:"title_en" gorm:"type:nvarchar(255)"`
	TitleAr   string     `json:"title_ar" gorm:"type:nvarchar(255)"`
	CreatedAt time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"type:datetime"`

	Titles []ControlLibraryTitle `json:"titles" gorm:"foreignKey:ControlLibraryID;constraint:OnDelete:CASCADE"`
}

func (ControlLibrary) TableName() string {
	return "control_library"
}
