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
