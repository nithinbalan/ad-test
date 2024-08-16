package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	"time"
)

type RiskMatrixImpactGuidelineListDTO struct {
	ID                        uint    `json:"id"`
	RiskImpactCategoryTitleEn string  `json:"risk_impact_category_title_en"`
	RiskImpactCategoryTitleAr string  `json:"risk_impact_category_title_ar"`
	RiskMatrixImpactTitleEn   string  `json:"risk_matrix_impact_title_en"`
	RiskMatrixImpactTitleAr   string  `json:"risk_matrix_impact_title_ar"`
	Guideline                 string  `json:"guideline"`
	CreatedAt                 string  `json:"created_at"`
	UpdatedAt                 string  `json:"updated_at"`
	CreatedUserID             uint    `json:"created_user_id"`
	CreatedUserFirstNameEn    string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr    string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn     string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr     *string `json:"created_user_last_name_ar"`
	UpdatedUserID             *uint   `json:"updated_user_id"`
	UpdatedUserFirstNameEn    *string `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr    *string `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn     *string `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr     *string `json:"updated_user_last_name_ar"`
}

type RiskMatrixImpactGuidelineListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
}

type RiskMatrixImpactGuidelineListResponseDTO struct {
	Items      []*RiskMatrixImpactGuidelineListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO        `json:"pagination"`
}

type RiskMatrixImpactGuidelineFormDTO struct {
	RiskImpactCategoryID uint     `json:"risk_impact_category_id" binding:"required"`
	RiskMatrixImpactID   uint     `json:"risk_matrix_impact_id" binding:"required"`
	Guidelines           []string `json:"guidelines" binding:"required,dive,max=500" gorm:"type:varchar(500);not null"`
	Status               string   `json:"status" binding:"required"`
}

type RiskMatrixImpactGuidelineUpdateFormDTO struct {
	RiskImpactCategoryID uint   `json:"risk_impact_category_id" binding:"required"`
	RiskMatrixImpactID   uint   `json:"risk_matrix_impact_id" binding:"required"`
	Guideline            string `json:"guideline" binding:"required,max=500" gorm:"type:varchar(500);not null"`
}

type OrganizationUnitDTO struct {
	ID      int    `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Status  string `json:"status"`
}

type UserMinimalDTO struct {
	ID                int                 `json:"id"`
	FirstNameEN       string              `json:"first_name_en"`
	LastNameEN        string              `json:"last_name_en"`
	FirstNameAR       string              `json:"first_name_ar"`
	LastNameAR        string              `json:"last_name_ar"`
	OrganizationUnit1 OrganizationUnitDTO `json:"organization_unit_1"`
}

type RiskMatrixImpactGuidelineDetailsDTO struct {
	ID                 uint                  `json:"id"`
	RiskImpactCategory RiskImpactCategoryDTO `json:"risk_impact_category"`
	RiskMatrixImpact   RiskMatrixImpactDTO   `json:"risk_matrix_impact"`
	Guideline          string                `json:"guideline"`
	Status             string                `json:"status"`
	CreatedAt          time.Time             `json:"created_at"`
	UpdatedAt          *time.Time            `json:"updated_at"`
	CreatedByUser      UserMinimalDTO        `json:"created_by"`
	UpdatedByUser      UserMinimalDTO        `json:"updated_by"`
}
