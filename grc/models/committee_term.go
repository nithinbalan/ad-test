package models

type CommitteeTerm struct {
	ID        uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	TitleEn   string `json:"title_en" binding:"required" gorm:"type:nvarchar(500);not null"`
	TitleAr   string `json:"title_ar" binding:"required" gorm:"type:nvarchar(500);not null"`
	Type      string `json:"type" binding:"required,oneof=temporary permanent" gorm:"type:varchar(100);not null;check:type IN ('temporary','permanent')"`
	ColorCode string `json:"color_code" gorm:"type:varchar(10)"`
	Status    string `json:"status" binding:"required,oneof=active inactive" gorm:"type:varchar(10);not null;default:'active';check:status IN ('active','inactive')"`
}

// TableName overrides the default table name used by GORM
func (CommitteeTerm) TableName() string {
	return "committee_terms"
}
