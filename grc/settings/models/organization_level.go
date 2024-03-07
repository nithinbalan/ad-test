package models

type OrganizationLevel struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	Level     uint8  `gorm:"not null;check:Level BETWEEN 1 AND 10"`
	TitleEn   string `gorm:"size:45;not null"`
	TitleAr   string `gorm:"size:45;not null"`
	IsPrimary bool   `gorm:"not null;default:false"`
	Type      string `gorm:"size:45;not null;check:Type IN ('division', 'department', 'section', 'subsection')"`
}

// TableName overrides the table name used by OrganizationLevel to `organization_levels`
func (OrganizationLevel) TableName() string {
	return "organization_levels"
}
