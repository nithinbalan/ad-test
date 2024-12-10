package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type AuditChecklistListDTO struct {
	ID                                   uint       `json:"id"`
	AuditID                              uint       `json:"audit_id"`
	AuditPlanChecklistID                 *uint      `json:"audit_plan_checklist_id"`
	AuditableItemID                      *uint      `json:"auditable_item_id"`
	AuditableItemType                    string     `json:"auditable_item_type"`
	AuditableItemTitleEn                 string     `json:"auditable_item_title_en"`
	AuditableItemTitleAr                 string     `json:"auditable_item_title_ar"`
	AuditableItemComment                 string     `json:"auditable_item_comment"`
	AuditableItemProcessID               *uint      `json:"auditable_item_process_id"`
	AuditableItemControlID               *uint      `json:"auditable_item_control_id"`
	AuditableItemDocumentID              *uint      `json:"auditable_item_document_id"`
	AuditableItemIncidentRegisterID      *uint      `json:"auditable_item_incident_register_id"`
	AuditableItemResponsibleUnitID       uint       `json:"auditable_item_responsible_unit_id"`
	AuditableItemComplianceRequirementID *uint      `json:"auditable_item_compliance_requirement_id"`
	AuditableItemRiskRegisterID          *uint      `json:"auditable_item_risk_register_id"`
	AuditableItemUserID                  *uint      `json:"auditable_item_user_id"`
	Question                             string     `json:"question"`
	ConformanceID                        *uint      `json:"conformance_id"`
	ConformanceTitleEn                   string     `json:"conformance_title_en"`
	ConformanceTitleAr                   string     `json:"conformance_title_ar"`
	ConformanceType                      string     `json:"conformance_type"`
	ConformanceColorCode                 string     `json:"conformance_color_code"`
	Comment                              *string    `json:"comment"`
	AnsweredByID                         *uint      `json:"answered_by_id"`
	AnsweredAt                           *time.Time `json:"answered_at"`
	CreatedByID                          uint       `json:"created_by_id"`
	CreatedAt                            time.Time  `json:"created_at"`
	UpdatedByID                          *uint      `json:"updated_by_id"`
	UpdatedAt                            *time.Time `json:"updated_at"`
	StatusID                             uint       `json:"status_id"`
	StatusTitleEn                        string     `json:"status_title_en"`
	StatusTitleAr                        string     `json:"status_title_ar"`
	StatusType                           string     `json:"status_type"`
	StatusColorCode                      string     `json:"status_color_code"`
	ProcessTitleEn                       string     `json:"process_title_en"`
	ProcessTitleAr                       string     `json:"process_title_ar"`
	ControlTitleEn                       string     `json:"control_title_en"`
	ControlTitleAr                       string     `json:"control_title_ar"`
	DocumentTitleEn                      string     `json:"document_title_en"`
	DocumentTitleAr                      string     `json:"document_title_ar"`
	UserFirstNameEn                      string     `json:"user_first_name_en"`
	UserFirstNameAr                      string     `json:"user_first_name_ar"`
	UserLastNameEn                       string     `json:"user_last_name_en"`
	UserLastNameAr                       string     `json:"user_last_name_ar"`
	UserInitial                          string     `json:"user_initial"`
	UserColor                            string     `json:"user_color"`
	UserEmail                            string     `json:"user_email"`
	UserProfileImagePath                 string     `json:"user_profile_image_path"`
	IncidentRegisterTitle                string     `json:"incident_register_title"`
	ComplianceRequirementTitle           string     `json:"compliance_requirement_title"`
	RiskRegisterTitle                    string     `json:"risk_register_title"`
	AnsweredUserFirstNameEn              string     `json:"answered_user_first_name_en"`
	AnsweredUserFirstNameAr              string     `json:"answered_user_first_name_ar"`
	AnsweredUserLastNameEn               string     `json:"answered_user_last_name_en"`
	AnsweredUserLastNameAr               string     `json:"answered_user_last_name_ar"`
	AnsweredUserInitial                  string     `json:"answered_user_initial"`
	AnsweredUserColor                    string     `json:"answered_user_color"`
	AnsweredUserEmail                    string     `json:"answered_user_email"`
	AnsweredUserProfileImagePath         string     `json:"answered_user_profile_image_path"`
	CreatedUserFirstNameEn               string     `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr               string     `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn                string     `json:"created_user_last_name_en"`
	CreatedUserLastNameAr                string     `json:"created_user_last_name_ar"`
	CreatedUserInitial                   string     `json:"created_user_initial"`
	CreatedUserColor                     string     `json:"created_user_color"`
	CreatedUserEmail                     string     `json:"created_user_email"`
	CreatedUserProfileImagePath          string     `json:"created_user_profile_image_path"`
	UpdatedUserID                        *uint      `json:"updated_user_id"`
	UpdatedUserFirstNameEn               *string    `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr               *string    `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn                *string    `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr                *string    `json:"updated_user_last_name_ar"`
	UpdatedUserInitial                   *string    `json:"updated_user_initial"`
	UpdatedUserColor                     *string    `json:"updated_user_color"`
	UpdatedUserEmail                     *string    `json:"updated_user_email"`
	UpdatedUserProfileImagePath          *string    `json:"updated_user_profile_image_path"`
}

type AuditChecklistAnswerFormDTO struct {
	ConformanceID uint                          `json:"conformance_id" binding:"required"`
	Comment       *string                       `json:"comment"`
	Documents     []*commonDTOs.DocumentFormDTO `json:"documents" binding:"omitempty"`
}

type AuditChecklistFormDTO struct {
	AuditableItemID *uint  `json:"auditable_item_id"`
	Question        string `json:"question" binding:"required"`
}

type AuditChecklistDetailDTO struct {
	ID                   uint                             `json:"id"`
	AuditID              uint                             `json:"audit_id"`
	AuditPlanChecklistID *uint                            `json:"audit_plan_checklist_id"`
	Question             string                           `json:"question"`
	ConformanceID        *uint                            `json:"conformance_id"`
	Comment              *string                          `json:"comment"`
	AnsweredByID         *uint                            `json:"answered_by_id"`
	AnsweredAt           *time.Time                       `json:"answered_at"`
	CreatedByID          uint                             `json:"created_by_id"`
	UpdatedByID          *uint                            `json:"updated_by_id"`
	StatusID             uint                             `json:"status_id"`
	AuditPlanChecklist   *AuditPlanChecklistDetailDTO     `json:"audit_plan_checklist"`
	AuditDraftFindings   []*AuditDraftFindingDetailDTO    `json:"audit_draft_findings"`
	AuditableItem        *AuditableItemDetailDTO          `json:"auditable_item"`
	Conformance          *masterDTOs.ConformanceDTO       `json:"conformance"`
	AnsweredByUser       *organizationDTOs.UserMinimalDTO `json:"answered_by"`
	CreatedAt            time.Time                        `json:"created_at"`
	UpdatedAt            *time.Time                       `json:"updated_at"`
	CreatedByUser        organizationDTOs.UserMinimalDTO  `json:"created_by"`
	UpdatedByUser        *organizationDTOs.UserMinimalDTO `json:"updated_by"`
	Status               *masterDTOs.StatusDTO            `json:"status"`
	Documents            []*AuditChecklistDocumentDTO     `json:"documents"`
}

type AuditChecklistListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Filters     map[string][]string
	QuickFilter []string
}

type AuditChecklistDocumentDTO struct {
	ID           *uint   `json:"id"`
	URL          *string `json:"url"`
	Name         string  `json:"name"`
	Token        string  `json:"token"`
	OriginalName string  `json:"original_name"`
	Type         string  `json:"type"`
	Size         int     `json:"size"`
}
