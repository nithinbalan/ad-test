package models

import (
	"time"
)

type OrganizationUnitReport struct {
	ID                 uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	OrganizationUnitID uint      `json:"organization_unit_id" gorm:"not null"`
	Title              string    `json:"title" gorm:"type:nvarchar(255);not null"`
	StartDate          time.Time `json:"start_date" gorm:"type:date;not null"`
	EndDate            time.Time `json:"end_date" gorm:"type:date;not null"`
	CreatedByID        uint      `json:"created_by_id" gorm:"not null"`
	CreatedAt          time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	OrganizationUnit *OrganizationUnit `json:"organization_unit" gorm:"foreignKey:OrganizationUnitID;constraint:OnDelete:NO ACTION,OnUpdate:NO ACTION"`
	CreatedByUser    *User             `json:"created_by_user" gorm:"foreignKey:CreatedByID;constraint:OnDelete:NO ACTION,OnUpdate:NO ACTION"`
}

// TableName specifies the actual database table name
func (OrganizationUnitReport) TableName() string {
	return "organization_unit_reports"
}
