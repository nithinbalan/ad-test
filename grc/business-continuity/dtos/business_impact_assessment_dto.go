package dtos

type BusinessImpactAssessmentListDTO struct {
	ID                          uint    `json:"id"`
	ProcessID                   *uint   `json:"process_id"`
	ProcessTitleEn              string  `json:"process_title_en"`
	ProcessTitleAr              string  `json:"process_title_ar"`
	ImpactCategoryID            uint    `json:"impact_category_id"`
	ImpactCategoryTitleEn       string  `json:"impact_category_title_en"`
	ImpactCategoryTitleAr       string  `json:"impact_category_title_ar"`
	ImpactCategoryType          string  `json:"impact_category_type"`
	RiskLevelID                 *uint   `json:"risk_level_id"`
	RiskLevelTitleEn            string  `json:"risk_level_title_en"`
	RiskLevelTitleAr            string  `json:"risk_level_title_ar"`
	RiskLevelColorCode          string  `json:"risk_level_color_code"`
	RiskLevelType               string  `json:"risk_level_type"`
	Cost                        *uint   `json:"cost"`
	Comments                    *string `json:"comments"`
	IsApplicable                bool    `json:"is_applicable"`
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
	ID             uint    `json:"id"`
	ProcessID      uint    `json:"process_id"`
	ProcessTitleEn *string `json:"process_title_en"`
	ProcessTitleAr *string `json:"process_title_ar"`

	ProcessDependencySubTypeID      uint    `json:"process_dependency_sub_type_id"`
	ProcessDependencySubTypeType    *string `json:"process_dependency_sub_type_type"`
	ProcessDependencySubTypeTitleEn *string `json:"process_dependency_sub_type_title_en"`
	ProcessDependencySubTypeTitleAr *string `json:"process_dependency_sub_type_title_ar"`
	ProcessDependencyTypeTitleEn    *string `json:"process_dependency_type_title_en"`
	ProcessDependencyTypeTitleAr    *string `json:"process_dependency_type_title_ar"`
	ProcessDependencyTypeType       *string `json:"process_dependency_type_type"`

	DependentProcessID      *uint   `json:"dependent_process_id"`
	DependentProcessTitleEn *string `json:"dependent_process_title_en"`
	DependentProcessTitleAr *string `json:"dependent_process_title_ar"`

	OrganizationUnitID      *uint   `json:"organization_unit_id"`
	OrganizationUnitTitleEn *string `json:"organization_unit_title_en"`
	OrganizationUnitTitleAr *string `json:"organization_unit_title_ar"`

	AssetID      *uint   `json:"asset_id"`
	AssetTitleEn *string `json:"asset_title_en"`
	AssetTitleAr *string `json:"asset_title_ar"`

	Title       string  `json:"title"`
	Description *string `json:"description"`

	Rto               *int    `json:"rto"`
	RtoUnitId         *uint   `json:"rto_unit_id"`
	RtoUnitTitleEn    *string `json:"rto_unit_title_en"`
	RtoUnitTitleAr    *string `json:"rto_unit_title_ar"`
	Mtpd              *int    `json:"mtpd"`
	MtpdUnitId        *uint   `json:"mtpd_unit_id"`
	MtpdUnitTitleEn   *string `json:"mtpd_unit_title_en"`
	MtpdUnitTitleAr   *string `json:"mtpd_unit_title_ar"`
	Rpo               *int    `json:"rpo"`
	RpoUnitId         *uint   `json:"rpo_unit_id"`
	RpoUnitTitleEn    *string `json:"rpo_unit_title_en"`
	RpoUnitTitleAr    *string `json:"rpo_unit_title_ar"`
	ImpactDescription *string `json:"impact_description"`

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

type BusinessImpactAssessmentListQueryParamsDTO struct {
	Sort    string
	SortBy  string
	Keyword string
	Filters map[string][]string
}

type BusinessImpactAssessmentRecoveryObjectiveListQueryParamsDTO struct {
	Sort    string
	SortBy  string
	Keyword string
	Filters map[string][]string
}

type BusinessImpactFormDTO struct {
	ProcessID *uint `json:"process_id"`
}

type BusinessImpactAssessmentFormDTO struct {
	ProcessID                 *uint                         `json:"process_id"`
	BusinessImpactAssessments []BusinessImpactAssessmentDTO `json:"business_impact_assessments"`
	RiskLevelID               *uint                         `json:"risk_level_id"`
	ImpactJustification       *string                       `json:"impact_justification"`
}

type BusinessImpactAssessmentDTO struct {
	ImpactCategoryID uint    `json:"impact_category_id"`
	RiskLevelID      *uint   `json:"risk_level_id"`
	Cost             *uint   `json:"cost"`
	Comments         *string `json:"comments"`
	IsApplicable     bool    `json:"is_applicable"`
}

type BusinessImpactAssessmentRecoveryObjectiveFormDTO struct {
	ProcessID          *uint                                          `json:"process_id"`
	Rto                *int                                           `json:"rto"`
	RtoUnitId          *uint                                          `json:"rto_unit_id"`
	Mtpd               *int                                           `json:"mtpd"`
	MtpdUnitId         *uint                                          `json:"mtpd_unit_id"`
	Rpo                *int                                           `json:"rpo"`
	RpoUnitId          *uint                                          `json:"rpo_unit_id"`
	RecoveryObjectives []BusinessImpactAssessmentRecoveryObjectiveDTO `json:"recovery_objectives"`
}

type BusinessImpactAssessmentProcessCriticalityFormDTO struct {
	ProcessID                 *uint   `json:"process_id"`
	ProcessCriticalityLevelID *uint   `json:"process_criticality_level_id"`
	CriticalityJustification  *string `json:"criticality_justification"`
}

type BusinessImpactAssessmentRecoveryObjectiveDTO struct {
	ProcessDependencyID *uint `json:"process_dependency_id"`
	Rto                 *int  `json:"rto"`
	RtoUnitId           *uint `json:"rto_unit_id"`
	Mtpd                *int  `json:"mtpd"`
	MtpdUnitId          *uint `json:"mtpd_unit_id"`
	Rpo                 *int  `json:"rpo"`
	RpoUnitId           *uint `json:"rpo_unit_id"`
}

type BIAMatrixResponseDTO struct {
	BusinessFunctionID                 uint    `json:"business_function_id"`
	BusinessFunctionTitle              string  `json:"business_function_title"`
	BusinessFunctionCriticalityTitleEn *string `json:"business_function_criticality_title_en"`
	BusinessFunctionCriticalityTitleAr *string `json:"business_function_criticality_title_ar"`
	BusinessFunctionCriticalityColor   *string `json:"business_function_criticality_color"`

	ProcessID                 uint    `json:"process_id"`
	ProcessTitleEn            string  `json:"process_title_en"`
	ProcessTitleAr            string  `json:"process_title_ar"`
	ProcessCriticalityTitleEn *string `json:"process_criticality_title_en"`
	ProcessCriticalityTitleAr *string `json:"process_criticality_title_ar"`
	ProcessCriticalityColor   *string `json:"process_criticality_color"`

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

	RiskLevelID        *uint   `json:"risk_level_id"`
	RiskLevelTitleEn   *string `json:"risk_level_title_en"`
	RiskLevelTitleAr   *string `json:"risk_level_title_ar"`
	RiskLevelColorCode *string `json:"risk_level_color_code"`
}

type BIAMatrixQueryParamsDTO struct {
	Sort     string
	SortBy   string
	Keyword  string
	Page     int
	PageSize int
	Filters  map[string][]string
	Export   bool
}
