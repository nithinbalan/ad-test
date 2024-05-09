package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type IncidentInvestigationListDTO struct {
	ID                                 uint      `json:"id"`
	IncidentRegisterID                 uint      `json:"incident_register_id"`
	IncidentTitle                      string    `json:"incident_title"`
	InitiatedAt                        time.Time `json:"initiated_at"`
	Description                        *string   `json:"description"`
	TeamLeadID                         uint      `json:"team_lead_id"`
	TeamLeadFirstNameEn                string    `json:"team_lead_first_name_en"`
	TeamLeadFirstNameAr                string    `json:"team_lead_first_name_ar"`
	TeamLeadLastNameEn                 string    `json:"team_lead_last_name_en"`
	TeamLeadLastNameAr                 string    `json:"team_lead_last_name_ar"`
	ActionRecommendation               *string   `json:"action_recommendation"`
	IncidentInvestigationStatusID      int       `json:"incident_investigation_status_id"`
	IncidentInvestigationStatusTitleEn string    `json:"incident_investigation_status_title_en"`
	IncidentInvestigationStatusTitleAr string    `json:"incident_investigation_status_title_ar"`
	CreatedAt                          string    `json:"created_at"`
	UpdatedAt                          string    `json:"updated_at"`
	CreatedUserFirstNameEn             string    `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr             string    `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn              string    `json:"created_user_last_name_en"`
	CreatedUserLastNameAr              string    `json:"created_user_last_name_ar"`
	UpdatedUserFirstNameEn             string    `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr             string    `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn              string    `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr              string    `json:"updated_user_last_name_ar"`
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
	IncidentRegisterID            uint                                    `json:"incident_register_id"`
	Description                   *string                                 `json:"description"`
	InitiatedAt                   string                                  `json:"initiated_at" binding:"required,datetime=2006-01-02"`
	TeamLeadID                    uint                                    `json:"team_lead_id" binding:"required"`
	ActionRecommendation          *string                                 `json:"action_recommendation"`
	IncidentInvestigationStatusID int                                     `json:"incident_investigation_status_id" binding:"required"`
	Documents                     []*IncidentInvestigationDocumentFormDTO `json:"documents" binding:"required"`
	Users                         []*uint                                 `json:"users"`
	Findings                      []*string                               `json:"findings" binding:"required,dive,max=500" gorm:"type:varchar(500);not null"`
}

type IncidentInvestigationDocumentFormDTO struct {
	ID            *uint  `json:"id,omitempty"`
	Language      string `json:"language" binding:"required"`
	URL           string `json:"url"`
	Icon          string `json:"icon" binding:"required"`
	FileName      string `json:"file_name"`
	OriginalName  string `json:"original_name"`
	Type          string `json:"type"`
	Size          int    `json:"size"`
	IsNewFile     *bool  `json:"is_new_file,omitempty"`
	IsDeletedFile *bool  `json:"is_deleted_file,omitempty"`
}
type IncidentInvestigationDocumentDTO struct {
	ID       *uint   `json:"id"`
	Language string  `json:"language"`
	URL      *string `json:"url"`
	Icon     string  `json:"icon"`
	Name     string  `json:"name"`
	Token    string  `json:"token"`
	Type     string  `json:"type"`
	Size     int     `json:"size"`
}

type IncidentInvestigationDetailDTO struct {
	IncidentRegister            IncidentRegisterDetailDTO           `json:"incident_register"`
	Description                 *string                             `json:"description"`
	InitiatedAt                 time.Time                           `json:"initiated_at"`
	ActionRecommendation        *string                             `json:"action_recommendation"`
	TeamLead                    organizationDTOs.UserMinimalDTO     `json:"team_lead"`
	IncidentInvestigationStatus IncidentInvestigationStatusDTO      `json:"incident_investigation_status"`
	Documents                   []*IncidentInvestigationDocumentDTO `json:"documents"`
	Users                       []*IncidentInvestigationUserDTO     `json:"users"`
	Findings                    []*IncidentInvestigationFindingDTO  `json:"findings"`
	CreatedAt                   time.Time                           `json:"created_at"`
	UpdatedAt                   *time.Time                          `json:"updated_at"`
	CreatedByUser               organizationDTOs.UserMinimalDTO     `json:"created_by"`
	UpdatedByUser               organizationDTOs.UserMinimalDTO     `json:"updated_by"`
}

type IncidentInvestigationParsedDates struct {
	InitiatedAt *time.Time
}
