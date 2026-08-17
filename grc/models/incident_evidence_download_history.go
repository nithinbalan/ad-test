package models

import "time"

type IncidentEvidenceDownloadHistory struct {
	ID                 uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	IncidentEvidenceID uint      `json:"incident_evidence_id" gorm:"not null"`
	DownloadedAt       time.Time `json:"downloaded_at" gorm:"type:datetime;not null"`
	DownloadedBy       uint      `json:"downloaded_by" gorm:"not null"`

	// Relationships
	IncidentEvidence *IncidentEvidence `gorm:"foreignKey:IncidentEvidenceID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	DownloadedByUser *User             `gorm:"foreignKey:DownloadedBy;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (IncidentEvidenceDownloadHistory) TableName() string {
	return "incident_evidence_download_histories"
}
