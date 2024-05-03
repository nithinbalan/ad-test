package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"

	"time"
)

type IncidentRegisterListDTO struct {
	ID                         uint       `json:"id"`
	ReferenceNumber            string     `json:"reference_number"`
	Title                      string     `json:"title"`
	Description                *string    `json:"description"`
	IncidentAt                 time.Time  `json:"incident_at"`
	Location                   *string    `json:"location"`
	ReportedAt                 *time.Time `json:"reported_at"`
	QuickAction                *string    `json:"quick_action"`
	IncidentCategoryID         uint       `json:"incident_category_id"`
	IncidentSubCategoryID      uint       `json:"incident_sub_category_id"`
	IncidentTypeID             uint       `json:"incident_type_id"`
	SeverityLevelID            uint       `json:"severity_level_id"`
	IncidentStatusID           uint       `json:"incident_status_id"`
	IncidentCategoryTitleEn    string     `json:"incident_category_title_en"`
	IncidentCategoryTitleAr    string     `json:"incident_category_title_ar"`
	IncidentSubCategoryTitleEn string     `json:"incident_sub_category_title_en"`
	IncidentSubCategoryTitleAr string     `json:"incident_sub_category_title_ar"`
	IncidentTypeTitleEn        string     `json:"incident_type_title_en"`
	IncidentTypeTitleAr        string     `json:"incident_type_title_ar"`
	SeverityLevelTitleEn       string     `json:"severity_level_title_en"`
	SeverityLevelTitleAr       string     `json:"severity_level_title_ar"`
	IncidentStatusTitleEn      string     `json:"incident_status_title_en"`
	IncidentStatusTitleAr      string     `json:"incident_status_title_ar"`
	IncidentStatusColorCode    string     `json:"incident_status_color_code"`
	SeverityLevelColorCode     string     `json:"severity_level_color_code"`
	ReportedUserFirstNameEn    string     `json:"reported_user_first_name_en"`
	ReportedUserFirstNameAr    string     `json:"reported_user_first_name_ar"`
	ReportedUserLastNameEn     string     `json:"reported_user_last_name_en"`
	ReportedUserLastNameAr     string     `json:"reported_user_last_name_ar"`
	CreatedBy                  uint       `json:"created_by"`
	UpdatedBy                  *uint      `json:"updated_by"`
	CreatedAt                  string     `json:"created_at"`
	UpdatedAt                  string     `json:"updated_at"`
	CreatedUserFirstNameEn     string     `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr     string     `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn      string     `json:"created_user_last_name_en"`
	CreatedUserLastNameAr      string     `json:"created_user_last_name_ar"`
	UpdatedUserFirstNameEn     string     `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr     string     `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn      string     `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr      string     `json:"updated_user_last_name_ar"`
}

type IncidentRegisterFormDTO struct {
	IncidentCategoryID     uint                       `json:"incident_category_id" binding:"required"`
	IncidentSubCategoryID  uint                       `json:"incident_sub_category_id" binding:"required"`
	IncidentTypeID         uint                       `json:"incident_type_id" binding:"required"`
	SeverityLevelID        uint                       `json:"severity_level_id" binding:"required"`
	IncidentStatusID       uint                       `json:"incident_status_id" binding:"required"`
	Title                  string                     `json:"title" binding:"required,max=255"`
	Description            *string                    `json:"description"`
	IncidentAt             string                     `json:"incident_at" binding:"required,datetime=2006-01-02"`
	Location               *string                    `json:"location"`
	ReportedBy             uint                       `json:"reported_by" binding:"required"`
	ReportedAt             string                     `json:"reported_at" binding:"required,datetime=2006-01-02"`
	QuickAction            *string                    `json:"quick_action"`
	Documents              []*IncidentDocumentFormDTO `json:"documents" binding:"required"`
	ComplianceRequirements []uint                     `json:"compliance_requirements"`
	Controls               []uint                     `json:"controls"`
}

type IncidentDocumentFormDTO struct {
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

type IncidentDocumentDTO struct {
	ID           int    `json:"id"`
	Language     string `json:"language"`
	URL          string `json:"url"`
	Icon         string `json:"icon"`
	FileName     string `json:"file_name"`
	OriginalName string `json:"original_name"`
	Type         string `json:"type"`
	Size         int    `json:"size"`
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
	ID                  int                          `json:"id"`
	ReferenceNumber     string                       `json:"reference_number"`
	Title               string                       `json:"title"`
	Description         string                       `json:"description"`
	IncidentAt          time.Time                    `json:"incident_at"`
	Location            *string                      `json:"location"`
	ReportedAt          *time.Time                   `json:"reported_at"`
	QuickAction         *string                      `json:"quick_action"`
	IncidentCategory    IncidentCategoryDTO          `json:"incident_category"`
	IncidentSubCategory IncidentSubCategoryDTO       `json:"incident_sub_category"`
	IncidentType        IncidentTypeDTO              `json:"incident_type"`
	SeverityLevel       *masterDTOs.SeverityLevelDTO `json:"severity_level"`
	IncidentStatus      IncidentStatusDTO            `json:"incident_status"`
	ReportedByUser      UserMinimalDTO               `json:"reported_by"`
	IncidentDocuments   []IncidentDocumentDTO        `json:"incident_documents"`
	CreatedAt           time.Time                    `json:"created_at"`
	UpdatedAt           *time.Time                   `json:"updated_at"`
	CreatedByUser       UserMinimalDTO               `json:"created_by"`
	UpdatedByUser       UserMinimalDTO               `json:"updated_by"`
}

type UserMinimalDTO struct {
	ID                int                 `json:"id"`
	FirstNameEN       string              `json:"first_name_en"`
	LastNameEN        string              `json:"last_name_en"`
	FirstNameAR       string              `json:"first_name_ar"`
	LastNameAR        string              `json:"last_name_ar"`
	OrganizationUnit1 OrganizationUnitDTO `json:"organization_unit_1"`
}

type OrganizationUnitDTO struct {
	ID      int    `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Status  string `json:"status"`
}

type IncidentRegisterParsedDates struct {
	ReportedAt *time.Time
	IncidentAt *time.Time
}
