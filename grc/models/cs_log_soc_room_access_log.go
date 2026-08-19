package models

import (
	"time"
)

type CSLogSocRoomAccessLog struct {
	ID                   uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ReferenceNumber      string     `json:"reference_number" binding:"omitempty,max=45" gorm:"type:nvarchar(45);not null;unique"`
	VisitorFullName      string     `json:"visitor_full_name" binding:"required,max=255" gorm:"type:nvarchar(255);not null"`
	EntryTime            *time.Time `json:"entry_time" binding:"required" gorm:"type:datetime;not null"`
	ExitTime             *time.Time `json:"exit_time" binding:"required" gorm:"type:datetime;not null"`
	VisitorIDNumber      string     `json:"visitor_id_number" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	DepartmentID         uint       `json:"department_id" binding:"required" gorm:"not null"`
	ExternalPartyName    string     `json:"external_party_name" binding:"required,max=255" gorm:"type:nvarchar(255);not null"`
	PurposeOfEntry       string     `json:"purpose_of_entry" binding:"required" gorm:"type:nvarchar(max);not null"`
	AuthorizingOfficerID uint       `json:"authorizing_officer_id" binding:"required" gorm:"not null"`
	Notes                *string    `json:"notes" gorm:"type:nvarchar(max)"`
	CreatedBYID          uint       `json:"created_by_id" binding:"required" gorm:"column:created_by_id;not null"`
	UpdatedBYID          *uint      `json:"updated_by_id" gorm:"column:updated_by_id"`
	CreatedAt            time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt            *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	Department         *OrganizationUnit                `gorm:"foreignKey:DepartmentID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	AuthorizingOfficer *User                            `gorm:"foreignKey:AuthorizingOfficerID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser      *User                            `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedByUser      *User                            `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Documents          []*CSLogSocRoomAccessLogDocument `gorm:"foreignKey:SocRoomAccessLogID"`
}

func (CSLogSocRoomAccessLog) TableName() string {
	return "cs_log_soc_room_access_logs"
}
