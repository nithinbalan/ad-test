package dtos

import (
	"time"
)

type AuditPlanChecklistListDTO struct {
	ID                          uint       `json:"id"`
	AuditPlanID                 uint       `json:"audit_plan_id"`
	AuditableItemID             uint       `json:"auditable_item_id"`
	AuditableItemTitleEn        string     `json:"auditable_item_title_en"`
	AuditableItemTitleAr        string     `json:"auditable_item_title_ar"`
	ProcessID                   *uint      `json:"process_id"`
	ControlID                   *uint      `json:"control_id"`
	DocumentID                  *uint      `json:"document_id"`
	IncidentRegisterID          *uint      `json:"incident_register_id"`
	ComplianceRequirementID     *uint      `json:"compliance_requirement_id"`
	RiskRegisterID              *uint      `json:"risk_register_id"`
	UserID                      *uint      `json:"user_id"`
	Question                    string     `json:"question"`
	CreatedAt                   time.Time  `json:"created_at"`
	UpdatedAt                   *time.Time `json:"updated_at"`
	CreatedUserID               uint       `json:"created_user_id"`
	CreatedUserFirstNameEn      string     `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string     `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string     `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       string     `json:"created_user_last_name_ar"`
	CreatedUserEmail            string     `json:"created_user_email"`
	CreatedUserProfileImagePath string     `json:"created_user_profile_image_path"`
	UpdatedUserFirstNameEn      string     `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr      string     `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn       string     `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr       string     `json:"updated_user_last_name_ar"`
	UpdatedUserEmail            string     `json:"updated_user_email"`
	UpdatedUserProfileImagePath string     `json:"updated_user_profile_image_path"`
}

type AuditPlanChecklistFormDTO struct {
	AuditableItemId uint   `json:"auditable_item_id" binding:"required"`
	Question        string `json:"question" binding:"required"`
}

type AuditPlanChecklistListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Filters     map[string][]string
	QuickFilter []string
}
