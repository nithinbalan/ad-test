package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type IncidentEvidenceListDTO struct {
	ID                        uint                       `json:"id"`
	IncidentRegisterID        uint                       `json:"incident_register_id"`
	IncidentEvidenceSourceID  uint                       `json:"incident_evidence_source_id"`
	IncidentEvidenceSourceEn  string                     `json:"incident_evidence_source_title_en"`
	IncidentEvidenceSourceAr  string                     `json:"incident_evidence_source_title_ar"`
	IncidentEvidenceTypeID    uint                       `json:"incident_evidence_type_id"`
	IncidentEvidenceTypeEn    string                     `json:"incident_evidence_type_title_en"`
	IncidentEvidenceTypeAr    string                     `json:"incident_evidence_type_title_ar"`
	Documents                 []*IncidentEvidenceListDTO `json:"documents" gorm:"-"`
	CreatedByID               uint                       `json:"created_by_id"`
	CreatedByFirstNameEn      string                     `json:"created_by_first_name_en"`
	CreatedByFirstNameAr      string                     `json:"created_by_first_name_ar"`
	CreatedByLastNameEn       string                     `json:"created_by_last_name_en"`
	CreatedByLastNameAr       string                     `json:"created_by_last_name_ar"`
	CreatedByInitial          string                     `json:"created_by_initial"`
	CreatedByColor            string                     `json:"created_by_color"`
	CreatedByEmail            string                     `json:"created_by_email"`
	CreatedByProfileImagePath *string                    `json:"created_by_profile_image_path"`
	CreatedAt                 time.Time                  `json:"created_at"`
	UpdatedByID               *uint                      `json:"updated_by_id"`
	UpdatedByFirstNameEn      *string                    `json:"updated_by_first_name_en"`
	UpdatedByFirstNameAr      *string                    `json:"updated_by_first_name_ar"`
	UpdatedByLastNameEn       *string                    `json:"updated_by_last_name_en"`
	UpdatedByLastNameAr       *string                    `json:"updated_by_last_name_ar"`
	UpdatedByInitial          *string                    `json:"updated_by_initial"`
	UpdatedByColor            *string                    `json:"updated_by_color"`
	UpdatedByEmail            *string                    `json:"updated_by_email"`
	UpdatedByProfileImagePath *string                    `json:"updated_by_profile_image_path"`
	UpdatedAt                 *time.Time                 `json:"updated_at"`
	URL                       *string                    `json:"url"`
	Name                      string                     `json:"name"`
	OriginalName              string                     `json:"original_name"`
	Token                     string                     `json:"token"`
	Type                      string                     `json:"type"`
	Size                      int                        `json:"size"`
	IsReadOnly                bool                       `json:"is_read_only"`
	IsDownload                bool                       `json:"is_download"`
}

type IncidentEvidenceFormDTO struct {
	IncidentEvidenceSourceID uint                          `json:"incident_evidence_source_id" binding:"required"`
	IncidentEvidenceTypeID   uint                          `json:"incident_evidence_type_id" binding:"required"`
	IsReadOnly               bool                          `json:"is_read_only"`
	IsDownload               bool                          `json:"is_download"`
	Documents                []*commonDTOs.DocumentFormDTO `json:"documents" binding:"required"` // At least one document expected
}

type IncidentEvidenceDownloadHistoryListDTO struct {
	ID                           uint      `json:"id"`
	IncidentEvidenceID           uint      `json:"incident_evidence_id"`
	DownloadedAt                 time.Time `json:"downloaded_at"`
	DownloadedByID               uint      `json:"downloaded_by_id"`
	DownloadedByFirstNameEn      string    `json:"downloaded_by_first_name_en"`
	DownloadedByFirstNameAr      string    `json:"downloaded_by_first_name_ar"`
	DownloadedByLastNameEn       string    `json:"downloaded_by_last_name_en"`
	DownloadedByLastNameAr       string    `json:"downloaded_by_last_name_ar"`
	DownloadedByInitial          string    `json:"downloaded_by_initial"`
	DownloadedByColor            string    `json:"downloaded_by_color"`
	DownloadedByEmail            string    `json:"downloaded_by_email"`
	DownloadedByProfileImagePath *string   `json:"downloaded_by_profile_image_path"`
}

type IncidentEvidenceDetailDTO struct {
	ID                     uint                             `json:"id"`
	IncidentRegisterID     uint                             `json:"incident_register_id"`
	IncidentEvidenceSource IncidentEvidenceSourceDTO        `json:"incident_evidence_source"`
	IncidentEvidenceType   IncidentEvidenceTypeDTO          `json:"incident_evidence_type"`
	URL                    string                           `json:"url"`
	Name                   string                           `json:"name"`
	Size                   int                              `json:"size"`
	Token                  string                           `json:"token"`
	OriginalName           string                           `json:"original_name"`
	Type                   string                           `json:"type"`
	IsReadOnly             bool                             `json:"is_read_only"`
	IsDownload             bool                             `json:"is_download"`
	CreatedAt              time.Time                        `json:"created_at"`
	UpdatedAt              *time.Time                       `json:"updated_at"`
	CreatedByUser          organizationDTOs.UserMinimalDTO  `json:"created_by"`
	UpdatedByUser          *organizationDTOs.UserMinimalDTO `json:"updated_by"`
}
