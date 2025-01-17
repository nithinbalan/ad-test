package dtos

type BusinessImpactAssessmentListDTO struct {
	ID                          uint    `json:"id"`
	ProcessID                   *uint   `json:"process_id"`
	ProcessTitleEn              string  `json:"process_title_en"`
	ProcessTitleAr              string  `json:"process_title_ar"`
	ImpactCategoryID            uint    `json:"impact_category_id"`
	ImpactCategoryTitleEn       string  `json:"impact_category_title_en"`
	ImpactCategoryTitleAr       string  `json:"impact_category_title_ar"`
	RiskLevelID                 *uint   `json:"risk_level_id"`
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
	UpdatedUserFirstNameEn      *string `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr      *string `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn       *string `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr       *string `json:"updated_user_last_name_ar"`
	UpdatedUserInitial          *string `json:"updated_user_initial"`
	UpdatedUserColor            *string `json:"updated_user_color"`
	UpdatedUserEmail            *string `json:"updated_user_email"`
	UpdatedUserProfileImagePath *string `json:"updated_user_profile_image_path"`
}

type BusinessImpactAssessmentRecoveryObjectiveListDTO struct {
	ID                           uint    `json:"id"`
	ProcessID                    *uint   `json:"process_id"`
	ProcessTitleEn               string  `json:"process_title_en"`
	ProcessTitleAr               string  `json:"process_title_ar"`
	AssetRegisterID              *uint   `json:"asset_register_id"`
	AssetRegisterTitleEn         *string `json:"asset_register_title_en"`
	AssetRegisterTitleAr         *string `json:"asset_register_title_ar"`
	AssetRegisterRto             *uint   `json:"asset_register_rto"`
	AssetRegisterRtoUnitId       *uint   `json:"asset_register_rto_unit_id"`
	AssetRegisterRtoUnitTitleEn  *string `json:"asset_register_rto_unit_title_en"`
	AssetRegisterRtoUnitTitleAr  *string `json:"asset_register_rto_unit_title_ar"`
	AssetRegisterMtpd            *uint   `json:"asset_register_mtpd"`
	AssetRegisterMtpdUnitId      *uint   `json:"asset_register_mtpd_unit_id"`
	AssetRegisterMtpdUnitTitleEn *string `json:"asset_register_mtpd_unit_title_en"`
	AssetRegisterMtpdUnitTitleAr *string `json:"asset_register_mtpd_unit_title_ar"`
	AssetRegisterRpo             *uint   `json:"asset_register_rpo"`
	AssetRegisterRpoUnitId       *uint   `json:"asset_register_rpo_unit_id"`
	AssetRegisterRpoUnitTitleEn  *string `json:"asset_register_rpo_unit_title_en"`
	AssetRegisterRpoUnitTitleAr  *string `json:"asset_register_rpo_unit_title_ar"`
	DocumentID                   *uint   `json:"document_id"`
	DocumentTitleEn              *string `json:"document_title_en"`
	DocumentTitleAr              *string `json:"document_title_ar"`
	DocumentRto                  *uint   `json:"document_rto"`
	DocumentRtoUnitId            *uint   `json:"document_rto_unit_id"`
	DocumentRtoUnitTitleEn       *string `json:"document_rto_unit_title_en"`
	DocumentRtoUnitTitleAr       *string `json:"document_rto_unit_title_ar"`
	DocumentMtpd                 *uint   `json:"document_mtpd"`
	DocumentMtpdUnitId           *uint   `json:"document_mtpd_unit_id"`
	DocumentMtpdUnitTitleEn      *string `json:"document_mtpd_unit_title_en"`
	DocumentMtpdUnitTitleAr      *string `json:"document_mtpd_unit_title_ar"`
	DocumentRpo                  *uint   `json:"document_rpo"`
	DocumentRpoUnitId            *uint   `json:"document_rpo_unit_id"`
	DocumentRpoUnitTitleEn       *string `json:"document_rpo_unit_title_en"`
	DocumentRpoUnitTitleAr       *string `json:"document_rpo_unit_title_ar"`
	ImpactDescription            *string `json:"impact_description"`
	CreatedAt                    string  `json:"created_at"`
	UpdatedAt                    string  `json:"updated_at"`
	CreatedUserID                uint    `json:"created_user_id"`
	CreatedUserFirstNameEn       string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr       string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn        string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr        string  `json:"created_user_last_name_ar"`
	CreatedUserInitial           string  `json:"created_user_initial"`
	CreatedUserColor             string  `json:"created_user_color"`
	CreatedUserEmail             string  `json:"created_user_email"`
	CreatedUserProfileImagePath  string  `json:"created_user_profile_image_path"`
	UpdatedUserID                *uint   `json:"updated_user_id"`
	UpdatedUserFirstNameEn       *string `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr       *string `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn        *string `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr        *string `json:"updated_user_last_name_ar"`
	UpdatedUserInitial           *string `json:"updated_user_initial"`
	UpdatedUserColor             *string `json:"updated_user_color"`
	UpdatedUserEmail             *string `json:"updated_user_email"`
	UpdatedUserProfileImagePath  *string `json:"updated_user_profile_image_path"`
}

type BusinessImpactAssessmentListQueryParamsDTO struct {
	Sort    string
	SortBy  string
	Keyword string
	Filters map[string][]string
}

// type BusinessImpactAssessmentListResponseDTO struct {
// 	Items []*BusinessImpactAssessmentListDTO `json:"items"`
// }

type BusinessImpactAssessmentRecoveryObjectiveListQueryParamsDTO struct {
	Sort    string
	SortBy  string
	Keyword string
	Filters map[string][]string
}

// type BusinessImpactAssessmentRecoveryObjectiveListResponseDTO struct {
// 	Items []*BusinessImpactAssessmentRecoveryObjectiveListDTO `json:"items"`
// }

type BusinessImpactFormDTO struct {
	ProcessID *uint `json:"process_id"`
}

type BusinessImpactAssessmentFormDTO struct {
	ProcessID                 *uint                         `json:"process_id"`
	BusinessImpactAssessments []BusinessImpactAssessmentDTO `json:"business_impact_assessments"`
}

type BusinessImpactAssessmentDTO struct {
	ImpactCategoryID uint    `json:"impact_category_id"`
	RiskLevelID      *uint   `json:"risk_level_id"`
	Cost             *uint   `json:"cost"`
	Comments         *string `json:"comments"`
}

type BusinessImpactAssessmentRecoveryObjectiveFormDTO struct {
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

type BIAMatrixResponseDTO struct {
	BusinessFunctionID     uint    `json:"business_function_id"`
	BusinessFunctionTitle  string  `json:"business_function_title"`
	ProcessID              uint    `json:"process_id"`
	ProcessTitleEn         string  `json:"process_title_en"`
	ProcessTitleAr         string  `json:"process_title_ar"`
	ProcessRTO             int     `json:"process_rto"`
	ProcessMTPD            int     `json:"process_mtpd"`
	ProcessRPO             int     `json:"process_rpo"`
	ProcessRTOUnitTitleEn  *string `json:"process_rto_unit_title_en"`
	ProcessRTOUnitTitleAr  *string `json:"process_rto_unit_title_ar"`
	ProcessMTPDUnitID      *uint   `json:"process_mtpd_unit_id"`
	ProcessMTPDUnitTitleEn *string `json:"process_mtpd_unit_title_en"`
	ProcessMTPDUnitTitleAr *string `json:"process_mtpd_unit_title_ar"`
	ProcessRPOUnitID       *uint   `json:"process_rpo_unit_id"`
	ProcessRPOUnitTitleEn  *string `json:"process_rpo_unit_title_en"`
	ProcessRPOUnitTitleAr  *string `json:"process_rpo_unit_title_ar"`
	RiskLevelID            *uint   `json:"risk_level_id"`
	RiskLevelTitleEn       *string `json:"risk_level_title_en"`
	RiskLevelTitleAr       *string `json:"risk_level_title_ar"`
}

type BIAMatrixQueryParamsDTO struct {
	Sort     string
	SortBy   string
	Keyword  string
	Page     int
	PageSize int
	Filters  map[string][]string
}

type BIAMatrixResponsePaginatedDTO struct {
	Items      []BIAMatrixResponseDTO       `json:"items"`
}
