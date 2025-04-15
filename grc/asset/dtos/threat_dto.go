package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	settingsDTOs "jethur/grc/settings/dtos"
	"time"
)

type ThreatListDTO struct {
	ID                              uint       `json:"id"`
	TitleEn                         string     `json:"title_en"`
	TitleAr                         string     `json:"title_ar"`
	ThreatTypeTitleEn               *string    `json:"threat_type_title_en"`
	ThreatTypeTitleAr               *string    `json:"threat_type_title_ar"`
	SeverityLevelTitleEn            *string    `json:"severity_level_title_en"`
	SeverityLevelTitleAr            *string    `json:"severity_level_title_ar"`
	SeverityLevelColorCode          *string    `json:"severity_level_color_code"`
	ThreatStatusTitleEn             string     `json:"threat_status_title_en"`
	ThreatStatusTitleAr             string     `json:"threat_status_title_ar"`
	ThreatStatusColorCode           string     `json:"threat_status_color_code"`
	ThreatStatusType                string     `json:"threat_status_type"`
	ResponsibleUnitTitleEn          string     `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr          string     `json:"responsible_unit_title_ar"`
	ResponsibleUserID               uint       `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn      string     `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr      string     `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn       string     `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr       string     `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          string     `json:"responsible_user_initial"`
	ResponsibleUserColor            string     `json:"responsible_user_color"`
	ResponsibleUserEmail            string     `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath string     `json:"responsible_user_profile_image_path"`
	Description                     *string    `json:"description,omitempty"`
	Source                          *string    `json:"source,omitempty"`
	DetectionDate                   *string    `json:"detection_date,omitempty"`
	ResolutionDate                  *string    `json:"resolution_date,omitempty"`
	CreatedUserID                   uint       `json:"created_user_id"`
	CreatedUserFirstNameEn          string     `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr          string     `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn           string     `json:"created_user_last_name_en"`
	CreatedUserLastNameAr           string     `json:"created_user_last_name_ar"`
	CreatedUserInitial              string     `json:"created_user_initial"`
	CreatedUserColor                string     `json:"created_user_color"`
	CreatedUserEmail                string     `json:"created_user_email"`
	CreatedUserProfileImagePath     string     `json:"created_user_profile_image_path"`
	CreatedAt                       time.Time  `json:"created_at"`
	UpdatedUserID                   *uint      `json:"updated_user_id"`
	UpdatedUserFirstNameEn          *string    `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr          *string    `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn           *string    `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr           *string    `json:"updated_user_last_name_ar"`
	UpdatedUserInitial              string     `json:"updated_user_initial"`
	UpdatedUserColor                string     `json:"updated_user_color"`
	UpdatedUserEmail                string     `json:"updated_user_email"`
	UpdatedUserProfileImagePath     string     `json:"updated_user_profile_image_path"`
	UpdatedAt                       *time.Time `json:"updated_at"`
}

type ThreatListQueryParamsDTO struct {
	Sort              string
	SortBy            string
	Keyword           string
	Status            string
	Page              int
	PageSize          int
	Filters           map[string][]string
	QuickFilter       []string
	DetectionDateFrom string
	DetectionDateTo   string
	Export            bool
}

type ThreatListResponseDTO struct {
	Items      []*ThreatListDTO             `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type ThreatFormDTO struct {
	TitleEn           string                        `json:"title_en" binding:"required"`
	TitleAr           string                        `json:"title_ar" binding:"required"`
	ThreatTypeID      uint                          `json:"threat_type_id" binding:"required"`
	SeverityLevelID   uint                          `json:"severity_level_id" binding:"required"`
	ThreatStatusID    uint                          `json:"threat_status_id" binding:"required"`
	ResponsibleUserID uint                          `json:"responsible_user_id" binding:"required"`
	VulnerabilityIDs  []*uint                       `json:"vulnerability_ids"`
	Description       *string                       `json:"description,omitempty"`
	Source            *string                       `json:"source,omitempty"`
	DetectionDate     *string                       `json:"detection_date,omitempty" binding:"omitempty,datetime=2006-01-02"`
	ResolutionDate    *string                       `json:"resolution_date,omitempty" binding:"omitempty,datetime=2006-01-02"`
	Documents         []*commonDTOs.DocumentFormDTO `json:"documents" binding:"omitempty"`
}

type ThreatParsedDates struct {
	DetectionDate  *time.Time
	ResolutionDate *time.Time
}

type ThreatDetailDTO struct {
	ID                   uint                             `json:"id"`
	TitleEn              string                           `json:"title_en"`
	TitleAr              string                           `json:"title_ar"`
	ThreatType           *ThreatTypeDTO                   `json:"threat_type"`
	SeverityLevel        masterDTOs.SeverityLevelDTO      `json:"severity_level"`
	ThreatStatus         ThreatStatusDTO                  `json:"threat_status"`
	ResponsibleUser      organizationDTOs.UserMinimalDTO  `json:"responsible_user"`
	ResponsibleUnit      settingsDTOs.OrganizationUnitDTO `json:"responsible_unit"`
	Description          *string                          `json:"description,omitempty"`
	Source               *string                          `json:"source,omitempty"`
	DetectionDateString  *string                          `json:"detection_date,omitempty"`
	ResolutionDateString *string                          `json:"resolution_date,omitempty"`
	CreatedAt            time.Time                        `json:"created_at"`
	UpdatedAt            *time.Time                       `json:"updated_at,omitempty"`
	CreatedByUser        organizationDTOs.UserMinimalDTO  `json:"created_by"`
	UpdatedByUser        *organizationDTOs.UserMinimalDTO `json:"updated_by,omitempty"`
	Documents            []*ThreatDocumentDTO             `json:"documents"`
	VulnerabilityThreats []*VulnerabilityThreatDTO        `json:"vulnerability_threats"`
}

type VulnerabilityThreatDTO struct {
	ID              uint              `json:"id"`
	VulnerabilityID uint              `json:"vulnerability_id"`
	ThreatID        uint              `json:"threat_id"`
	Vulnerability   *VulnerabilityDTO `json:"vulnerability"`
	CreatedByID     uint              `json:"created_by_id"`
	CreatedAt       time.Time         `json:"created_at"`
}

type ThreatDocumentDTO struct {
	ID           *uint   `json:"id,omitempty"`
	URL          *string `json:"url"`
	Name         string  `json:"file_name"`
	OriginalName string  `json:"original_name"`
	Type         string  `json:"type"`
	Size         int     `json:"size"`
}

type ThreatSummaryDTO struct {
	OpenCount   int `json:"open_count"`
	ClosedCount int `json:"closed_count"`
}
