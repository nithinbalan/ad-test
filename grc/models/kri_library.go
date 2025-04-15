package models

type KRILibrary struct {
	ID          uint    `json:"id" gorm:"primaryKey;autoIncrement"`
	KRICategory string  `json:"kri_category" gorm:"type:nvarchar(255);not null"`
	Title       string  `json:"title" gorm:"type:nvarchar(500);not null"`
	Description *string `json:"description" gorm:"type:nvarchar(max)"`
	KRIUnit     *string `json:"kri_unit" gorm:"type:nvarchar(45)"`
}

func (KRILibrary) TableName() string {
	return "kri_library"
}
