package models

type OrganizationUnit struct {
	ID                  uint               `gorm:"primaryKey;autoIncrement"`
	OrganizationLevelID uint               `gorm:"not null;index:,foreignKey:OrganizationLevelID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	OrganizationUnitID  *uint              `gorm:""`
	TitleEn             string             `gorm:"size:100;not null"`
	TitleAr             string             `gorm:"size:100;not null"`
	Status              string             `gorm:"size:10;not null;check:Status IN ('active', 'inactive')"`
	ParentUnit          *OrganizationUnit  `gorm:"foreignKey:OrganizationUnitID"`
	ChildUnits          []OrganizationUnit `gorm:"foreignKey:OrganizationUnitID"`
}

// TableName overrides the table name used by OrganizationUnit to `organization_units`
func (OrganizationUnit) TableName() string {
	return "organization_units"
}
