package models

import "time"

type ComplianceLibraryVersion struct {
	ID                  uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ComplianceLibraryID uint       `json:"compliance_library_id"`
	Version             string     `json:"version" gorm:"type:nvarchar(50);not null"`
	Year                int        `json:"year" gorm:"not null"`
	Overview            string    `json:"overview" gorm:"type:nvarchar(max)"`
	Introduction        string    `json:"introduction" gorm:"type:nvarchar(max)"`
	Domains             string    `json:"domains" gorm:"type:nvarchar(max)"`
	CreatedAt           time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt           *time.Time `json:"updated_at" gorm:"type:datetime"`

	Library *ComplianceLibrary `json:"library" gorm:"foreignKey:ComplianceLibraryID;constraint:OnDelete:CASCADE"`
}

func (ComplianceLibraryVersion) TableName() string {
	return "compliance_library_versions"
}
