package dtos

import (
	"time"

	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
)

type ProjectSecurityGateChecklistCommentDTO struct {
	ID            uint                                              `json:"id"`
	Comment       string                                            `json:"comment"`
	CreatedByUser *organizationDTOs.UserMinimalDTO                  `json:"created_by"`
	CreatedAt     time.Time                                         `json:"created_at"`
	Documents     []*ProjectSecurityGateChecklistCommentDocumentDTO `json:"documents"`
}

type ProjectSecurityGateChecklistCommentDocumentDTO struct {
	ID           uint   `json:"id"`
	URL          string `json:"url"`
	Name         string `json:"name"`
	Size         int    `json:"size"`
	Token        string `json:"token"`
	OriginalName string `json:"original_name"`
	Type         string `json:"type"`
}

type ProjectSecurityGateChecklistCommentFormDTO struct {
	Comment   string                        `json:"comment" binding:"required,max=1000"`
	Documents []*commonDTOs.DocumentFormDTO `json:"documents" binding:"omitempty"`
}
