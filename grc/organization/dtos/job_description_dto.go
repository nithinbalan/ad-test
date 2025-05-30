package dtos

import "time"

type JobDescriptionDTO struct {
	ID                       int    `json:"id"`
	TitleEn                  string `json:"title_en"`
	TitleAr                  string `json:"title_ar"`
	DescriptionEn            string `json:"description_en"`
	DescriptionAr            string `json:"description_ar"`
	RolesAndResponsibilities string `json:"roles_and_responsibilities"`
	Status                   string `json:"status"`
}

type JobDescriptionDetailDTO struct {
	ID                       uint                                `json:"id"`
	TitleEn                  string                              `json:"title_en"`
	TitleAr                  string                              `json:"title_ar"`
	DescriptionEn            string                              `json:"description_en"`
	DescriptionAr            string                              `json:"description_ar"`
	RolesAndResponsibilities string                              `json:"roles_and_responsibilities"`
	Status                   string                              `json:"status"`
	CreatedAt                time.Time                           `json:"created_at"`
	UpdatedAt                *time.Time                          `json:"updated_at,omitempty"`
	CreatedByUser            *UserMinimalDTO                     `json:"created_by_user"`
	UpdatedByUser            *UserMinimalDTO                     `json:"updated_by_user,omitempty"`
	Competencies             []JobDescriptionCompetencyDetailDTO `json:"competencies"`
}
