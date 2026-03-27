package models

import (
	"time"
)

type InternalCommitteeVersionNomineeDepartmentParticipant struct {
	ID                                          uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	InternalCommitteeVersionNomineeDepartmentID uint       `json:"internal_committee_version_nominee_department_id" binding:"required" gorm:"not null"`
	UserID                                      uint       `json:"user_id" binding:"required" gorm:"not null"`
	DepartmentID                                uint       `json:"department_id" gorm:"not null"`
	CreatedBYID                                 *uint      `json:"created_by_id"`
	UpdatedBYID                                 *uint      `json:"updated_by_id"`
	CreatedAt                                   time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                                   *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	NomineeDepartment *InternalCommitteeVersionNomineeDepartment `gorm:"foreignKey:InternalCommitteeVersionNomineeDepartmentID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	User              *User                                      `gorm:"foreignKey:UserID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Department        *OrganizationUnit                          `gorm:"foreignKey:DepartmentID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser     *User                                      `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedByUser     *User                                      `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

// TableName overrides the table name
func (InternalCommitteeVersionNomineeDepartmentParticipant) TableName() string {
	return "internal_committee_version_nominee_department_participants"
}
