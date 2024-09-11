package dtos

import (
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type AuditDraftFindingConversationListDTO struct {
	ID                          uint   `json:"id"`
	AuditDraftFindingID         uint   `json:"audit_draft_finding_id"`
	Description                 string `json:"description"`
	CreatedAt                   string `json:"created_at"`
	CreatedUserID               uint   `json:"created_user_id"`
	CreatedUserFirstNameEn      string `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       string `json:"created_user_last_name_ar"`
	CreatedUserInitial          string `json:"created_user_initial"`
	CreatedUserColor            string `json:"created_user_color"`
	CreatedUserEmail            string `json:"created_user_email"`
	CreatedUserProfileImagePath string `json:"created_user_profile_image_path"`
}

type AuditDraftFindingConversationFormDTO struct {
	Description string `json:"description" binding:"required,max=500"`
	Type        string `json:"type" binding:"required,max=10"`
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
	ID                int                             `json:"id"`
	Type              string                          `json:"type"`
	Description       string                          `json:"description"`
	AuditDraftFinding AuditDraftFindingDetailDTO      `json:"audit_draft_finding"`
	CreatedAt         time.Time                       `json:"created_at"`
	UpdatedAt         *time.Time                      `json:"updated_at"`
	CreatedByUser     organizationDTOs.UserMinimalDTO `json:"created_by_user"`
	UpdatedByUser     organizationDTOs.UserMinimalDTO `json:"updated_by_user"`
}
