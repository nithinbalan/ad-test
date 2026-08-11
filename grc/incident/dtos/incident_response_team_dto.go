package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"

	"time"
)

type IncidentResponseTeamListDTO struct {
	ID                           uint                              `json:"id"`
	TitleEn                      string                            `json:"title_en"`
	TitleAr                      string                            `json:"title_ar"`
	Scope                        *string                           `json:"scope"`
	Objective                    *string                           `json:"objective"`
	TeamLeadUserID               *uint                             `json:"team_lead_user_id"`
	TeamLeadUserFirstNameEn      string                            `json:"team_lead_user_first_name_en"`
	TeamLeadUserLastNameEn       string                            `json:"team_lead_user_last_name_en"`
	TeamLeadUserFirstNameAr      string                            `json:"team_lead_user_first_name_ar"`
	TeamLeadUserLastNameAr       string                            `json:"team_lead_user_last_name_ar"`
	TeamLeadUserEmail            string                            `json:"team_lead_user_email"`
	TeamLeadUserInitial          string                            `json:"team_lead_user_initial"`
	TeamLeadUserColor            string                            `json:"team_lead_user_color"`
	TeamLeadUserProfileImagePath *string                           `json:"team_lead_user_profile_image_path"`
	CreatedUserID                uint                              `json:"created_user_id"`
	CreatedUserFirstNameEn       string                            `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr       string                            `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn        string                            `json:"created_user_last_name_en"`
	CreatedUserLastNameAr        string                            `json:"created_user_last_name_ar"`
	CreatedUserEmail             string                            `json:"created_user_email"`
	CreatedUserInitial           string                            `json:"created_user_initial"`
	CreatedUserColor             string                            `json:"created_user_color"`
	CreatedUserProfileImagePath  *string                           `json:"created_user_profile_image_path"`
	UpdatedUserID                *uint                             `json:"updated_user_id"`
	UpdatedUserFirstNameEn       string                            `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr       string                            `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn        string                            `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr        string                            `json:"updated_user_last_name_ar"`
	UpdatedUserEmail             string                            `json:"updated_user_email"`
	UpdatedUserInitial           string                            `json:"updated_user_initial"`
	UpdatedUserColor             string                            `json:"updated_user_color"`
	UpdatedUserProfileImagePath  *string                           `json:"updated_user_profile_image_path"`
	CreatedAt                    string                            `json:"created_at"`
	UpdatedAt                    string                            `json:"updated_at"`
	Members                      []*IncidentResponseTeamMembersDTO `json:"members" gorm:"-"`
}

type IncidentResponseTeamListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
}

type IncidentResponseTeamListResponseDTO struct {
	Items      []*IncidentResponseTeamListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO   `json:"pagination"`
}

type IncidentResponseTeamFormDTO struct {
	TitleEn        string  `json:"title_en" binding:"required,max=500"`
	TitleAr        string  `json:"title_ar" binding:"required,max=500"`
	TeamLeadUserID *uint   `json:"team_lead_user_id"`
	UserIDs        []*uint `json:"user_ids" binding:"required"`
	Scope          *string `json:"scope"`
	Objective      *string `json:"objective"`
}

type IncidentResponseTeamDetailDTO struct {
	ID                          uint                              `json:"id"`
	TitleEn                     string                            `json:"title_en"`
	TitleAr                     string                            `json:"title_ar"`
	Scope                       *string                           `json:"scope"`
	Objective                   *string                           `json:"objective"`
	TeamLeadUser                *organizationDTOs.UserMinimalDTO  `json:"team_lead_user"`
	IncidentResponseTeamMembers []*IncidentResponseTeamMembersDTO `json:"members"`
	CreatedAt                   time.Time                         `json:"created_at"`
	UpdatedAt                   *time.Time                        `json:"updated_at"`
	CreatedByUser               *organizationDTOs.UserMinimalDTO  `json:"created_by"`
	UpdatedByUser               *organizationDTOs.UserMinimalDTO  `json:"updated_by"`
}

type IncidentResponseTeamMembersDTO struct {
	ID                     uint                            `json:"id"`
	IncidentResponseTeamID uint                            `json:"incident_response_team_id"`
	UserID                 uint                            `json:"user_id"`
	User                   organizationDTOs.UserMinimalDTO `json:"user"`
}
