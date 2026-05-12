package models

type CommitteeMemberType struct {
	ID      uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	TitleEn string `json:"title_en" binding:"required" gorm:"type:nvarchar(500);not null"`
	TitleAr string `json:"title_ar" binding:"required" gorm:"type:nvarchar(500);not null"`
	Status  string `json:"status" binding:"required,oneof=active inactive" gorm:"type:varchar(10);not null;default:'active';check:status IN ('active','inactive')"`
}

// TableName overrides the default table name used by GORM
func (CommitteeMemberType) TableName() string {
	return "committee_member_types"
}
