package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	"time"
)

type RCMListDTO struct {
	ID                       uint             `json:"id"`
	AuditPlanAuditableItemID uint             `json:"audit_plan_auditable_item_id"`
	RiskRegisterID           uint             `json:"risk_register_id"`
	RiskTitle                string           `json:"risk_title"`
	RiskType                 string           `json:"risk_type"`
	CreatedAt                time.Time        `json:"created_at"`
	UpdatedAt                *time.Time       `json:"updated_at,omitempty"`
	Controls                 []*RCMControlDTO `json:"controls"`
}

type RCMControlDTO struct {
	ID                                   uint    `json:"id"`
	ControlID                            uint    `json:"control_id"`
	ControlTitleEn                       string  `json:"control_title_en"`
	ControlTitleAr                       string  `json:"control_title_ar"`
	ControlType                          string  `json:"control_type"`
	ControlImplementationStatusTitleEn   *string `json:"control_implementation_status_title_en"`
	ControlImplementationStatusTitleAr   *string `json:"control_implementation_status_title_ar"`
	ControlImplementationStatusColorCode *string `json:"control_implementation_status_color_code"`
	ControlEfficiencyStatusTitleEn       *string `json:"control_efficiency_status_title_en"`
	ControlEfficiencyStatusTitleAr       *string `json:"control_efficiency_status_title_ar"`
	ControlEfficiencyStatusColorCode     *string `json:"control_efficiency_status_color_code"`
}

type RCMListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
}

type RCMListResponseDTO struct {
	Items      []*RCMListDTO                `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}
