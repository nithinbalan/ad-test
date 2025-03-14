package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	complianceDTOs "jethur/grc/compliance/dtos"
	documentsDTOs "jethur/grc/documents/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

// AssetPolicyListDTO - Used for listing Asset Policies
type AssetPolicyListDTO struct {
	ID                        uint      `json:"id"`
	AssetRegisterID           uint      `json:"asset_register_id"`
	AssetRegisterTitleEn      string    `json:"asset_register_title_en"`
	AssetRegisterTitleAr      string    `json:"asset_register_title_ar"`
	AssetCategoryTitleEn      string    `json:"asset_category_title_en"`
	AssetCategoryTitleAr      string    `json:"asset_category_title_ar"`
	AssetSubCategoryTitleEn   string    `json:"asset_sub_category_title_en"`
	AssetSubCategoryTitleAr   string    `json:"asset_sub_category_title_ar"`
	ComplianceStatusTitleEn   string    `json:"compliance_status_title_en"`
	ComplianceStatusTitleAr   string    `json:"compliance_status_title_ar"`
	ComplianceStatusColorCode string    `json:"compliance_status_color_code"`
	DocumentTitleEn           string    `json:"document_title_en"`
	DocumentTitleAr           string    `json:"document_title_ar"`
	DocumentURL               string    `json:"document_url"`
	DocumentName              string    `json:"document_name"`
	DocumentSize              int       `json:"document_size"`
	DocumentToken             string    `json:"document_token"`
	DocumentOriginalName      *string   `json:"document_original_name"`
	CreatedAt                 time.Time `json:"created_at"`
}

// AssetPolicyFormDTO - Used for creating an Asset Policy
type AssetPolicyFormDTO struct {
	DocumentIDs []uint `json:"document_ids" binding:"required"`
}

// AssetPolicyListQueryParamsDTO - Used for filtering and sorting
type AssetPolicyListQueryParamsDTO struct {
	Sort     string
	SortBy   string
	Keyword  string
	Page     int
	PageSize int
	Filters  map[string][]string
}

// AssetPolicyListResponseDTO - Response DTO for listing
type AssetPolicyListResponseDTO struct {
	Items      []*AssetPolicyListDTO        `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type AssetPolicyDetailDTO struct {
	ID                    uint                               `json:"id"`
	AssetRegister         AssetRegisterDetailDTO             `json:"asset_register"`
	ComplianceStatus      complianceDTOs.ComplianceStatusDTO `json:"compliance_status"`
	AssetPolicyChecklists []*AssetPolicyChecklistDTO         `json:"asset_policy_checklists"`
	Document              documentsDTOs.DocumentDetailsDTO   `json:"document"`
	CreatedByUser         organizationDTOs.UserMinimalDTO    `json:"created_by"`
	UpdatedByUser         *organizationDTOs.UserMinimalDTO   `json:"updated_by,omitempty"`
	CreatedAt             time.Time                          `json:"created_at"`
	UpdatedAt             *time.Time                         `json:"updated_at,omitempty"`
}

type AssetPolicyByComplianceDTO struct {
	TotalPolicies          int     `json:"total_policies"`
	Compliant              int     `json:"compliant"`
	NonCompliant           int     `json:"non_compliant"`
	CompliantPercentage    float64 `json:"compliant_percentage"`
	NonCompliantPercentage float64 `json:"non_compliant_percentage"`
}

type AssetPolicyStatusUpdateFormDTO struct {
	ComplianceStatusID uint `json:"compliance_status_id" binding:"required"`
}
