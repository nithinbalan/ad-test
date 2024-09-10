package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	settingsDTOs "jethur/grc/settings/dtos"
	"time"
)

type SLARegisterListDTO struct {
	ID                              uint      `json:"id"`
	ReferenceNumber                 string    `json:"reference_number"`
	LifecycleStatusID               uint      `json:"lifecycle_status_id"`
	TitleEn                         string    `json:"title_en"`
	TitleAr                         string    `json:"title_ar"`
	Description                     string    `json:"description"`
	Client                          string    `json:"client"`
	CreatedAt                       time.Time `json:"created_at"`
	UpdatedAt                       time.Time `json:"updated_at"`
	SLACategoryTitleEn              string    `json:"sla_category_title_en"`
	SLACategoryTitleAr              string    `json:"sla_category_title_ar"`
	SeverityLevelTitleEn            string    `json:"severity_level_title_en"`
	SeverityLevelTitleAr            string    `json:"severity_level_title_ar"`
	SeverityLevelColorCode          string    `json:"severity_level_color_code"`
	ResponsibleUnitTitleEn          string    `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr          string    `json:"responsible_unit_title_ar"`
	ResponsibleUserID               uint      `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn      string    `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr      string    `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn       string    `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr       string    `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          string    `json:"responsible_user_initial"`
	ResponsibleUserColor            string    `json:"responsible_user_color"`
	ResponsibleUserEmail            string    `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath string    `json:"responsible_user_profile_image_path"`
	LifecycleStatusTitleEn          string    `json:"lifecycle_status_title_en"`
	LifecycleStatusTitleAr          string    `json:"lifecycle_status_title_ar"`
	LifecycleStatusColorCode        string    `json:"lifecycle_status_color_code"`
	ComplianceStatusTitleEn         string    `json:"compliance_status_title_en"`
	ComplianceStatusTitleAr         string    `json:"compliance_status_title_ar"`
	ComplianceStatusColorCode       string    `json:"compliance_status_color_code"`
	CreatedUserID                   uint      `json:"created_user_id"`
	CreatedUserFirstNameEn          string    `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr          string    `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn           string    `json:"created_user_last_name_en"`
	CreatedUserLastNameAr           string    `json:"created_user_last_name_ar"`
	CreatedUserInitial              string    `json:"created_user_initial"`
	CreatedUserColor                string    `json:"created_user_color"`
	CreatedUserEmail                string    `json:"created_user_email"`
	CreatedUserProfileImagePath     string    `json:"created_user_profile_image_path"`
	UpdatedUserID                   *uint     `json:"updated_user_id"`
	UpdatedUserFirstNameEn          string    `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr          string    `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn           string    `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr           string    `json:"updated_user_last_name_ar"`
	UpdatedUserInitial              string    `json:"updated_user_initial"`
	UpdatedUserColor                string    `json:"updated_user_color"`
	UpdatedUserEmail                string    `json:"updated_user_email"`
	UpdatedUserProfileImagePath     string    `json:"updated_user_profile_image_path"`
	// Fields from sla_documents table
	Version      string `json:"version"`
	IssueDate    string `json:"issue_date"`
	ExpiryDate   string `json:"expiry_date"`
	SA1Date      string `json:"sa1_date"`
	SA2Date      string `json:"sa2_date"`
	Comment      string `json:"comment"`
	URL          string `json:"url"`
	Name         string `json:"name"`
	Size         int    `json:"size"`
	Token        string `json:"token"`
	Type         string `json:"type"`
	OriginalName string `json:"original_name"`
	IsLatest     bool   `json:"is_latest"`
}

type SLARegisterListQueryParamsDTO struct {
	Sort           string
	SortBy         string
	Keyword        string
	Status         string
	Page           int
	PageSize       int
	Filters        map[string][]string
	QuickFilter    []string
	IssueDateFrom  string
	IssueDateTo    string
	ExpiryDateFrom string
	ExpiryDateTo   string
}

type SLARegisterListResponseDTO struct {
	Items      []*SLARegisterListDTO        `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type SLARegisterFormDTO struct {
	SLACategoryID     uint               `json:"sla_category_id" binding:"required"`
	SeverityLevelID   uint               `json:"severity_level_id"`
	ResponsibleUserID uint               `json:"responsible_user_id"`
	LifecycleStatusID uint               `json:"lifecycle_status_id"`
	TitleEn           string             `json:"title_en" binding:"required,max=250"`
	TitleAr           string             `json:"title_ar" binding:"required,max=250"`
	Description       *string            `json:"description"`
	Client            string             `json:"client" binding:"max=250"`
	Version           string             `json:"version" binding:"required"`
	IssueDate         string             `json:"issue_date" binding:"omitempty,datetime=2006-01-02"`
	ExpiryDate        string             `json:"expiry_date" binding:"omitempty,datetime=2006-01-02"`
	SA1Date           string             `json:"sa1_date" binding:"omitempty,datetime=2006-01-02"`
	SA2Date           string             `json:"sa2_date" binding:"omitempty,datetime=2006-01-02"`
	Comment           *string            `json:"comment"`
	Document          SLADocumentFormDTO `json:"document" binding:"required"`
	ProcessIDs        []*uint            `json:"process_ids"`
}

type SLADocumentFormDTO struct {
	ID           *uint  `json:"id,omitempty"`
	URL          string `json:"url"`
	FileName     string `json:"file_name"`
	OriginalName string `json:"original_name"`
	Type         string `json:"type"`
	Size         int    `json:"size"`
	IsNewFile    *bool  `json:"is_new_file,omitempty"` // Optional field for indicating if the file is new or updated
}

type SLADocumentDTO struct {
	ID               *uint      `json:"id"`
	Version          string     `json:"version"`
	IssueDate        *time.Time `json:"-"`
	IssueDateString  *string    `json:"issue_date,omitempty"`
	ExpiryDate       *time.Time `json:"-"`
	ExpiryDateString *string    `json:"expiry_date,omitempty"`
	SA1Date          *time.Time `json:"-"`
	SA1DateString    *string    `json:"sa1_date,omitempty"`
	SA2Date          *time.Time `json:"-"`
	SA2DateString    *string    `json:"sa2_date,omitempty"`
	Comment          *string    `json:"comment,omitempty"`
	URL              string     `json:"url"`
	Name             string     `json:"name"`
	Size             int        `json:"size"`
	Token            string     `json:"token"`
	OriginalName     string     `json:"original_name"`
	Type             string     `json:"type"`
	IsLatest         bool       `json:"is_latest"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

type ParsedSLADates struct {
	IssueDate  *time.Time
	ExpiryDate *time.Time
	SA1Date    *time.Time
	SA2Date    *time.Time
}

type SLARegisterDetailDTO struct {
	ID               uint                             `json:"id"`
	ReferenceNumber  string                           `json:"reference_number"`
	SLACategory      SLACategoryDTO                   `json:"sla_category"`
	SeverityLevel    *masterDTOs.SeverityLevelDTO     `json:"severity_level"`
	ResponsibleUnit  settingsDTOs.OrganizationUnitDTO `json:"responsible_unit"`
	ResponsibleUser  organizationDTOs.UserMinimalDTO  `json:"responsible_user"`
	LifecycleStatus  masterDTOs.LifecycleStatusDTO    `json:"lifecycle_status"`
	ComplianceStatus *ComplianceStatusDTO             `json:"compliance_status,omitempty"`
	TitleEn          string                           `json:"title_en"`
	TitleAr          string                           `json:"title_ar"`
	Description      *string                          `json:"description,omitempty"`
	Client           string                           `json:"client"`
	CreatedAt        time.Time                        `json:"created_at"`
	UpdatedAt        *time.Time                       `json:"updated_at"`
	CreatedByUser    organizationDTOs.UserMinimalDTO  `json:"created_by"`
	UpdatedByUser    *organizationDTOs.UserMinimalDTO `json:"updated_by,omitempty"`
	Documents        []*SLADocumentDTO                `json:"documents"`
	Checks           []*SLAChecksDTO                  `json:"checks"`
	Processes        []*SLAProcessDTO                 `json:"processes"`
}

type SLAChecksDTO struct {
	ID            *uint     `json:"id"`
	SLARegisterID uint      `json:"sla_register_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type SLAProcessDTO struct {
	ID            *uint                             `json:"id"`
	SLARegisterID uint                              `json:"sla_register_id"`
	ProcessID     uint                              `json:"process_id"`
	Process       organizationDTOs.ProcessDetailDTO `json:"process"`
}

type SLACheckAnswerStatusCountResponseDTO struct {
	ID        int    `json:"id"`
	Type      string `json:"type"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	ColorCode string `json:"color_code"`
	Count     int    `json:"count"`
}
