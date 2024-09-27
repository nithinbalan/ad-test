package dtos

import (
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type AuditDraftReportCommentListDTO struct {
	ID                          uint                                  `json:"id"`
	Comment                     string                                `json:"comment" binding:"required"`
	Type                        string                                `json:"type"`
	IsResolved                  bool                                  `json:"is_resolved"`
	AuditDraftReportID          uint                                  `json:"audit_draft_report_id"`
	AuditDraftFindingID         *uint                                 `json:"audit_draft_finding_id"`
	CreatedAt                   string                                `json:"created_at"`
	UpdatedAt                   *string                               `json:"updated_at"`
	CreatedBYID                 uint                                  `json:"created_by_id"`
	CreatedUserFirstNameEn      string                                `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string                                `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string                                `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       string                                `json:"created_user_last_name_ar"`
	CreatedUserInitial          string                                `json:"created_user_initial"`
	CreatedUserColor            string                                `json:"created_user_color"`
	CreatedUserEmail            string                                `json:"created_user_email"`
	CreatedUserProfileImagePath string                                `json:"created_user_profile_image_path"`
	UpdatedBYID                 *uint                                 `json:"updated_by_id"`
	UpdatedUserFirstNameEn      *string                               `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr      *string                               `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn       *string                               `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr       *string                               `json:"updated_user_last_name_ar"`
	UpdatedUserInitial          *string                               `json:"updated_user_initial"`
	UpdatedUserColor            *string                               `json:"updated_user_color"`
	UpdatedUserEmail            *string                               `json:"updated_user_email"`
	UpdatedUserProfileImagePath *string                               `json:"updated_user_profile_image_path"`
	Documents                   []*AuditDraftReportCommentDocumentDTO `json:"documents" gorm:"-"`
}

type AuditDraftReportCommentListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Filters     map[string][]string
	QuickFilter []string
}

type AuditDraftReportCommentFormDTO struct {
	AuditDraftFindingID *uint                               `json:"audit_draft_finding_id" binding:"omitempty"`
	Comment             string                              `json:"comment" binding:"required,max=1000"`
	Type                string                              `json:"type" binding:"required,oneof=executive-summary finding conclusion"`
	Documents           []*AuditDraftFindingDocumentFormDTO `json:"documents"`
}
type AuditDraftReportCommentStatusUpdateFormDTO struct {
	IsResolved *bool `json:"is_resolved" binding:"required"`
}
type AuditDraftReportCommentDocumentFormDTO struct {
	ID            *uint  `json:"id,omitempty"`
	URL           string `json:"url"`
	FileName      string `json:"file_name"`
	OriginalName  string `json:"original_name"`
	Type          string `json:"type"`
	Size          int    `json:"size"`
	IsNewFile     *bool  `json:"is_new_file,omitempty"`
	IsDeletedFile *bool  `json:"is_deleted_file,omitempty"`
}

type AuditDraftReportCommentDetailDTO struct {
	ID                  uint                                  `json:"id"`
	Comment             string                                `json:"comment"`
	IsResolved          bool                                  `json:"is_resolved"`
	AuditDraftReportID  uint                                  `json:"audit_draft_report_id"`
	AuditDraftFindingID *uint                                 `json:"audit_draft_finding_id"`
	Documents           []*AuditDraftReportCommentDocumentDTO `json:"documents"`
	CreatedAt           time.Time                             `json:"created_at"`
	UpdatedAt           *time.Time                            `json:"updated_at"`
	CreatedByUser       organizationDTOs.UserMinimalDTO       `json:"created_by"`
	UpdatedByUser       *organizationDTOs.UserMinimalDTO      `json:"updated_by"`
}

type AuditDraftReportCommentDocumentDTO struct {
	ID           *uint  `json:"id"`
	URL          string `json:"url"`
	Name         string `json:"name"`
	Token        string `json:"token"`
	OriginalName string `json:"original_name"`
	Type         string `json:"type"`
	Size         int    `json:"size"`
}
