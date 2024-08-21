package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	complianceDTOs "jethur/grc/compliance/dtos"
	incidentDTOs "jethur/grc/incident/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	riskDTOs "jethur/grc/risk/dtos"
	"time"
)

type MeetingListDTO struct {
	ID                            uint       `json:"id"`
	ReferenceNumber               string     `json:"reference_number"`
	MeetingID                     *uint      `json:"meeting_id"`
	MeetingTypeTitleEn            string     `json:"meeting_type_title_en"`
	MeetingTypeTitleAr            string     `json:"meeting_type_title_ar"`
	MeetingFrequencyTitleEn       string     `json:"meeting_frequency_title_en"`
	MeetingFrequencyTitleAr       string     `json:"meeting_frequency_title_ar"`
	Title                         string     `json:"title"`
	Description                   *string    `json:"description"`
	Location                      *string    `json:"location"`
	Date                          time.Time  `json:"date"`
	StartTime                     *string    `json:"start_time"`
	EndTime                       *string    `json:"end_time"`
	OrganizerUserID               uint       `json:"organizer_user_id"`
	OrganizerUserFirstNameEn      string     `json:"organizer_user_first_name_en"`
	OrganizerUserFirstNameAr      string     `json:"organizer_user_first_name_ar"`
	OrganizerUserLastNameEn       string     `json:"organizer_user_last_name_en"`
	OrganizerUserLastNameAr       string     `json:"organizer_user_last_name_ar"`
	OrganizerUserInitial          string     `json:"organizer_user_initial"`
	OrganizerUserColor            string     `json:"organizer_user_color"`
	OrganizerUserEmail            string     `json:"organizer_user_email"`
	OrganizerUserProfileImagePath string     `json:"organizer_user_profile_image_path"`
	CreatedAt                     time.Time  `json:"created_at"`
	UpdatedAt                     *time.Time `json:"updated_at"`
	MeetingStatusTitleEn          string     `json:"meeting_status_title_en"`
	MeetingStatusTitleAr          string     `json:"meeting_status_title_ar"`
	CreatedUserID                 uint       `json:"created_user_id"`
	CreatedUserFirstNameEn        string     `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr        string     `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn         string     `json:"created_user_last_name_en"`
	CreatedUserLastNameAr         string     `json:"created_user_last_name_ar"`
	CreatedUserInitial            string     `json:"created_user_initial"`
	CreatedUserColor              string     `json:"created_user_color"`
	CreatedUserEmail              string     `json:"created_user_email"`
	CreatedUserProfileImagePath   string     `json:"created_user_profile_image_path"`
	UpdatedUserID                 *uint      `json:"updated_user_id"`
	UpdatedUserFirstNameEn        string     `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr        string     `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn         string     `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr         string     `json:"updated_user_last_name_ar"`
	UpdatedUserInitial            string     `json:"updated_user_initial"`
	UpdatedUserColor              string     `json:"updated_user_color"`
	UpdatedUserEmail              string     `json:"updated_user_email"`
	UpdatedUserProfileImagePath   string     `json:"updated_user_profile_image_path"`
	IsPublic                      bool       `json:"is_public"`
}

type MeetingListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
	DateFrom    string
	DateTo      string
}

type MeetingListResponseDTO struct {
	Items      []*MeetingListDTO            `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type MeetingDetailDTO struct {
	ID                     uint                               `json:"id"`
	ReferenceNumber        string                             `json:"reference_number"`
	MeetingID              *uint                              `json:"meeting_id"`
	MeetingTypeID          uint                               `json:"meeting_type_id"`
	MeetingFrequencyID     *uint                              `json:"meeting_frequency_id"`
	Meeting                *MeetingDetailDTO                  `json:"meeting"`
	MeetingType            MeetingTypeDTO                     `json:"meeting_type"`
	MeetingFrequency       *MeetingFrequencyDTO               `json:"meeting_frequency"`
	Title                  string                             `json:"title"`
	Description            *string                            `json:"description"`
	Location               *string                            `json:"location"`
	DateString             *string                            `json:"date"`
	StartTime              *string                            `json:"start_time"`
	EndTime                *string                            `json:"end_time"`
	IsPublic               bool                               `json:"is_public"`
	Organizer              organizationDTOs.UserMinimalDTO    `json:"organizer"`
	CreatedByUser          organizationDTOs.UserMinimalDTO    `json:"created_by"`
	UpdatedByUser          *organizationDTOs.UserMinimalDTO   `json:"updated_by"`
	CreatedAt              time.Time                          `json:"created_at"`
	UpdatedAt              *time.Time                         `json:"updated_at"`
	MeetingStatusID        uint                               `json:"meeting_status_id"`
	MeetingStatus          MeetingStatusDTO                   `json:"meeting_status"`
	Agendas                []*MeetingAgendaDetailDTO          `json:"agendas"`
	Participants           []*MeetingParticipantDTO           `json:"participants"`
	Controls               []*MeetingControlDTO               `json:"controls"`
	ComplianceRequirements []*MeetingComplianceRequirementDTO `json:"compliance_requirements"`
	Risks                  []*MeetingRiskDTO                  `json:"risks"`
	Processes              []*MeetingProcessDTO               `json:"processes"`
	Incidents              []*MeetingIncidentDTO              `json:"incidents"`
}

type MeetingParticipantDTO struct {
	ID        *uint                           `json:"id"`
	MeetingID uint                            `json:"meeting_id"`
	UserID    uint                            `json:"user_id"`
	User      organizationDTOs.UserMinimalDTO `json:"user"`
}

type MeetingFormDTO struct {
	Title                    string  `json:"title" binding:"required"`
	Description              *string `json:"description"`
	MeetingTypeID            uint    `json:"meeting_type_id" binding:"required"`
	MeetingFrequencyID       *uint   `json:"meeting_frequency_id"`
	MeetingID                *uint   `json:"meeting_id"`
	Location                 *string `json:"location"`
	Date                     string  `json:"date" binding:"required,datetime=2006-01-02"`
	StartTime                *string `json:"start_time"`
	EndTime                  *string `json:"end_time"`
	OrganizerID              uint    `json:"organizer_id" binding:"required"`
	IsPublic                 bool    `json:"is_public"`
	MeetingStatusID          uint    `json:"meeting_status_id" binding:"required"`
	UserIDs                  []*uint `json:"user_ids"`
	ComplianceRequirementIds []*uint `json:"compliance_requirement_ids"`
	ControlIds               []*uint `json:"control_ids"`
	ProcessIds               []*uint `json:"process_ids"`
	RiskIds                  []*uint `json:"risk_ids"`
	IncidentIds              []*uint `json:"incident_ids"`
	// Agendas                  []*MeetingAgendasFormDTO    `json:"agendas" binding:"required"`
	// MeetingMinutes           []*MeetingMinuteFormDTO     `json:"meeting_minutes"`
	ActionPlans []*MeetingActionPlanFormDTO `json:"action_plans"`
}

type MeetingParsedDates struct {
	Date      time.Time
	StartTime time.Time
	EndTime   time.Time
}
type MeetingAgendasFormDTO struct {
	ID       *uint   `json:"id,omitempty"`
	Title    string  `json:"title" binding:"required"`
	TimeFrom *string `json:"time_from"`
	TimeTo   *string `json:"time_to"`
}

type MeetingControlDTO struct {
	ID        *uint                             `json:"id"`
	MeetingID uint                              `json:"meeting_id"`
	ControlID uint                              `json:"control_id"`
	Control   organizationDTOs.ControlDetailDTO `json:"control"`
}

type MeetingComplianceRequirementDTO struct {
	ID                      *uint                                         `json:"id"`
	MeetingID               uint                                          `json:"meeting_id"`
	ComplianceRequirementID uint                                          `json:"compliance_requirement_id"`
	ComplianceRequirement   complianceDTOs.ComplianceRequirementDetailDTO `json:"compliance_requirement"`
}

type MeetingRiskDTO struct {
	ID             *uint                          `json:"id"`
	MeetingID      uint                           `json:"meeting_id"`
	RiskRegisterID uint                           `json:"risk_register_id"`
	RiskRegister   riskDTOs.RiskRegisterDetailDTO `json:"risk"`
}

type MeetingProcessDTO struct {
	ID         *uint                             `json:"id"`
	MeetingID  uint                              `json:"meeting_id"`
	ProocessID uint                              `json:"process_id"`
	Process    organizationDTOs.ProcessDetailDTO `json:"process"`
}

type MeetingIncidentDTO struct {
	ID                 *uint                                  `json:"id"`
	MeetingID          uint                                   `json:"meeting_id"`
	IncidentRegisterID uint                                   `json:"incident_register_id"`
	IncidentRegister   incidentDTOs.IncidentRegisterDetailDTO `json:"incident_register"`
}

type MeetingsCountResponseDTO struct {
	Agenda     int `json:"agenda"`      // Count of meeting agendas
	Mom        int `json:"mom"`         // Count of minutes of meeting (MoM)
	ActionPlan int `json:"action_plan"` // Count of meeting action plans
}
