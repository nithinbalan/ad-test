package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	settingsDTOs "jethur/grc/settings/dtos"
	"time"
)

type ThreatListDTO struct {
	ID                       uint      `json:"id"`
	TitleEn                  string    `json:"title_en"`
	TitleAr                  string    `json:"title_ar"`
	ThreatTypeTitleEn        *string   `json:"threat_type_title_en"`
	ThreatTypeTitleAr        *string   `json:"threat_type_title_ar"`
	SeverityLevelTitleEn     *string   `json:"severity_level_title_en"`
	SeverityLevelTitleAr     *string   `json:"severity_level_title_ar"`
	ThreatStatusTitleEn      string    `json:"threat_status_title_en"`
	ThreatStatusTitleAr      string    `json:"threat_status_title_ar"`
	ResponsibleUnitTitleEn   string    `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr   string    `json:"responsible_unit_title_ar"`
	ResponsibleUserID        uint      `json:"responsible_user_id"`
	ResponsibleUserFirstName string    `json:"responsible_user_first_name"`
	ResponsibleUserLastName  string    `json:"responsible_user_last_name"`
	DescriptionEn            *string   `json:"description_en,omitempty"`
	DescriptionAr            *string   `json:"description_ar,omitempty"`
	DetectionDate            *string   `json:"detection_date,omitempty"`
	ResolutionDate           *string   `json:"resolution_date,omitempty"`
	CreatedByID              uint      `json:"created_by_id"`
	CreatedAt                time.Time `json:"created_at"`
	UpdatedByID              *uint     `json:"updated_by_id,omitempty"`
	UpdatedAt                *time.Time `json:"updated_at,omitempty"`
}

type ThreatListQueryParamsDTO struct {
	Sort             string
	SortBy           string
	Keyword          string
	Status           string
	Page             int
	PageSize         int
	Filters          map[string][]string
	QuickFilter      []string
	DetectionDateFrom string
	DetectionDateTo   string
	Export           bool
}

type ThreatListResponseDTO struct {
	Items      []*ThreatListDTO           `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type ThreatFormDTO struct {
	TitleEn           string  `json:"title_en" binding:"required"`
	TitleAr           string  `json:"title_ar" binding:"required"`
	ThreatTypeID      uint    `json:"threat_type_id" binding:"required"`
	SeverityLevelID   uint    `json:"severity_level_id" binding:"required"`
	ThreatStatusID    uint    `json:"threat_status_id" binding:"required"`
	ResponsibleUserID uint    `json:"responsible_user_id" binding:"required"`
	ResponsibleUnitID uint    `json:"responsible_unit_id" binding:"required"`
	DescriptionEn     *string `json:"description_en,omitempty"`
	DescriptionAr     *string `json:"description_ar,omitempty"`
	DetectionDate     *string `json:"detection_date,omitempty" binding:"omitempty,datetime=2006-01-02"`
	ResolutionDate    *string `json:"resolution_date,omitempty" binding:"omitempty,datetime=2006-01-02"`
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
	DescriptionEn        *string                          `json:"description_en,omitempty"`
	DescriptionAr        *string                          `json:"description_ar,omitempty"`
	DetectionDateString  *string                          `json:"detection_date,omitempty"`
	ResolutionDateString *string                          `json:"resolution_date,omitempty"`
	CreatedAt            time.Time                        `json:"created_at"`
	UpdatedAt            *time.Time                       `json:"updated_at,omitempty"`
	CreatedByUser        organizationDTOs.UserMinimalDTO  `json:"created_by"`
	UpdatedByUser        *organizationDTOs.UserMinimalDTO `json:"updated_by,omitempty"`
}
