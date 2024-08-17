package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	complianceDTOs "jethur/grc/compliance/dtos"
	documentDTOs "jethur/grc/documents/dtos"
	incidentDTOs "jethur/grc/incident/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	riskDTOs "jethur/grc/risk/dtos"
	settingsDTOs "jethur/grc/settings/dtos"
	"time"
)

type AuditableItemListDTO struct {
	ID                          uint   `json:"id"`
	AuditorRiskRatingID         uint   `json:"auditor_risk_rating_id"`
	AuditStatusID               uint   `json:"audit_status_id"`
	TitleEn                     string `json:"title_en"`
	TitleAr                     string `json:"title_ar"`
	Comment                     string `json:"comment"`
	ProcessTitleEn              string `json:"process_title_en"`
	ProcessTitleAr              string `json:"process_title_ar"`
	ControlTitleEn              string `json:"control_title_en"`
	ControlTitleAr              string `json:"control_title_ar"`
	DocumentTitleEn             string `json:"document_title_en"`
	DocumentTitleAr             string `json:"document_title_ar"`
	ResponsibleUnitTitleEn      string `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr      string `json:"responsible_unit_title_ar"`
	IncidentRegisterTitle       string `json:"incident_register_title"`
	ComplianceRequirementTitle  string `json:"compliance_requirement_title"`
	RiskRegisterTitle           string `json:"risk_register_title"`
	UserID                      uint   `json:"user_id"`
	UserFirstNameEn             string `json:"user_first_name_en"`
	UserFirstNameAr             string `json:"user_first_name_ar"`
	UserLastNameEn              string `json:"user_last_name_en"`
	UserLastNameAr              string `json:"user_last_name_ar"`
	UserEmail                   string `json:"user_email"`
	UserProfileImagePath        string `json:"user_profile_image_path"`
	CreatedUserID               uint   `json:"created_user_id"`
	CreatedUserFirstNameEn      string `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       string `json:"created_user_last_name_ar"`
	CreatedUserEmail            string `json:"created_user_email"`
	CreatedUserProfileImagePath string `json:"created_user_profile_image_path"`
	UpdatedUserID               *uint  `json:"updated_user_id"`
	UpdatedUserFirstNameEn      string `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr      string `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn       string `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr       string `json:"updated_user_last_name_ar"`
	UpdatedUserEmail            string `json:"updated_user_email"`
	UpdatedUserProfileImagePath string `json:"updated_user_profile_image_path"`
	AuditStatusTitleEn          string `json:"audit_status_title_en"`
	AuditStatusTitleAr          string `json:"audit_status_title_ar"`
	AuditStatusColorCode        string `json:"audit_status_color_code"`
	AuditStatusType             string `json:"audit_status_type"`
	AuditorRiskRatingTitleEn    string `json:"auditor_risk_rating_title_en"`
	AuditorRiskRatingTitleAr    string `json:"auditor_risk_rating_title_ar"`
	AuditorRiskRatingColorCode  string `json:"auditor_risk_rating_color_code"`
	AuditorRiskRatingType       string `json:"auditor_risk_rating_type"`
	CreatedAt                   string `json:"created_at"`
	UpdatedAt                   string `json:"updated_at"`
}
type AuditableItemListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
}

type AuditableItemListResponseDTO struct {
	Items      []*AuditableItemListDTO      `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type AuditableItemFormDTO struct {
	TitleEn             string `json:"title_en" binding:"required,max=500"`
	TitleAr             string `json:"title_ar" binding:"required,max=500"`
	ResponsibleUnitID   uint   `json:"responsible_unit_id" binding:"required"`
	Comment             string `json:"comment" binding:"required"`
	AuditorRiskRatingID uint   `json:"auditor_risk_rating_id" binding:"required"`
}

type AuditableItemDetailDTO struct {
	ID                    uint                                           `json:"id"`
	TitleEn               string                                         `json:"title_en"`
	TitleAr               string                                         `json:"title_ar"`
	Comment               *string                                        `json:"comment"`
	AuditStatus           *AuditStatusDTO                                `json:"audit_status"`
	AuditorRiskRating     *AuditorRiskRatingDTO                          `json:"auditor_risk_rating"`
	Process               *organizationDTOs.ProcessDetailDTO             `json:"process"`
	Control               *organizationDTOs.ControlDetailDTO             `json:"control"`
	Document              *documentDTOs.DocumentDetailsDTO               `json:"document"`
	IncidentRegister      *incidentDTOs.IncidentRegisterDetailDTO        `json:"incident_register"`
	ComplianceRequirement *complianceDTOs.ComplianceRequirementDetailDTO `json:"compliance_requirement"`
	RiskRegister          *riskDTOs.RiskRegisterDetailDTO                `json:"risk_register"`
	User                  *organizationDTOs.UserMinimalDTO               `json:"user"`
	ResponsibleUnit       settingsDTOs.OrganizationUnitDTO               `json:"responsible_unit"`
	CreatedAt             time.Time                                      `json:"created_at"`
	UpdatedAt             *time.Time                                     `json:"updated_at"`
	CreatedByUser         *organizationDTOs.UserMinimalDTO               `json:"created_by"`
	UpdatedByUser         *organizationDTOs.UserMinimalDTO               `json:"updated_by"`
}

type AuditableItemStatusUpdateFormDTO struct {
	AuditStatusID uint `json:"audit_status_id"`
}
