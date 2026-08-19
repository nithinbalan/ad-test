package models

import "time"

type IncidentEvidence struct {
	ID                       uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	IncidentRegisterID       uint       `json:"incident_register_id" gorm:"not null"`
	IncidentEvidenceSourceID uint       `json:"incident_evidence_source_id" gorm:"not null"`
	IncidentEvidenceTypeID   uint       `json:"incident_evidence_type_id" gorm:"not null"`
	URL                      *string    `json:"url,omitempty" gorm:"type:varchar(1000);default:null"`
	Name                     *string    `json:"name,omitempty" gorm:"type:nvarchar(255);default:null"`
	Size                     *int       `json:"size,omitempty" gorm:"default:null"`
	Token                    *string    `json:"token,omitempty" gorm:"type:varchar(255);default:null"`
	OriginalName             *string    `json:"original_name,omitempty" gorm:"type:nvarchar(255);default:null"`
	Type                     *string    `json:"type,omitempty" gorm:"type:varchar(255);default:null"`
	IsReadOnly               bool       `json:"is_read_only" gorm:"type:bit;not null;default:0"`
	IsDownload               bool       `json:"is_download" gorm:"type:bit;not null;default:0"`
	CreatedBYID              uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID              *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt                time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                *time.Time `json:"updated_at,omitempty" gorm:"type:datetime;default:null"`

	// Relationships
	IncidentRegister       *IncidentRegister       `gorm:"foreignKey:IncidentRegisterID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	IncidentEvidenceSource *IncidentEvidenceSource `gorm:"foreignKey:IncidentEvidenceSourceID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	IncidentEvidenceType   *IncidentEvidenceType   `gorm:"foreignKey:IncidentEvidenceTypeID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser          *User                   `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedByUser          *User                   `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (IncidentEvidence) TableName() string {
	return "incident_evidences"
}
