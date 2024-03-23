package models

type OrganizationLevel struct {
	ID        uint   `json:"id" gorm:"primary_key;auto_increment"`
	Level     uint8  `json:"level" gorm:"not null;check:level BETWEEN 1 AND 10"`
	TitleEn   string `json:"title_en" binding:"required,max=100" gorm:"type:varchar(45);not null"`
	TitleAr   string `json:"title_ar" binding:"required,max=100" gorm:"type:nvarchar(45);not null"`
	IsPrimary bool   `json:"is_primary" gorm:"not null;default:false"`
}

// TableName overrides the table name used by OrganizationLevel to `organization_levels`
func (OrganizationLevel) TableName() string {
	return "organization_levels"
}
