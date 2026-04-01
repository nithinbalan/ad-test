package models

import "time"

type OrganizationUnitAttribute struct {
	ID           uint       `gorm:"primaryKey;autoIncrement"`
	DepartmentID uint       `json:"department_id" gorm:"not null"`
	Title        *string    `gorm:"type:text"`
	Description  *string    `gorm:"type:text"`
	CreatedBYID  uint       `json:"created_by_id" gorm:"not null;column:created_by_id" binding:"required"`
	UpdatedBYID  *uint      `json:"updated_by_id" gorm:"default:null;column:updated_by_id"`
	CreatedAt    time.Time  `json:"created_at" gorm:"type:datetime;not null" binding:"-"`
	UpdatedAt    *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	Department    *OrganizationUnit `json:"department" gorm:"foreignKey:DepartmentID;constraint:OnDelete:NO ACTION,OnUpdate:NO ACTION"`
	CreatedByUser *User             `json:"created_by_user" gorm:"foreignKey:CreatedBYID;constraint:OnDelete:NO ACTION,OnUpdate:NO ACTION"`
	UpdatedByUser *User             `json:"updated_by_user" gorm:"foreignKey:UpdatedBYID;constraint:OnDelete:NO ACTION,OnUpdate:NO ACTION"`
}

// TableName overrides the table name used by OrganizationUnit to `organization_units`
func (OrganizationUnitAttribute) TableName() string {
	return "organization_unit_attributes"
}
