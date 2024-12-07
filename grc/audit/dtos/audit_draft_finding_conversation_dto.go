package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type AuditDraftFindingConversationListDTO struct {
	ID                          uint                                        `json:"id"`
	AuditDraftFindingID         uint                                        `json:"audit_draft_finding_id"`
	Comment                     string                                      `json:"comment"`
	IsAccepted                  bool                                        `json:"is_accepted"`
	IsDismissed                 bool                                        `json:"is_dismissed"`
	IsReopened                  bool                                        `json:"is_reopened"`
	CreatedAt                   string                                      `json:"created_at"`
	CreatedUserID               uint                                        `json:"created_user_id"`
	CreatedUserFirstNameEn      string                                      `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string                                      `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string                                      `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       string                                      `json:"created_user_last_name_ar"`
	CreatedUserInitial          string                                      `json:"created_user_initial"`
	CreatedUserColor            string                                      `json:"created_user_color"`
	CreatedUserEmail            string                                      `json:"created_user_email"`
	CreatedUserProfileImagePath string                                      `json:"created_user_profile_image_path"`
	Documents                   []*AuditDraftFindingConversationDocumentDTO `json:"documents" gorm:"-"`
}

type AuditDraftFindingConversationFormDTO struct {
	Comment     string                        `json:"comment" binding:"required,max=1000"`
	Type        string                        `json:"type" binding:"required,max=10"`
	IsAccepted  *bool                         `json:"is_accepted" binding:"omitempty"`
	IsDismissed *bool                         `json:"is_dismissed" binding:"omitempty"`
	IsReopened  *bool                         `json:"is_reopened" binding:"omitempty"`
	Documents   []*commonDTOs.DocumentFormDTO `json:"documents" binding:"omitempty"`
}

type AuditDraftFindingConversationListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
	CreatedAt   string
}

type AuditDraftFindingConversationDetailDTO struct {
	ID                int                                         `json:"id"`
	Type              string                                      `json:"type"`
	Comment           string                                      `json:"comment"`
	AuditDraftFinding AuditDraftFindingDetailDTO                  `json:"audit_draft_finding"`
	IsAccepted        bool                                        `json:"is_accepted"`
	IsDismissed       bool                                        `json:"is_dismissed"`
	IsReopened        bool                                        `json:"is_reopened"`
	CreatedAt         time.Time                                   `json:"created_at"`
	UpdatedAt         *time.Time                                  `json:"updated_at"`
	CreatedByUser     organizationDTOs.UserMinimalDTO             `json:"created_by_user"`
	UpdatedByUser     organizationDTOs.UserMinimalDTO             `json:"updated_by_user"`
	Documents         []*AuditDraftFindingConversationDocumentDTO `json:"documents"`
}

type AuditDraftFindingConversationDocumentDTO struct {
	ID           uint   `json:"id"`
	URL          string `json:"url"`
	Name         string `json:"name"`
	Size         int    `json:"size"`
	Token        string `json:"token"`
	OriginalName string `json:"original_name"`
	Type         string `json:"type"`
}
