package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type RiskComplianceRequirementListDTO struct {
	ComplianceRegisterID                 uint   `json:"compliance_register_id"`
	ComplianceRegisterTitleEn            string `json:"compliance_register_title_en"`
	ComplianceRegisterTitleAr            string `json:"compliance_register_title_ar"`
	ComplianceRequirementID              uint   `json:"compliance_requirement_id"`
	ComplianceClauseTitleEn              string `json:"compliance_clause_title_en"`
	ComplianceClauseTitleAr              string `json:"compliance_clause_title_ar"`
	ComplianceClauseNo                   string `json:"compliance_clause_no"`
	ComplianceRequirementCategoryTitleEn string `json:"compliance_requirement_category_title_en"`
	ComplianceRequirementCategoryTitleAr string `json:"compliance_requirement_category_title_ar"`
	ComplianceRequirementTitleEn         string `json:"compliance_requirement_title_en"`
	ComplianceRequirementTitleAr         string `json:"compliance_requirement_title_ar"`
	ComplianceStatusTitleEn              string `json:"compliance_status_title_en"`
	ComplianceStatusTitleAr              string `json:"compliance_status_title_ar"`
	ComplianceStatusColorCode            string `json:"compliance_status_color_code"`
}

type RiskStrategicObjectiveListDTO struct {
	StrategicObjectiveID      uint    `json:"strategic_objective_id"`
	StrategicObjectiveTitleEn string  `json:"strategic_objective_title_en"`
	StrategicObjectiveTitleAr *string `json:"strategic_objective_title_ar"`
	StrategicObjectiveStatus  *string `json:"strategic_objective_status"`
}
type RiskSLAListDTO struct {
	ReferenceNumber                 string `json:"reference_number"`
	LifecycleStatusID               uint   `json:"lifecycle_status_id"`
	LifecycleStatusTitleEn          string `json:"lifecycle_status_title_en"`
	LifecycleStatusTitleAr          string `json:"lifecycle_status_title_ar"`
	LifecycleStatusColorCode        string `json:"lifecycle_status_color_code"`
	TitleEn                         string `json:"title_en"`
	TitleAr                         string `json:"title_ar"`
	Description                     string `json:"description"`
	Client                          string `json:"client"`
	SLACategoryTitleEn              string `json:"sla_category_title_en"`
	SLACategoryTitleAr              string `json:"sla_category_title_ar"`
	SeverityLevelTitleEn            string `json:"severity_level_title_en"`
	SeverityLevelTitleAr            string `json:"severity_level_title_ar"`
	ResponsibleUnitTitleEn          string `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr          string `json:"responsible_unit_title_ar"`
	ResponsibleUserID               uint   `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn      string `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr      string `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn       string `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr       string `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          string `json:"responsible_user_initial"`
	ResponsibleUserColor            string `json:"responsible_user_color"`
	ResponsibleUserEmail            string `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath string `json:"responsible_user_profile_image_path"`
	ComplianceStatusTitleEn         string `json:"compliance_status_title_en"`
	ComplianceStatusTitleAr         string `json:"compliance_status_title_ar"`
	ComplianceStatusColorCode       string `json:"compliance_status_color_code"`
}

type RiskIncidentListDTO struct {
	ReferenceNumber            string     `json:"reference_number"`
	Title                      string     `json:"title"`
	Description                *string    `json:"description"`
	IncidentAt                 time.Time  `json:"incident_at"`
	Location                   *string    `json:"location"`
	ReportedAt                 *time.Time `json:"reported_at"`
	QuickAction                *string    `json:"quick_action"`
	IncidentCategoryTitleEn    string     `json:"incident_category_title_en"`
	IncidentCategoryTitleAr    string     `json:"incident_category_title_ar"`
	IncidentSubCategoryTitleEn string     `json:"incident_sub_category_title_en"`
	IncidentSubCategoryTitleAr string     `json:"incident_sub_category_title_ar"`
	IncidentTypeTitleEn        string     `json:"incident_type_title_en"`
	IncidentTypeTitleAr        string     `json:"incident_type_title_ar"`
	SeverityLevelTitleEn       string     `json:"severity_level_title_en"`
	SeverityLevelTitleAr       string     `json:"severity_level_title_ar"`
	SeverityLevelColorCode     string     `json:"severity_level_color_code"`
	IncidentStatusTitleEn      string     `json:"incident_status_title_en"`
	IncidentStatusTitleAr      string     `json:"incident_status_title_ar"`
	IncidentStatusColorCode    string     `json:"incident_status_color_code"`
}

type RiskAssetListDTO struct {
	ID                              uint       `json:"id"`
	ReferenceNumber                 string     `json:"reference_number"`
	TitleEn                         string     `json:"title_en"`
	TitleAr                         string     `json:"title_ar"`
	AssetTypeTitleEn                *string    `json:"asset_type_title_en"`
	AssetTypeTitleAr                *string    `json:"asset_type_title_ar"`
	AssetCategoryTitleEn            *string    `json:"asset_category_title_en"`
	AssetCategoryTitleAr            *string    `json:"asset_category_title_ar"`
	AssetStatusTitleEn              string     `json:"asset_status_title_en"`
	AssetStatusTitleAr              string     `json:"asset_status_title_ar"`
	AssetStatusColorCode            string     `json:"asset_status_color_code"`
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
	Location                        *string    `json:"location"`
	Description                     *string    `json:"description"`
	AcquisitionDate                 string     `json:"acquisition_date"`
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
	ResponsibleUnitTitleEn          *string    `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr          *string    `json:"responsible_unit_title_ar"`
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
	UpdatedUserInitial              *string    `json:"updated_user_initial"`
	UpdatedUserColor                *string    `json:"updated_user_color"`
	UpdatedUserEmail                *string    `json:"updated_user_email"`
	UpdatedUserProfileImagePath     *string    `json:"updated_user_profile_image_path"`
	UpdatedAt                       *time.Time `json:"updated_at"`
}

type RiskAssociatedItemListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	QuickFilter []string
}

type RiskAssociatedItemControlListResponseDTO struct {
	Items      []*commonDTOs.AssociatedControlListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO           `json:"pagination"`
}

type RiskAssociatedItemComplianceRequirementListResponseDTO struct {
	Items      []*organizationDTOs.ComplianceRegisterDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO              `json:"pagination"`
}

type RiskAssociatedItemStrategicObjectiveListResponseDTO struct {
	Items      []*RiskStrategicObjectiveListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO     `json:"pagination"`
}

type RiskAssociatedItemSLAListResponseDTO struct {
	Items      []*RiskSLAListDTO            `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type RiskAssociatedItemIncidentListResponseDTO struct {
	Items      []*RiskIncidentListDTO       `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type RiskAssociatedItemAssetListResponseDTO struct {
	Items      []*RiskAssetListDTO          `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type RiskAssociatedItemThreatListResponseDTO struct {
	Items      []*ThreatListDTO             `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}
