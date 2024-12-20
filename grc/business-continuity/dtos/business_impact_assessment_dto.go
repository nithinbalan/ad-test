package dtos

import commonDTOs "jethur/grc/common/dtos"

type BusinessImpactAssessmentListDTO struct {
	ID                          uint    `json:"id"`
	ProcessID                   *uint   `json:"process_id"`
	ProcessTitleEn              string  `json:"process_title_en"`
	ProcessTitleAr              string  `json:"process_title_ar"`
	AssetRegisterID             *uint   `json:"asset_register_id"`
	AssetRegisterTitleEn        string  `json:"asset_register_title_en"`
	AssetRegisterTitleAr        string  `json:"asset_register_title_ar"`
	ImpactCategoryID            uint    `json:"impact_category_id"`
	ImpactCategoryTitleEn       string  `json:"impact_category_title_en"`
	ImpactCategoryTitleAr       string  `json:"impact_category_title_ar"`
	RiskLevelID                 uint    `json:"risk_level_id"`
	RiskLevelTitleEn            string  `json:"risk_level_title_en"`
	RiskLevelTitleAr            string  `json:"risk_level_title_ar"`
	RiskLevelColorCode          string  `json:"risk_level_color_code"`
	Cost                        *uint   `json:"cost"`
	Comments                    *string `json:"comments"`
	CreatedAt                   string  `json:"created_at"`
	UpdatedAt                   string  `json:"updated_at"`
	CreatedUserID               uint    `json:"created_user_id"`
	CreatedUserFirstNameEn      string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       string  `json:"created_user_last_name_ar"`
	CreatedUserInitial          string  `json:"created_user_initial"`
	CreatedUserColor            string  `json:"created_user_color"`
	CreatedUserEmail            string  `json:"created_user_email"`
	CreatedUserProfileImagePath string  `json:"created_user_profile_image_path"`
	UpdatedUserID               *uint   `json:"updated_user_id"`
	UpdatedUserFirstNameEn      string  `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr      string  `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn       string  `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr       string  `json:"updated_user_last_name_ar"`
	UpdatedUserInitial          string  `json:"updated_user_initial"`
	UpdatedUserColor            string  `json:"updated_user_color"`
	UpdatedUserEmail            string  `json:"updated_user_email"`
	UpdatedUserProfileImagePath string  `json:"updated_user_profile_image_path"`
}

type BusinessImpactAssessmentListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
	Export      bool
}

type BusinessImpactAssessmentListResponseDTO struct {
	Items      []*BusinessImpactAssessmentListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO       `json:"pagination"`
}

type BusinessImpactAssessmentFormDTO struct {
	Type                      string                        `json:"type"`
	ProcessID                 *uint                         `json:"process_id"`
	AssetRegisterID           *uint                         `json:"asset_register_id"`
	BusinessImpactAssessments []BusinessImpactAssessmentDTO `json:"business_impact_assessments"`
}

type BusinessImpactAssessmentDTO struct {
	ImpactCategoryID uint    `json:"impact_category_id"`
	RiskLevelID      uint    `json:"risk_level_id"`
	Cost             *uint   `json:"cost"`
	Comments         *string `json:"comments"`
}

type BusinessImpactAssessmentRecoveryObjectiveFormDTO struct {
	Type               string                                         `json:"type"`
	ProcessID          *uint                                          `json:"process_id"`
	Rto                *uint                                          `json:"rto"`
	RtoUnitId          *uint                                          `json:"rto_unit_id"`
	Mtpd               *uint                                          `json:"mtpd"`
	MtpdUnitId         *uint                                          `json:"mtpd_unit_id"`
	Rpo                *uint                                          `json:"rpo"`
	RpoUnitId          *uint                                          `json:"rpo_unit_id"`
	RecoveryObjectives []BusinessImpactAssessmentRecoveryObjectiveDTO `json:"recovery_objectives"`
}

type BusinessImpactAssessmentRecoveryObjectiveDTO struct {
	AssetRegisterID   *uint   `json:"asset_register_id"`
	DocumentID        *uint   `json:"document_id"`
	RoleInProcess     *string `json:"role_in_process" gorm:"type:nvarchar(max)"`
	Rto               *uint   `json:"rto"`
	RtoUnitId         *uint   `json:"rto_unit_id"`
	Mtpd              *uint   `json:"mtpd"`
	MtpdUnitId        *uint   `json:"mtpd_unit_id"`
	Rpo               *uint   `json:"rpo"`
	RpoUnitId         *uint   `json:"rpo_unit_id"`
	ImpactDescription *string `json:"impact_description"`
}
