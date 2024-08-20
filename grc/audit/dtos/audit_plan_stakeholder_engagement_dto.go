package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type AuditPlanStakeholderEngagementListDTO struct {
	ID                              uint    `json:"id"`
	ReferenceNumber                 string  `json:"reference_number"`
	AuditPlanID                     uint    `json:"audit_plan_id"`
	AuditEngagementTypeID           uint    `json:"engagement_type_id"`
	AuditEngagementTypeTitleEn      string  `json:"audit_engagement_type_title_en"`
	AuditEngagementTypeTitleAr      string  `json:"audit_engagement_type_title_ar"`
	Title                           string  `json:"title"`
	Objectives                      *string `json:"objectives"`
	ResponsibleUserID               *uint   `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn      string  `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr      string  `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn       string  `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr       string  `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          string  `json:"responsible_user_initial"`
	ResponsibleUserColor            string  `json:"responsible_user_color"`
	ResponsibleUserEmail            string  `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath string  `json:"responsible_user_profile_image_path"`
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

type AuditPlanStakeholderEngagementFormDTO struct {
	AuditEngagementTypeID uint    `json:"engagement_type_id" binding:"required"`
	Title                 string  `json:"title" binding:"required,max=255"`
	Objectives            *string `json:"objectives"`
	ResponsibleUserID     *uint   `json:"responsible_user_id"`
}

type AuditPlanStakeholderEngagementDetailDTO struct {
	ID                  uint                             `json:"id"`
	ReferenceNumber     string                           `json:"reference_number"`
	AuditPlan           AuditPlanDetailDTO               `json:"audit_plan"`
	AuditEngagementType AuditEngagementTypeDTO           `json:"audit_engagement_type"`
	Title               string                           `json:"title"`
	Objectives          *string                          `json:"objectives"`
	ResponsibleUser     *organizationDTOs.UserMinimalDTO `json:"responsible_user"`
	AuditScheduleStatus AuditScheduleStatusDTO           `json:"audit_schedule_status"`
	CreatedAt           time.Time                        `json:"created_at"`
	UpdatedAt           *time.Time                       `json:"updated_at"`
	CreatedByUser       organizationDTOs.UserMinimalDTO  `json:"created_by_id"`
	UpdatedByUser       *organizationDTOs.UserMinimalDTO `json:"created_by_id"`
}

type AuditPlanStakeholderEngagementListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
}

type AuditPlanStakeholderEngagementListResponseDTO struct {
	Items      []*AuditPlanStakeholderEngagementListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO             `json:"pagination"`
}
