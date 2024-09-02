package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

// AuditScheduleListDTO represents a list view of audit schedules.
type AuditScheduleListDTO struct {
	ID                              uint    `json:"id"`
	AuditID                         uint    `json:"audit_id"`
	AuditActivityTypeID             uint    `json:"audit_activity_type_id"`
	AuditActivityTypeTitleEn        string  `json:"audit_activity_type_title_en"`
	AuditActivityTypeTitleAr        string  `json:"audit_activity_type_title_ar"`
	ResponsibleUserID               *uint   `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn      string  `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr      string  `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn       string  `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr       string  `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          string  `json:"responsible_user_initial"`
	ResponsibleUserColor            string  `json:"responsible_user_color"`
	ResponsibleUserEmail            string  `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath string  `json:"responsible_user_profile_image_path"`
	Title                           string  `json:"title"`
	Description                     *string `json:"description"`
	StartDate                       *string `json:"start_date"`
	TargetDate                      *string `json:"target_date"`
	AuditScheduleStatusID           uint    `json:"audit_schedule_status_id"`
	AuditScheduleStatusTitleEn      string  `json:"audit_schedule_status_title_en"`
	AuditScheduleStatusTitleAr      string  `json:"audit_schedule_status_title_ar"`
	AuditScheduleStatusColorCode    string  `json:"audit_schedule_status_color_code"`
	CreatedUserID                   uint    `json:"created_user_id"`
	CreatedUserFirstNameEn          string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr          string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn           string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr           string  `json:"created_user_last_name_ar"`
	CreatedUserInitial              string  `json:"created_user_initial"`
	CreatedUserColor                string  `json:"created_user_color"`
	CreatedUserEmail                string  `json:"created_user_email"`
	CreatedUserProfileImagePath     string  `json:"created_user_profile_image_path"`
	UpdatedUserID                   *uint   `json:"updated_user_id"`
	UpdatedUserFirstNameEn          string  `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr          string  `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn           string  `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr           string  `json:"updated_user_last_name_ar"`
	UpdatedUserInitial              string  `json:"updated_user_initial"`
	UpdatedUserColor                string  `json:"updated_user_color"`
	UpdatedUserEmail                string  `json:"updated_user_email"`
	UpdatedUserProfileImagePath     string  `json:"updated_user_profile_image_path"`
	CreatedAt                       string  `json:"created_at"`
	UpdatedAt                       *string `json:"updated_at"`
}

// AuditScheduleFormDTO represents the form data for creating or updating an audit schedule.
type AuditScheduleFormDTO struct {
	AuditActivityTypeID uint    `json:"audit_activity_type_id" binding:"required"`
	ResponsibleUserID   *uint   `json:"responsible_user_id"`
	Title               string  `json:"title" binding:"required,max=255"`
	Description         *string `json:"description"`
	StartDate           *string `json:"start_date"`
	TargetDate          *string `json:"target_date"`
	AuditableItemIds    []*uint `json:"auditable_item_ids"`
}

// AuditScheduleDetailDTO represents the detailed view of an audit schedule.
type AuditScheduleDetailDTO struct {
	ID                  uint                             `json:"id"`
	Audit               AuditDetailDTO                   `json:"audit"`
	AuditPlan           AuditPlanDetailDTO               `json:"audit_plan"`
	AuditActivityType   AuditActivityTypeDTO             `json:"audit_activity_type"`
	ResponsibleUser     *organizationDTOs.UserMinimalDTO `json:"responsible_user"`
	Title               string                           `json:"title"`
	Description         *string                          `json:"description"`
	StartDateString     string                           `json:"start_date"`
	TargetDateString    string                           `json:"target_date"`
	AuditScheduleStatus AuditScheduleStatusDTO           `json:"audit_schedule_status"`
	AuditableItems      []*AuditScheduleAuditableItemDTO `json:"auditable_items"`
	CreatedAt           time.Time                        `json:"created_at"`
	UpdatedAt           *time.Time                       `json:"updated_at"`
	CreatedByUser       organizationDTOs.UserMinimalDTO  `json:"created_by_id"`
	UpdatedByUser       *organizationDTOs.UserMinimalDTO `json:"updated_by_id"`
}

// AuditScheduleListQueryParamsDTO represents the query parameters for filtering, sorting, and pagination.
type AuditScheduleListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
	StartDate   string
	TargetDate  string
}

// AuditScheduleListResponseDTO represents the paginated response format for listing audit schedules.
type AuditScheduleListResponseDTO struct {
	Items      []*AuditScheduleListDTO      `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

// AuditScheduleParsedDates represents parsed date fields for easier manipulation in services.
type AuditScheduleParsedDates struct {
	StartDate  *time.Time
	TargetDate *time.Time
}

// AuditScheduleAuditableItemDTO represents the auditable items associated with an audit schedule.
type AuditScheduleAuditableItemDTO struct {
	ID              uint                             `json:"id"`
	AuditScheduleID uint                             `json:"audit_schedule_id"`
	AuditableItemID uint                             `json:"auditable_item_id"`
	AuditableItem   AuditableItemDetailDTO           `json:"auditable_item"`
	CreatedAt       time.Time                        `json:"created_at"`
	CreatedByUser   organizationDTOs.UserMinimalDTO  `json:"created_by_id"`
	UpdatedByUser   *organizationDTOs.UserMinimalDTO `json:"updated_by_id"`
}
