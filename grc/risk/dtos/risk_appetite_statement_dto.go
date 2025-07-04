package dtos

import (
	commonDTOs "jethur/grc/common/dtos"

	"time"
)

type RiskAppetiteStatementListDTO struct {
	ID                          uint                   `json:"id"`
	RiskCategoryID              uint                   `json:"risk_category_id"`
	RiskCategoryTitleEn         string                 `json:"risk_category_title_en"`
	RiskCategoryTitleAr         string                 `json:"risk_category_title_ar"`
	ReferenceNumber             string                 `json:"reference_number"`
	StatementEn                 string                 `json:"statement_en"`
	StatementAr                 string                 `json:"statement_ar"`
	Status                      string                 `json:"status"`
	CreatedAt                   string                 `json:"created_at"`
	UpdatedAt                   string                 `json:"updated_at"`
	CreatedUserID               uint                   `json:"created_user_id"`
	CreatedUserFirstNameEn      string                 `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string                 `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string                 `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       *string                `json:"created_user_last_name_ar"`
	CreatedUserInitial          string                 `json:"created_user_initial"`
	CreatedUserColor            string                 `json:"created_user_color"`
	CreatedUserEmail            string                 `json:"created_user_email"`
	CreatedUserProfileImagePath string                 `json:"created_user_profile_image_path"`
	UpdatedUserID               *uint                  `json:"updated_user_id"`
	UpdatedUserFirstNameEn      *string                `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr      *string                `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn       *string                `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr       *string                `json:"updated_user_last_name_ar"`
	UpdatedUserInitial          string                 `json:"updated_user_initial"`
	UpdatedUserColor            string                 `json:"updated_user_color"`
	UpdatedUserEmail            string                 `json:"updated_user_email"`
	UpdatedUserProfileImagePath string                 `json:"updated_user_profile_image_path"`
	Risks                       []*RiskRegisterListDTO `json:"risks" gorm:"-"`
}

type RiskAppetiteStatementListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
}

type RiskAppetiteStatementListResponseDTO struct {
	Items      []*RiskAppetiteStatementListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO    `json:"pagination"`
}

type RiskAppetiteStatementGroupedDTO struct {
	RiskCategoryID      uint                            `json:"risk_category_id"`
	RiskCategoryTitleEn string                          `json:"risk_category_title_en"`
	RiskCategoryTitleAr string                          `json:"risk_category_title_ar"`
	Statements          []*RiskAppetiteStatementListDTO `json:"statements"`
}

type RiskAppetiteStatementListGroupedResponseDTO struct {
	Items      []*RiskAppetiteStatementGroupedDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO       `json:"pagination"`
}

type RiskAppetiteStatementFormDTO struct {
	RiskCategoryID uint   `json:"risk_category_id" binding:"required"`
	StatementEn    string `json:"statement_en" binding:"required"`
	StatementAr    string `json:"statement_ar" binding:"required"`
	Status         string `json:"status" binding:"required"`
}

type RiskAppetiteStatementDetailsDTO struct {
	ID             uint            `json:"id"`
	RiskCategoryID uint            `json:"risk_category_id"`
	RiskCategory   RiskCategoryDTO `json:"risk_category"`
	StatementEn    string          `json:"statement_en"`
	StatementAr    string          `json:"statement_ar"`
	Status         string          `json:"status"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      *time.Time      `json:"updated_at"`
}

type RiskAppetiteStatementSummaryDTO struct {
	ID                           uint   `json:"id"`
	RiskCategoryTitleEn          string `json:"risk_category_title_en"`
	RiskCategoryTitleAr          string `json:"risk_category_title_ar"`
	RiskAppetiteStatementTitleEn string `json:"risk_appetite_statement_title_en"`
	RiskAppetiteStatementTitleAr string `json:"risk_appetite_statement_title_ar"`
	RisksMapped                  int    `json:"risks_mapped"`
	AboveCount                   int    `json:"above_count"`
	WithinCount                  int    `json:"within_count"`
	BelowCount                   int    `json:"below_count"`
}

type RiskAppetiteStatementSummaryQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Filters     map[string][]string
	QuickFilter []string
}

type RiskAppetiteDashboardSummaryDTO struct {
	ActiveStatementsCount     int `json:"active_statements_count"`
	AlignedRisksCount         int `json:"aligned_risks_count"`
	TotalRisksCount           int `json:"total_risks_count"`
	AboveToleranceRisksCount  int `json:"above_tolerance_risks_count"`
	WithinToleranceRisksCount int `json:"within_tolerance_risks_count"`
	BelowToleranceRisksCount  int `json:"below_tolerance_risks_count"`
}

type MappedRiskListDTO struct {
	RiskID          uint    `json:"risk_id"`
	ReferenceNumber string  `json:"reference_number"`
	Title           string  `json:"title"`
	Description     *string `json:"description"`

	RiskCategoryTitleEn *string `json:"risk_category_title_en"`
	RiskCategoryTitleAr *string `json:"risk_category_title_ar"`

	RiskSourceTitleEn *string `json:"risk_source_title_en"`
	RiskSourceTitleAr *string `json:"risk_source_title_ar"`

	OwnerUnitTitleEn *string `json:"owner_unit_title_en"`
	OwnerUnitTitleAr *string `json:"owner_unit_title_ar"`

	OwnerID               *uint   `json:"owner_id"`
	OwnerFirstNameEn      *string `json:"owner_first_name_en"`
	OwnerFirstNameAr      *string `json:"owner_first_name_ar"`
	OwnerLastNameEn       *string `json:"owner_last_name_en"`
	OwnerLastNameAr       *string `json:"owner_last_name_ar"`
	OwnerInitial          *string `json:"owner_initial"`
	OwnerColor            *string `json:"owner_color"`
	OwnerEmail            *string `json:"owner_email"`
	OwnerProfileImagePath *string `json:"owner_profile_image_path"`

	RiskMatrixLikelihoodTitleEn *string `json:"risk_matrix_likelihood_title_en"`
	RiskMatrixLikelihoodTitleAr *string `json:"risk_matrix_likelihood_title_ar"`
	RiskMatrixLikelihoodScore   *int    `json:"risk_matrix_likelihood_score"`

	RiskMatrixImpactTitleEn *string `json:"risk_matrix_impact_title_en"`
	RiskMatrixImpactTitleAr *string `json:"risk_matrix_impact_title_ar"`
	RiskMatrixImpactScore   *int    `json:"risk_matrix_impact_score"`

	RiskMatrixRiskRatingTitleEn   *string `json:"risk_matrix_risk_rating_title_en"`
	RiskMatrixRiskRatingTitleAr   *string `json:"risk_matrix_risk_rating_title_ar"`
	RiskMatrixRiskRatingType      *string `json:"risk_matrix_risk_rating_type"`
	RiskMatrixRiskRatingColorCode *string `json:"risk_matrix_risk_rating_color_code"`

	RiskStatusTitleEn   *string `json:"risk_status_title_en"`
	RiskStatusTitleAr   *string `json:"risk_status_title_ar"`
	RiskStatusColorCode *string `json:"risk_status_color_code"`

	AppetiteAlignment *string `json:"appetite_alignment"`
}

type MappedRiskListResponseDTO struct {
	Items      []*MappedRiskListDTO         `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type MappedRiskListQueryParamsDTO struct {
	Sort     string `json:"sort"`
	SortBy   string `json:"sort_by"`
	Keyword  string `json:"keyword"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
}
