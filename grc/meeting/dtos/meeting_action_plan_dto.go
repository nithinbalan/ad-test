package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"

	"time"
)

type MeetingActionPlanListDTO struct {
	ID                                  uint   `json:"id"`
	Title                               string `json:"title"`
	MeetingAgendaID                     *uint  `json:"meeting_agenda_id"`
	MeetingAgendaTitle                  string `json:"meeting_agenda_title"`
	Description                         string `json:"description"`
	SeverityLevelTitleEn                string `json:"severity_level_title_en"`
	SeverityLevelTitleAr                string `json:"severity_level_title_ar"`
	SeverityLevelColorCode              string `json:"severity_level_color_code"`
	ResponsibleUserID                   uint   `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn          string `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr          string `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn           string `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr           string `json:"responsible_user_last_name_ar"`
	ResponsibleUserUserInitial          string `json:"responsible_user_user_initial"`
	ResponsibleUserUserColor            string `json:"responsible_user_user_color"`
	ResponsibleUserUserEmail            string `json:"responsible_user_user_email"`
	ResponsibleUserUserProfileImagePath string `json:"responsible_user_user_profile_image_path"`
	TargetDate                          string `json:"target_date"`
	Percentage                          uint   `json:"percentage"`
	Comments                            string `json:"comments"`
	ActionPlanStatusTitleEn             string `json:"action_plan_status_title_en"`
	ActionPlanStatusTitleAr             string `json:"action_plan_status_title_ar"`
	ActionPlanStatusColorCode           string `json:"action_plan_status_color_code"`
	CreatedUserID                       uint   `json:"created_user_id"`
	CreatedUserFirstNameEn              string `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr              string `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn               string `json:"created_user_last_name_en"`
	CreatedUserLastNameAr               string `json:"created_user_last_name_ar"`
	CreatedUserInitial                  string `json:"created_user_initial"`
	CreatedUserColor                    string `json:"created_user_color"`
	CreatedUserEmail                    string `json:"created_user_email"`
	CreatedUserProfileImagePath         string `json:"created_user_profile_image_path"`
	UpdatedUserID                       uint   `json:"updated_user_id"`
	UpdatedUserFirstNameEn              string `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr              string `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn               string `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr               string `json:"updated_user_last_name_ar"`
	UpdatedUserInitial                  string `json:"updated_user_initial"`
	UpdatedUserColor                    string `json:"updated_user_color"`
	UpdatedUserEmail                    string `json:"updated_user_email"`
	UpdatedUserProfileImagePath         string `json:"updated_user_profile_image_path"`
}

type MeetingActionPlanListQueryParamsDTO struct {
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
}

type MeetingActionPlanListResponseDTO struct {
	Items      []*MeetingActionPlanListDTO  `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type MeetingActionPlanFormDTO struct {
	Title              string  `json:"title" binding:"required"`
	Description        *string `json:"description" binding:"omitempty"`
	SeverityLevelID    *uint   `json:"severity_level_id" binding:"omitempty"`
	ResponsibleUserID  uint    `json:"responsible_user_id" binding:"required"`
	TargetDate         *string `json:"target_date" binding:"omitempty,datetime=2006-01-02"`
	Percentage         *uint   `json:"percentage" binding:"omitempty"`
	Comments           *string `json:"comments" binding:"omitempty"`
	ActionPlanStatusID uint    `json:"action_plan_status_id" binding:"required"`
	MeetingAgendaID    *uint   `json:"meeting_agenda_id"`
}

type MeetingActionPlanParsedDates struct {
	TargetDate *time.Time
}

type MeetingActionPlanDetailDTO struct {
	ID                      uint                                `json:"id"`
	Title                   string                              `json:"title"`
	Description             string                              `json:"description"`
	SeverityLevel           *masterDTOs.SeverityLevelDTO        `json:"severity_level"`
	ResponsibleUser         organizationDTOs.UserMinimalDTO     `json:"responsible_user"`
	ActionPlanStatus        masterDTOs.ActionPlanStatusDTO      `json:"action_plan_status"`
	MeetingAgenda           *MeetingAgendaDetailDTO             `json:"agenda"`
	Percentage              *uint                               `json:"percentage"`
	Comments                *string                             `json:"comments"`
	TargetDate              *time.Time                          `json:"target_date"`
	CreatedAt               time.Time                           `json:"created_at"`
	UpdatedAt               *time.Time                          `json:"updated_at"`
	CreatedByUser           organizationDTOs.UserMinimalDTO     `json:"created_by_user"`
	UpdatedByUser           *organizationDTOs.UserMinimalDTO    `json:"updated_by_user"`
	ActionPlanStatusUpdates []*MeetingActionPlanStatusUpdateDTO `json:"action_plan_status_updates"`
}
type MeetingActionPlanStatusUpdateDTO struct {
	ID               uint                                        `json:"id"`
	ActionPlanStatus masterDTOs.ActionPlanStatusDTO              `json:"action_plan_status"`
	Comment          string                                      `json:"comment"`
	Percentage       *uint                                       `json:"percentage"`
	CreatedAt        time.Time                                   `json:"created_at"`
	UpdatedAt        *time.Time                                  `json:"updated_at"`
	CreatedByUser    organizationDTOs.UserMinimalDTO             `json:"created_by_user"`
	UpdatedByUser    *organizationDTOs.UserMinimalDTO            `json:"updated_by_user"`
	Documents        []*MeetingActionPlanStatusUpdateDocumentDTO `json:"documents"`
}
type MeetingActionPlanStatusUpdateDocumentDTO struct {
	ID           uint   `json:"id"`
	URL          string `json:"url"`
	Name         string `json:"name"`
	Size         int    `json:"size"`
	Token        string `json:"token"`
	OriginalName string `json:"original_name"`
	Type         string `json:"type"`
}

type SendCreateActionPlanMailDTO struct {
	FirstNameEn string `json:"first_name_en"`
	LastNameEn  string `json:"last_name_en"`
	Subject     string `json:"subject"`
	MessageEn   string `json:"message_en"`
	RedirectURL string `json:"redirect_url"`
}
