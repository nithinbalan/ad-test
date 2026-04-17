package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"

	"time"
)

type RiskActionPlanListDTO struct {
	ID                              uint    `json:"id"`
	RiskRegisterID                  uint    `json:"risk_register_id"`
	RiskReferenceNumber             string  `json:"risk_reference_number"`
	RiskTitle                       string  `json:"risk_title"`
	Title                           string  `json:"title"`
	Description                     string  `json:"description"`
	SeverityLevelTitleEn            string  `json:"severity_level_title_en"`
	SeverityLevelTitleAr            string  `json:"severity_level_title_ar"`
	SeverityLevelColorCode          string  `json:"severity_level_color_code"`
	ResponsibleUnitID               *uint   `json:"responsible_unit_id"`
	ResponsibleUnitTitleEn          string  `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr          string  `json:"responsible_unit_title_ar"`
	ResponsibleUserID               uint    `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn      string  `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr      string  `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn       string  `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr       string  `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          string  `json:"responsible_user_initial"`
	ResponsibleUserColor            string  `json:"responsible_user_color"`
	ResponsibleUserEmail            string  `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath string  `json:"responsible_user_profile_image_path"`
	ReviewUserFirstNameEn           string  `json:"review_user_first_name_en"`
	ReviewUserID                    *uint   `json:"review_user_id"`
	ReviewUserFirstNameAr           *string `json:"review_user_first_name_ar"`
	ReviewUserLastNameEn            *string `json:"review_user_last_name_en"`
	ReviewUserLastNameAr            *string `json:"review_user_last_name_ar"`
	ReviewUserInitial               *string `json:"review_user_initial"`
	ReviewUserColor                 *string `json:"review_user_color"`
	ReviewUserEmail                 *string `json:"review_user_email"`
	ReviewUserProfileImagePath      *string `json:"review_user_profile_image_path"`
	StartDate                       string  `json:"start_date"`
	TargetDate                      string  `json:"target_date"`
	Cost                            *int    `json:"cost"`
	Percentage                      uint    `json:"percentage"`
	Comments                        string  `json:"comments"`
	Notes                           string  `json:"notes"`
	ActionPlanStatusTitleEn         string  `json:"action_plan_status_title_en"`
	ActionPlanStatusTitleAr         string  `json:"action_plan_status_title_ar"`
	ActionPlanStatusColorCode       string  `json:"action_plan_status_color_code"`
	CreatedUserID                   uint    `json:"created_user_id"`
	CreatedUserFirstNameEn          string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr          string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn           string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr           string  `json:"created_user_last_name_ar"`
	CreatedUserInitial              string  `json:"created_user_initial"`
	CreatedUserColor                string  `json:"created_user_color"`
	CreatedUserEmail                string  `json:"created_user_email"`
	CreatedUserProfileImagePath     string  `json:"created_user_profile_image_path"`
	UpdatedUserID                   uint    `json:"updated_user_id"`
	UpdatedUserFirstNameEn          string  `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr          string  `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn           string  `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr           string  `json:"updated_user_last_name_ar"`
	UpdatedUserInitial              string  `json:"updated_user_initial"`
	UpdatedUserColor                string  `json:"updated_user_color"`
	UpdatedUserEmail                string  `json:"updated_user_email"`
	UpdatedUserProfileImagePath     string  `json:"updated_user_profile_image_path"`
}

type RiskActionPlanListQueryParamsDTO struct {
	Sort           string
	SortBy         string
	Keyword        string
	Status         string
	Page           int
	PageSize       int
	Filters        map[string][]string
	QuickFilter    []string
	TargetDateFrom string
	TargetDateTo   string
	Export         bool
	DateFrom       string
	DateTo         string
}

type RiskActionPlanListResponseDTO struct {
	Items      []*RiskActionPlanListDTO     `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type RiskActionPlanFormDTO struct {
	Title              string            `json:"title" binding:"required"`
	Description        *string           `json:"description" binding:"omitempty"`
	SeverityLevelID    *uint             `json:"severity_level_id" binding:"omitempty"`
	ResponsibleUnitID  *uint             `json:"responsible_unit_id" binding:"required"`
	ResponsibleUserID  uint              `json:"responsible_user_id" binding:"required"`
	StartDate          *string           `json:"start_date" binding:"omitempty,datetime=2006-01-02"`
	TargetDate         *string           `json:"target_date" binding:"omitempty,datetime=2006-01-02"`
	Cost               *int              `json:"cost" binding:"omitempty"`
	Percentage         *uint             `json:"percentage" binding:"omitempty"`
	Comments           *string           `json:"comments" binding:"omitempty"`
	Notes              *string           `json:"notes" binding:"omitempty"`
	ReviewUserID       *uint             `json:"review_user_id"`
	ActionPlanStatusID uint              `json:"action_plan_status_id"`
	Controls           []*ControlFormDTO `json:"controls"`
}

type ControlFormDTO struct {
	ID            *uint  `json:"id"`
	ControlTypeID uint   `json:"control_type_id" binding:"required"`
	OwnerID       uint   `json:"owner_id" binding:"required"`
	TitleEn       string `json:"title_en" binding:"required"`
	TitleAr       string `json:"title_ar" binding:"required"`
}

type RiskActionPlanParsedDates struct {
	StartDate  *time.Time
	TargetDate *time.Time
}

type RiskActionPlanDetailDTO struct {
	ID                      int                                   `json:"id"`
	Title                   string                                `json:"title"`
	Description             string                                `json:"description"`
	RiskRegister            RiskRegisterDetailDTO                 `json:"risk_register"`
	SeverityLevel           *masterDTOs.SeverityLevelDTO          `json:"severity_level"`
	ResponsibleUnit         *organizationDTOs.OrganizationUnitDTO `json:"responsible_unit"`
	ResponsibleUnitID       *uint                                 `json:"responsible_unit_id"`
	ResponsibleUser         organizationDTOs.UserMinimalDTO       `json:"responsible_user"`
	ResponsibleUserID       uint                                  `json:"responsible_user_id"`
	ReviewUser              *UserMinimalDTO                       `json:"review_user"`
	ReviewUserID            *uint                                 `json:"review_user_id"`
	ActionPlanStatus        masterDTOs.ActionPlanStatusDTO        `json:"action_plan_status"`
	Percentage              *uint                                 `json:"percentage"`
	Comments                *string                               `json:"comments"`
	Notes                   *string                               `json:"notes"`
	StartDateString         *string                               `json:"start_date"`
	TargetDateString        *string                               `json:"target_date"`
	Cost                    *int                                  `json:"cost"`
	DelegateeUser           *organizationDTOs.UserMinimalDTO      `json:"delegatee_user"`
	CreatedAt               time.Time                             `json:"created_at"`
	UpdatedAt               *time.Time                            `json:"updated_at"`
	CreatedByUser           organizationDTOs.UserMinimalDTO       `json:"created_by_user"`
	UpdatedByUser           *organizationDTOs.UserMinimalDTO      `json:"updated_by_user"`
	ActionPlanStatusUpdates []*RiskActionPlanStatusUpdateDTO      `json:"action_plan_status_updates"`
	Controls                []*RiskActionPlanControlDTO           `json:"controls"`
}

type RiskActionPlanStatusUpdateDTO struct {
	ID               uint                                     `json:"id"`
	ActionPlanStatus masterDTOs.ActionPlanStatusDTO           `json:"action_plan_status"`
	Comment          string                                   `json:"comment"`
	Percentage       *uint                                    `json:"percentage"`
	CreatedAt        time.Time                                `json:"created_at"`
	UpdatedAt        *time.Time                               `json:"updated_at"`
	CreatedByUser    organizationDTOs.UserMinimalDTO          `json:"created_by_user"`
	UpdatedByUser    *organizationDTOs.UserMinimalDTO         `json:"updated_by_user"`
	Documents        []*RiskActionPlanStatusUpdateDocumentDTO `json:"documents"`
}

type RiskActionPlanStatusUpdateDocumentDTO struct {
	ID           uint   `json:"id"`
	URL          string `json:"url"`
	Name         string `json:"name"`
	Size         int    `json:"size"`
	Token        string `json:"token"`
	OriginalName string `json:"original_name"`
	Type         string `json:"type"`
}

type RiskActionPlanControlDTO struct {
	ID               *uint                             `json:"id"`
	RiskActionPlanID uint                              `json:"risk_action_plan_id"`
	ControlID        uint                              `json:"control_id"`
	Control          organizationDTOs.ControlDetailDTO `json:"control"`
}
