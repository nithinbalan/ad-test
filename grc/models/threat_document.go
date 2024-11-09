package models

import (
	"time"
)

type ThreatDocument struct {
	ID           uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	ThreatID     uint       `json:"threat_id" gorm:"not null"`
	URL          *string    `json:"url" gorm:"type:varchar(1000);default:null"`
	Name         *string    `json:"name" gorm:"type:nvarchar(255);default:null"`
	Size         *int       `json:"size" gorm:"default:null"`
	Token        *string    `json:"token" gorm:"type:varchar(255);default:null"`
	OriginalName *string    `json:"original_name" gorm:"type:nvarchar(255);default:null"`
	Type         *string    `json:"type" gorm:"type:varchar(255);default:null"`
	CreatedBYID  uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID  *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt    time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt    *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	Threat        *Threat `gorm:"foreignKey:ThreatID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser *User   `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	UpdatedByUser *User   `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// TableName overrides the default table name.
func (ThreatDocument) TableName() string {
	return "threat_documents"
}
