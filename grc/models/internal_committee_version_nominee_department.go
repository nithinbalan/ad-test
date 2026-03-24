package models

import (
	"time"
)

type InternalCommitteeVersionNomineeDepartment struct {
	ID                         uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	InternalCommitteeVersionID uint       `json:"internal_committee_version_id" binding:"required" gorm:"not null"`
	NomineeDepartmentID        uint       `json:"nominee_department_id" binding:"required" gorm:"not null"`
	MembershipType             *string    `json:"membership_type" gorm:"type:nvarchar(max)"`
	CreatedBYID                *uint      `json:"created_by_id"`
	UpdatedBYID                *uint      `json:"updated_by_id"`
	CreatedAt                  time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                  *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	InternalCommitteeVersion *InternalCommitteeVersion                              `gorm:"foreignKey:InternalCommitteeVersionID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	NomineeDepartment        *OrganizationUnit                                      `gorm:"foreignKey:NomineeDepartmentID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser            *User                                                  `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedByUser            *User                                                  `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Participants             []InternalCommitteeVersionNomineeDepartmentParticipant `gorm:"foreignKey:InternalCommitteeVersionNomineeDepartmentID;constraint:OnUpdate:CASCADE"`
}

// TableName overrides the table name
func (InternalCommitteeVersionNomineeDepartment) TableName() string {
	return "internal_committee_version_nominee_departments"
}
