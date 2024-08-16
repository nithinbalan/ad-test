package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	complianceDTOs "jethur/grc/compliance/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	riskDTOs "jethur/grc/risk/dtos"
	settingsDTOs "jethur/grc/settings/dtos"

	"time"
)

type IncidentRegisterListDTO struct {
	ID                              uint    `json:"id"`
	ReferenceNumber                 string  `json:"reference_number"`
	Title                           string  `json:"title"`
	Description                     *string `json:"description"`
	IncidentAt                      string  `json:"incident_at"`
	Location                        *string `json:"location"`
	ReportedAt                      *string `json:"reported_at"`
	QuickAction                     *string `json:"quick_action"`
	IncidentCategoryID              uint    `json:"incident_category_id"`
	IncidentTypeID                  uint    `json:"incident_type_id"`
	SeverityLevelID                 uint    `json:"severity_level_id"`
	IncidentStatusID                uint    `json:"incident_status_id"`
	ResponsibleUnitTitleEn          string  `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr          string  `json:"responsible_unit_title_ar"`
	ResponsibleUserFirstNameEn      string  `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr      string  `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn       string  `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr       string  `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          string  `json:"responsible_user_initial"`
	ResponsibleUserColor            string  `json:"responsible_user_color"`
	ResponsibleUserEmail            string  `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath string  `json:"responsible_user_profile_image_path"`
	IncidentCategoryTitleEn         string  `json:"incident_category_title_en"`
	IncidentCategoryTitleAr         string  `json:"incident_category_title_ar"`
	IncidentTypeTitleEn             string  `json:"incident_type_title_en"`
	IncidentTypeTitleAr             string  `json:"incident_type_title_ar"`
	SeverityLevelTitleEn            string  `json:"severity_level_title_en"`
	SeverityLevelTitleAr            string  `json:"severity_level_title_ar"`
	IncidentStatusTitleEn           string  `json:"incident_status_title_en"`
	IncidentStatusTitleAr           string  `json:"incident_status_title_ar"`
	IncidentStatusColorCode         string  `json:"incident_status_color_code"`
	SeverityLevelColorCode          string  `json:"severity_level_color_code"`
	ReportedUserFirstNameEn         string  `json:"reported_user_first_name_en"`
	ReportedUserFirstNameAr         string  `json:"reported_user_first_name_ar"`
	ReportedUserLastNameEn          string  `json:"reported_user_last_name_en"`
	ReportedUserLastNameAr          string  `json:"reported_user_last_name_ar"`
	CreatedAt                       string  `json:"created_at"`
	UpdatedAt                       string  `json:"updated_at"`
	CreatedUserID                   uint    `json:"created_user_id"`
	CreatedUserFirstNameEn          string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr          string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn           string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr           string  `json:"created_user_last_name_ar"`
	UpdatedUserID                   *uint   `json:"updated_user_id"`
	UpdatedUserFirstNameEn          string  `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr          string  `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn           string  `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr           string  `json:"updated_user_last_name_ar"`
}

type IncidentRegisterFormDTO struct {
	IncidentCategoryID       uint                       `json:"incident_category_id" binding:"required"`
	IncidentTypeID           *uint                      `json:"incident_type_id"`
	SeverityLevelID          *uint                      `json:"severity_level_id"`
	IncidentStatusID         *uint                      `json:"incident_status_id"`
	Title                    string                     `json:"title" binding:"required,max=255"`
	ResponsibleUserID        uint                       `json:"responsible_user_id" binding:"required"`
	Description              *string                    `json:"description"`
	IncidentAt               string                     `json:"incident_at" binding:"required,datetime=2006-01-02"`
	Location                 *string                    `json:"location"`
	ReportedBy               uint                       `json:"reported_by" binding:"required"`
	ReportedAt               string                     `json:"reported_at" binding:"required,datetime=2006-01-02"`
	QuickAction              *string                    `json:"quick_action"`
	Documents                []*IncidentDocumentFormDTO `json:"documents" binding:"required"`
	ComplianceRequirementIds []*uint                    `json:"compliance_requirement_ids"`
	ControlIds               []*uint                    `json:"control_ids"`
	RiskIds                  []*uint                    `json:"risk_ids"`
	ProcessIDs               []*uint                    `json:"process_ids"`
	SLARegisterIDs           []*uint                    `json:"sla_register_ids"`
}

type IncidentDocumentFormDTO struct {
	ID            *uint  `json:"id,omitempty"`
	URL           string `json:"url"`
	FileName      string `json:"file_name"`
	OriginalName  string `json:"original_name"`
	Type          string `json:"type"`
	Size          int    `json:"size"`
	IsNewFile     *bool  `json:"is_new_file,omitempty"`
	IsDeletedFile *bool  `json:"is_deleted_file,omitempty"`
}

type IncidentDocumentDTO struct {
	ID           *uint   `json:"id"`
	URL          *string `json:"url"`
	Name         string  `json:"name"`
	Token        string  `json:"token"`
	OriginalName string  `json:"original_name"`
	Type         string  `json:"type"`
	Size         int     `json:"size"`
}

type IncidentRegisterListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
	IncidentAt  string
	ReportedAt  string
}

type IncidentRegisterListResponseDTO struct {
	Items      []*IncidentRegisterListDTO   `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type IncidentRegisterDetailDTO struct {
	ID                     uint                                `json:"id"`
	ReferenceNumber        string                              `json:"reference_number"`
	Title                  string                              `json:"title"`
	Description            string                              `json:"description"`
	IncidentAtString       string                              `json:"incident_at"`
	Location               *string                             `json:"location"`
	ReportedAtStaring      string                              `json:"reported_at"`
	QuickAction            *string                             `json:"quick_action"`
	IncidentCategory       IncidentCategoryDTO                 `json:"incident_category"`
	IncidentType           IncidentTypeDTO                     `json:"incident_type"`
	SeverityLevel          *masterDTOs.SeverityLevelDTO        `json:"severity_level"`
	IncidentStatus         IncidentStatusDTO                   `json:"incident_status"`
	ReportedByUser         organizationDTOs.UserMinimalDTO     `json:"reported_by"`
	Documents              []*IncidentDocumentDTO              `json:"documents"`
	Controls               []*IncidentControlDTO               `json:"controls"`
	ComplianceRequirements []*IncidentComplianceRequirementDTO `json:"compliance_requirements"`
	Risks                  []*IncidentRiskDTO                  `json:"risks"`
	Processes              []*IncidentProcessDTO               `json:"processes"`
	ResponsibleUnit        settingsDTOs.OrganizationUnitDTO    `json:"responsible_unit"`
	ResponsibleUser        organizationDTOs.UserMinimalDTO     `json:"responsible_user"`
	SLARegisters           []*IncidentSLARegisterDTO           `json:"sla_registers"`
	CreatedAt              time.Time                           `json:"created_at"`
	UpdatedAt              *time.Time                          `json:"updated_at"`
	CreatedByUser          organizationDTOs.UserMinimalDTO     `json:"created_by"`
	UpdatedByUser          organizationDTOs.UserMinimalDTO     `json:"updated_by"`
}
type IncidentRegisterParsedDates struct {
	ReportedAt *time.Time
	IncidentAt *time.Time
}

type IncidentControlDTO struct {
	ID                 *uint                             `json:"id"`
	IncidentRegisterID uint                              `json:"incident_register_id"`
	ControlID          uint                              `json:"control_id"`
	Control            organizationDTOs.ControlDetailDTO `json:"control"`
}

type IncidentComplianceRequirementDTO struct {
	ID                      *uint                                         `json:"id"`
	IncidentRegisterID      uint                                          `json:"incident_register_id"`
	ComplianceRequirementID uint                                          `json:"compliance_requirement_id"`
	ComplianceRequirement   complianceDTOs.ComplianceRequirementDetailDTO `json:"compliance_requirement"`
}

type IncidentRiskDTO struct {
	ID                 *uint                          `json:"id"`
	IncidentRegisterID uint                           `json:"incident_register_id"`
	RiskRegisterID     uint                           `json:"risk_register_id"`
	RiskRegister       riskDTOs.RiskRegisterDetailDTO `json:"risk"`
}

type IncidentProcessDTO struct {
	ID                 *uint                             `json:"id"`
	IncidentRegisterID uint                              `json:"incident_register_id"`
	ProcessID          uint                              `json:"process_id"`
	Process            organizationDTOs.ProcessDetailDTO `json:"process"`
}
type IncidentSLARegisterDTO struct {
	ID                 *uint                               `json:"id"`
	IncidentRegisterID uint                                `json:"incident_register_id"`
	SLARegisterID      uint                                `json:"sla_register_id"`
	SLARegister        complianceDTOs.SLARegisterDetailDTO `json:"sla_register"`
}
