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
	ProcessTitleEn              *string    `json:"process_title_en"`
	ProcessTitleAr              *string    `json:"process_title_ar"`
	ControlTitleEn              *string    `json:"control_title_en"`
	ControlTitleAr              *string    `json:"control_title_ar"`
	DocumentTitleEn             *string    `json:"document_title_en"`
	DocumentTitleAr             *string    `json:"document_title_ar"`
	UserFirstNameEn             *string    `json:"user_first_name_en"`
	UserFirstNameAr             *string    `json:"user_first_name_ar"`
	UserLastNameEn              *string    `json:"user_last_name_en"`
	UserLastNameAr              *string    `json:"user_last_name_ar"`
	UserInitial                 *string    `json:"user_initial"`
	UserColor                   *string    `json:"user_color"`
	UserEmail                   *string    `json:"user_email"`
	UserProfileImagePath        *string    `json:"user_profile_image_path"`
	IncidentRegisterTitle       *string    `json:"incident_register_title"`
	ComplianceRequirementTitle  *string    `json:"compliance_requirement_title"`
	RiskRegisterTitle           *string    `json:"risk_register_title"`
	CreatedAt                   time.Time  `json:"created_at"`
	UpdatedAt                   *time.Time `json:"updated_at"`
	CreatedUserID               uint       `json:"created_user_id"`
	CreatedUserFirstNameEn      string     `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string     `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string     `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       string     `json:"created_user_last_name_ar"`
	CreatedUserInitial          string     `json:"created_user_initial"`
	CreatedUserColor            string     `json:"created_user_color"`
	CreatedUserEmail            string     `json:"created_user_email"`
	CreatedUserProfileImagePath string     `json:"created_user_profile_image_path"`
	UpdatedUserID               *uint      `json:"updated_user_id"`
	UpdatedUserFirstNameEn      string     `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr      string     `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn       string     `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr       string     `json:"updated_user_last_name_ar"`
	UpdatedUserInitial          string     `json:"updated_user_initial"`
	UpdatedUserColor            string     `json:"updated_user_color"`
	UpdatedUserEmail            string     `json:"updated_user_email"`
	UpdatedUserProfileImagePath string     `json:"updated_user_profile_image_path"`
}

type AuditPlanChecklistDetailDTO struct {
	ID              uint                   `json:"id"`
	AuditPlanID     uint                   `json:"audit_plan_id"`
	AuditableItemID uint                   `json:"auditable_item_id"`
	AuditableItem   AuditableItemDetailDTO `json:"auditable_item"`
	Question        string                 `json:"question"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       *time.Time             `json:"updated_at"`
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
