package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type IncidentInvestigationListDTO struct {
	ID                                   uint      `json:"id"`
	IncidentRegisterID                   uint      `json:"incident_register_id"`
	IncidentTitle                        string    `json:"incident_title"`
	InitiatedAt                          time.Time `json:"initiated_at"`
	Description                          *string   `json:"description"`
	TeamLeadID                           uint      `json:"team_lead_id"`
	TeamLeadFirstNameEn                  string    `json:"team_lead_first_name_en"`
	TeamLeadFirstNameAr                  string    `json:"team_lead_first_name_ar"`
	TeamLeadLastNameEn                   string    `json:"team_lead_last_name_en"`
	TeamLeadLastNameAr                   string    `json:"team_lead_last_name_ar"`
	ActionRecommendation                 *string   `json:"action_recommendation"`
	IncidentInvestigationStatusID        int       `json:"incident_investigation_status_id"`
	IncidentInvestigationStatusTitleEn   string    `json:"incident_investigation_status_title_en"`
	IncidentInvestigationStatusTitleAr   string    `json:"incident_investigation_status_title_ar"`
	IncidentInvestigationStatusColorCode string    `json:"incident_investigation_status_color_code"`
	CreatedAt                            string    `json:"created_at"`
	UpdatedAt                            string    `json:"updated_at"`
	CreatedUserFirstNameEn               string    `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr               string    `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn                string    `json:"created_user_last_name_en"`
	CreatedUserLastNameAr                string    `json:"created_user_last_name_ar"`
	UpdatedUserFirstNameEn               string    `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr               string    `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn                string    `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr                string    `json:"updated_user_last_name_ar"`
}
type IncidentInvestigationListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
	InitiatedAt string
}

type IncidentInvestigationListResponseDTO struct {
	Items      []*IncidentInvestigationListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO    `json:"pagination"`
}

type IncidentInvestigationFormDTO struct {
	IncidentRegisterID            uint                          `json:"incident_register_id"`
	Description                   *string                       `json:"description"`
	InitiatedAt                   string                        `json:"initiated_at" binding:"required,datetime=2006-01-02"`
	TeamLeadID                    uint                          `json:"team_lead_id" binding:"required"`
	ActionRecommendation          *string                       `json:"action_recommendation"`
	IncidentInvestigationStatusID int                           `json:"incident_investigation_status_id" binding:"required"`
	Documents                     []*commonDTOs.DocumentFormDTO `json:"documents" binding:"required"`
	UserIds                       []*uint                       `json:"user_ids"`
	Findings                      []*string                     `json:"findings" binding:"required,dive,max=500" gorm:"type:varchar(500);not null"`
}
type IncidentInvestigationDocumentDTO struct {
	ID           *uint   `json:"id"`
	URL          *string `json:"url"`
	Name         string  `json:"name"`
	Token        string  `json:"token"`
	OriginalName string  `json:"original_name"`
	Type         string  `json:"type"`
	Size         int     `json:"size"`
}

type IncidentInvestigationDetailDTO struct {
	ID                          uint                                `json:"id"`
	IncidentRegister            IncidentRegisterDetailDTO           `json:"incident_register"`
	Description                 *string                             `json:"description"`
	InitiatedAt                 time.Time                           `json:"initiated_at"`
	ActionRecommendation        *string                             `json:"action_recommendation"`
	TeamLead                    organizationDTOs.UserMinimalDTO     `json:"team_lead"`
	IncidentInvestigationStatus IncidentInvestigationStatusDTO      `json:"incident_investigation_status"`
	Documents                   []*IncidentInvestigationDocumentDTO `json:"documents"`
	Users                       []*IncidentInvestigationUsersDTO    `json:"users"`
	Findings                    []*IncidentInvestigationFindingsDTO `json:"findings"`
	CreatedAt                   time.Time                           `json:"created_at"`
	UpdatedAt                   *time.Time                          `json:"updated_at"`
	CreatedByUser               organizationDTOs.UserMinimalDTO     `json:"created_by"`
	UpdatedByUser               organizationDTOs.UserMinimalDTO     `json:"updated_by"`
}

type IncidentInvestigationParsedDates struct {
	InitiatedAt *time.Time
}

type IncidentInvestigationUsersDTO struct {
	ID                      uint                            `json:"id"`
	IncidentInvestigationID uint                            `json:"incident_investigation_id"`
	UserID                  int                             `json:"user_id"`
	User                    organizationDTOs.UserDetailsDTO `json:"user"`
}

type IncidentInvestigationFindingsDTO struct {
	ID                      uint   `json:"id"`
	IncidentInvestigationID uint   `json:"incident_investigation_id"`
	Title                   string `json:"title"`
}
