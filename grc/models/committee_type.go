package models

type CommitteeType struct {
	ID      uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	TitleEn string `json:"title_en" binding:"required" gorm:"type:nvarchar(500);not null"`
	TitleAr string `json:"title_ar" binding:"required" gorm:"type:nvarchar(500);not null"`
	Type    string `json:"type" binding:"required,oneof=statutory_reference no_statutory_reference other" gorm:"type:varchar(100);not null;check:type IN ('statutory_reference','no_statutory_reference','other')"`
	Status  string `json:"status" binding:"required,oneof=active inactive" gorm:"type:varchar(10);not null;default:'active';check:status IN ('active','inactive')"`
}

// TableName overrides the default table name used by GORM
func (CommitteeType) TableName() string {
	return "committee_types"
}
