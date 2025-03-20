package dtos

import commonDTOs "jethur/grc/common/dtos"

type BcpRiskForm struct {
	RiskRegisterIDs []uint `json:"risk_register_ids" binding:"required"`
}

type BcpRiskListDTO struct {
	ID                              uint    `json:"id"`
	RiskRegisterID                  uint    `json:"risk_register_id"`
	ProcessID                       uint    `json:"process_id"`
	ReferenceNumber                 string  `json:"reference_number"`
	RiskCategoryTitleEn             string  `json:"risk_category_title_en"`
	RiskCategoryTitleAr             string  `json:"risk_category_title_ar"`
	RiskSourceTitleEn               string  `json:"risk_source_title_en"`
	RiskSourceTitleAr               string  `json:"risk_source_title_ar"`
	ResponsibleUnitTitleEn          string  `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr          string  `json:"responsible_unit_title_ar"`
	ResponsibleUserID               uint    `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn      string  `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr      string  `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn       string  `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr       string  `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          string  `json:"responsible_user_initial"`
	ResponsibleUserColor            string  `json:"responsible_user_color"`
	ResponsibleUserEmail            string  `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath string  `json:"responsible_user_profile_image_path"`
	RiskStatusTitleEn               string  `json:"risk_status_title_en"`
	RiskStatusTitleAr               string  `json:"risk_status_title_ar"`
	RiskStatusColorCode             string  `json:"risk_status_color_code"`
	Title                           string  `json:"title"`
	Description                     *string `json:"description"`
	RiskMatrixLikelihoodTitleEn     string  `json:"risk_matrix_likelihood_title_en"`
	RiskMatrixLikelihoodTitleAr     string  `json:"risk_matrix_likelihood_title_ar"`
	RiskMatrixLikelihoodColorCode   string  `json:"risk_matrix_likelihood_score"`
	RiskMatrixImpactTitleEn         string  `json:"risk_matrix_impact_title_en"`
	RiskMatrixImpactTitleAr         string  `json:"risk_matrix_impact_title_ar"`
	RiskMatrixImpactColorCode       string  `json:"risk_matrix_impact_score"`
	RiskMatrixRiskRatingTitleEn     string  `json:"risk_matrix_risk_rating_title_en"`
	RiskMatrixRiskRatingTitleAr     string  `json:"risk_matrix_risk_rating_title_ar"`
	RiskMatrixRiskRatingType        string  `json:"risk_matrix_risk_rating_type"`
	RiskMatrixRiskRatingColorCode   string  `json:"risk_matrix_risk_rating_color_code"`
	RiskMatrixRiskRatingValue       string  `json:"risk_matrix_risk_rating_value"`
}

type BcpRiskListResponseDTO struct {
	Items      []*BcpRiskListDTO            `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type BcpRiskListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	Value       int
	QuickFilter []string
}
