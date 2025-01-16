package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
)

type AuditFinalReportAppendixFormDTO struct {
	Documents []*commonDTOs.DocumentFormDTO `json:"documents" binding:"required"`
}

type AuditFinalReportAppendixDetailDTO struct {
	ID                 uint                             `json:"id"`
	AuditFinalReportID uint                             `json:"audit_final_report_id"`
	URL                string                           `json:"url"`
	Name               string                           `json:"name"`
	Size               int                              `json:"size"`
	Token              string                           `json:"token"`
	OriginalName       string                           `json:"original_name"`
	Type               string                           `json:"type"`
	CreatedBYID        uint                             `json:"created_by_id"`
	UpdatedBYID        *uint                            `json:"updated_by_id"`
	CreatedByUser      organizationDTOs.UserMinimalDTO  `json:"created_by"`
	UpdatedByUser      *organizationDTOs.UserMinimalDTO `json:"updated_by"`
	CreatedAt          string                           `json:"created_at"`
	UpdatedAt          *string                          `json:"updated_at"`
}
