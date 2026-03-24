package models

import (
	"time"
)

type InternalCommittee struct {
	ID                 uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	ReferenceNumber    string     `json:"reference_number" binding:"required,max=45" gorm:"type:varchar(45);not null;unique"`
	RequestingEntityID uint       `json:"requesting_entity_id" binding:"required" gorm:"not null"`
	CreatedBYID        uint       `json:"created_by_id" binding:"required" gorm:"not null"`
	UpdatedBYID        *uint      `json:"updated_by_id"`
	CreatedAt          time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt          *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	RequestingEntity           *OrganizationUnit            `gorm:"foreignKey:RequestingEntityID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser              *User                        `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedByUser              *User                        `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	InternalCommitteeVersions  []*InternalCommitteeVersion  `gorm:"foreignKey:InternalCommitteeID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	InternalCommitteeWorkflows []*InternalCommitteeWorkflow `gorm:"foreignKey:InternalCommitteeID"`
}

// TableName overrides the table name used by InternalCommittee
func (InternalCommittee) TableName() string {
	return "internal_committees"
}
