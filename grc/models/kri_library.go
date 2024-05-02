package models

type KRILibrary struct {
	ID                uint    `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	KRICategoryID     uint    `json:"kri_category_id" binding:"required" gorm:"not null"`
	ReviewFrequencyID *uint   `json:"review_frequency_id" binding:"required" gorm:"default:null"`
	Title             string  `json:"title" binding:"required,max=255" gorm:"type:nvarchar(255);not null"`
	Description       *string `json:"description" gorm:"type:nvarchar(max)"`
	CalculationType   string  `json:"calculation_type" binding:"required,oneof=automatic, manual" gorm:"type:varchar(10);not null"`
	Source            *string `json:"source" gorm:"type:varchar(45)"`
	Status            string  `json:"status" binding:"required,oneof=active inactive" gorm:"type:varchar(10);not null"`

	// Relationships
	KRICategory     *KRICategory     `gorm:"foreignkey:KRICategoryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ReviewFrequency *ReviewFrequency `gorm:"foreignkey:ReviewFrequencyID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// TableName overrides the table name used by KRILibrary to `kri_library`
func (KRILibrary) TableName() string {
	return "kri_library"
}
