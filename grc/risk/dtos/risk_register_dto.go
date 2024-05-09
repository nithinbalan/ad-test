package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"

	"time"
)

type RiskRegisterListDTO struct {
	ID                      uint       `json:"id"`
	ReferenceNumber         string     `json:"reference_number"`
	RiskCategoryTitleEn     string     `json:"risk_category_title_en"`
	RiskCategoryTitleAr     string     `json:"risk_category_title_ar"`
	RiskSourceTitleEn       string     `json:"risk_source_title_en"`
	RiskSourceTitleAr       string     `json:"risk_source_title_ar"`
	OwnerUnitTitleEn        string     `json:"owner_unit_title_en"`
	OwnerUnitTitleAr        string     `json:"owner_unit_title_ar"`
	OwnerFirstNameEn        string     `json:"owner_first_name_en"`
	OwnerFirstNameAr        string     `json:"owner_first_name_ar"`
	OwnerLastNameEn         string     `json:"owner_last_name_en"`
	OwnerLastNameAr         string     `json:"owner_last_name_ar"`
	RiskStatusTitleEn       string     `json:"risk_status_title_en"`
	RiskStatusTitleAr       string     `json:"risk_status_title_ar"`
	RiskStatusColorCode     string     `json:"risk_status_color_code"`
	Title                   string     `json:"title"`
	Description             *string    `json:"description"`
	ReportedAt              *time.Time `json:"reported_at"`
	ReportedUserFirstNameEn *string    `json:"reported_user_first_name_en"`
	ReportedUserFirstNameAr *string    `json:"reported_user_first_name_ar"`
	ReportedUserLastNameEn  *string    `json:"reported_user_last_name_en"`
	ReportedUserLastNameAr  *string    `json:"reported_user_last_name_ar"`
	ReviewFrequencyTitleEn  *string    `json:"review_frequency_title_en"`
	ReviewFrequencyTitleAr  *string    `json:"review_frequency_title_ar"`
	NextReviewDate          *string    `json:"next_review_date"`
	Comments                *string    `json:"comments"`
	CreatedAt               string     `json:"created_at"`
	UpdatedAt               *string    `json:"updated_at"`
	CreatedUserFirstNameEn  string     `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr  string     `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn   string     `json:"created_user_last_name_en"`
	CreatedUserLastNameAr   string     `json:"created_user_last_name_ar"`
	UpdatedUserFirstNameEn  *string    `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr  *string    `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn   *string    `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr   *string    `json:"updated_user_last_name_ar"`
}

type RiskRegisterListQueryParamsDTO struct {
	Sort               string
	SortBy             string
	Keyword            string
	Status             string
	Page               int
	PageSize           int
	Filters            map[string][]string
	QuickFilter        []string
	ReportedAtFrom     string
	ReportedAtTo       string
	NextReviewDateFrom string
	NextReviewDateTo   string
}

type RiskRegisterListResponseDTO struct {
	Items      []*RiskRegisterListDTO       `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type RiskDocumentDTO struct {
	ID    *uint   `json:"id"`
	URL   *string `json:"url"`
	Icon  string  `json:"icon"`
	Name  string  `json:"name"`
	Token string  `json:"token"`
	Type  string  `json:"type"`
	Size  int     `json:"size"`
}

type RiskRegisterDetailDTO struct {
	ID              int                                  `json:"id"`
	ReferenceNumber string                               `json:"reference_number"`
	RiskCategory    RiskCategoryDTO                      `json:"risk_category"`
	RiskSource      RiskSourceDTO                        `json:"risk_source"`
	OwnerUnit       organizationDTOs.OrganizationUnitDTO `json:"owner_unit"`
	Owner           organizationDTOs.UserMinimalDTO      `json:"owner"`
	RiskStatus      RiskStatusDTO                        `json:"risk_status"`
	Title           string                               `json:"title"`
	Description     *string                              `json:"description"`
	ReportedAt      *time.Time                           `json:"reported_at"`
	ReportedByUser  *organizationDTOs.UserMinimalDTO     `json:"reported_by"`
	ReviewFrequency *masterDTOs.ReviewFrequencyDTO       `json:"review_frequency"`
	Comments        *string                              `json:"comments"`
	Documents       []*RiskDocumentDTO                   `json:"documents"`
	CreatedAt       time.Time                            `json:"created_at"`
	UpdatedAt       *time.Time                           `json:"updated_at"`
	CreatedByUser   organizationDTOs.UserMinimalDTO      `json:"created_by"`
	UpdatedByUser   organizationDTOs.UserMinimalDTO      `json:"updated_by"`
	// Controls               []*RiskControlDTO               `json:"controls"`
	// ComplianceRequirements []*RiskComplianceRequirementDTO `json:"compliance_requirements"`
}

type RiskRegisterFormDTO struct {
	RiskCategoryID    uint                          `json:"risk_category_id" binding:"required"`
	RiskSourceID      uint                          `json:"risk_source_id" binding:"required"`
	OwnerUnitID       uint                          `json:"owner_unit_id" binding:"required"`
	OwnerID           uint                          `json:"owner_id" binding:"required"`
	RiskStatusID      uint                          `json:"risk_status_id" binding:"required"`
	Title             string                        `json:"title" binding:"required,max=255"`
	Description       *string                       `json:"description"`
	ReportedAt        *string                       `json:"reported_at" binding:"required,datetime=2006-01-02"`
	ReportedByID      *uint                         `json:"reported_by" binding:"required"`
	ReviewFrequencyID *uint                         `json:"review_frequency_id"`
	NextReviewDate    *string                       `json:"next_review_date" binding:"omitempty,datetime=2006-01-02"`
	Comments          *string                       `json:"comments"`
	Documents         []*commonDTOs.DocumentFormDTO `json:"documents"`
	// ComplianceRequirements []*uint                       `json:"compliance_requirements"`
	// Controls               []*uint                       `json:"controls"`
}

type RiskRegisterParsedDates struct {
	ReportedAt     *time.Time
	NextReviewDate *time.Time
}
