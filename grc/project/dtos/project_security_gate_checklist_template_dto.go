package dtos

import (
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type ProjectSecurityGateChecklistTemplateDTO struct {
	ID      uint                         `json:"id"`
	StageID uint                         `json:"stage_id"`
	Stage   *ProjectSecurityGateStageDTO `json:"stage"`
	TitleEn string                       `json:"title_en"`
	TitleAr string                       `json:"title_ar"`
}

type ProjectSecurityGateChecklistTemplateDetailDTO struct {
	ID              uint                                                    `json:"id"`
	StageID         uint                                                    `json:"stage_id"`
	Stage           *ProjectSecurityGateStageDTO                            `json:"stage"`
	TitleEn         string                                                  `json:"title_en"`
	TitleAr         string                                                  `json:"title_ar"`
	LatestVersionID *uint                                                   `json:"latest_version_id"`
	Versions        []*ProjectSecurityGateChecklistTemplateVersionDetailDTO `json:"versions"`
}

type ProjectSecurityGateChecklistTemplateVersionDetailDTO struct {
	ID              uint                                                  `json:"id"`
	TemplateID      uint                                                  `json:"template_id"`
	Version         int                                                   `json:"version"`
	Status          string                                                `json:"status"`
	CreatedByUser   *organizationDTOs.UserMinimalDTO                      `json:"created_by_user"`
	CreatedAt       time.Time                                             `json:"created_at"`
	PublishedByUser *organizationDTOs.UserMinimalDTO                      `json:"published_by_user"`
	PublishedAt     *time.Time                                            `json:"published_at"`
	Requirements    []*ProjectSecurityGateChecklistTemplateRequirementDTO `json:"requirements"`
}

type ProjectSecurityGateChecklistTemplateRequirementDTO struct {
	ID                uint   `json:"id"`
	TemplateVersionID uint   `json:"template_version_id"`
	TitleEn           string `json:"title_en"`
	TitleAr           string `json:"title_ar"`
	IsCritical        bool   `json:"is_critical"`
}

type ProjectSecurityGateChecklistTemplateRequirementFormDTO struct {
	TitleEn    string `json:"title_en" binding:"required,max=500"`
	TitleAr    string `json:"title_ar" binding:"required,max=500"`
	IsCritical bool   `json:"is_critical"`
}
