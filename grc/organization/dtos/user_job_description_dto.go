package dtos

import "time"

type UserJobDescriptionListDTO struct {
	ID                       uint   `json:"id"`
	TitleEn                  string `json:"title_en"`
	TitleAr                  string `json:"title_ar"`
	DescriptionEn            string `json:"description_en"`
	DescriptionAr            string `json:"description_ar"`
	RolesAndResponsibilities string `json:"roles_and_responsibilities"`
	Status                   string `json:"status"`
	CreatedAt                string `json:"created_at"`
}

type AssignJobDescriptionsDTO struct {
	JobDescriptionIDs []uint `json:"job_description_ids" binding:"required,min=1,dive,required"`
}

type EvaluateUserJDCompetencyDTO struct {
	ActualLevelID uint                       `json:"actual_level_id" binding:"required"`
	Comments      *string                    `json:"comments"`
	Documents     []UserCertificateUploadDTO `json:"documents"` // One for each required certification
}

type UserCertificateUploadDTO struct {
	ID              uint   `json:"id"`
	CertificationID uint   `json:"certification_id" binding:"required"`
	FileName        string `json:"file_name" binding:"required"`
	OriginalName    string `json:"original_name"`
	Size            int    `json:"size"`
	Type            string `json:"type"`
	IsNewFile       *bool  `json:"is_new_file"`
	IsDeletedFile   *bool  `json:"is_deleted_file"`
}

type UserJobDescriptionDetailDTO struct {
	ID                     uint                                    `json:"id"`
	UserID                 uint                                    `json:"user_id"`
	JobDescriptionID       uint                                    `json:"job_description_id"`
	JobDescription         *JobDescriptionDTO                      `json:"job_description"`
	Competencies           []UserJobDescriptionCompetencyDetailDTO `json:"competencies"`
	CompetencyTitleEn      *string                                 `json:"competency_title_en"`
	CompetencyTitleAr      *string                                 `json:"competency_title_ar"`
	CompetencyGroupTitleEn *string                                 `json:"competency_group_title_en"`
	CompetencyGroupTitleAr *string                                 `json:"competency_group_title_ar"`
	CreatedAt              time.Time                               `json:"created_at"`
}

type UserJobDescriptionCompetencyDetailDTO struct {
	ID                     uint                                     `json:"id"`
	CompetencyID           uint                                     `json:"competency_id"`
	Competency             *CompetencyDTO                           `json:"competency"`
	ActualLevelID          uint                                     `json:"actual_level_id"`
	ActualLevel            *CompetencyLevelDTO                      `json:"actual_level"`
	Comments               *string                                  `json:"comments"`
	Certificates           []UserJobDescriptionCertificateDetailDTO `json:"certificates"`
	CompetencyTitleEn      *string                                  `json:"competency_title_en"`
	CompetencyTitleAr      *string                                  `json:"competency_title_ar"`
	CompetencyGroupTitleEn *string                                  `json:"competency_group_title_en"`
	CompetencyGroupTitleAr *string                                  `json:"competency_group_title_ar"`
}

type UserJobDescriptionCertificateDetailDTO struct {
	ID                                      uint    `json:"id"`
	JobDescriptionCompetencyCertificationID uint    `json:"certification_id"`
	CertificateTitle                        string  `json:"certificate_title"`
	Token                                   *string `json:"token"`
	OriginalName                            *string `json:"original_name"`
	URL                                     *string `json:"url"`
	Type                                    *string `json:"type"`
	Size                                    *string `json:"size"`
}
