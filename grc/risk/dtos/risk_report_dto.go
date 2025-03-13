package dtos

import "time"

// RiskReportDetailDTO represents the complete risk report details
type RiskReportDetailDTO struct {
	ID                     uint                   `json:"id"`
	Title                  string                 `json:"title"`
	StartDate              time.Time              `json:"start_date"`
	EndDate                time.Time              `json:"end_date"`
	CreatedUserFirstNameEn string                 `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr string                 `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn  string                 `json:"created_user_last_name_en"`
	CreatedUserLastNameAr  string                 `json:"created_user_last_name_ar"`
	CreatedAt              time.Time              `json:"created_at"`
	Sections               []RiskReportSectionDTO `json:"sections"`
}

// RiskReportSectionDTO represents a section inside the risk report
type RiskReportSectionDTO struct {
	Type            string                  `json:"type"`
	TitleEn         string                  `json:"title_en"`
	TitleAr         string                  `json:"title_ar"`
	DescriptionEn   string                  `json:"description_en"`
	DescriptionAr   string                  `json:"description_ar"`
	OrderNum        int                     `json:"order_num"`
	Identifications []RiskIdentificationDTO `json:"identifications,omitempty"`
	ActionPlans     []RiskActionPlanDTO     `json:"action_plans,omitempty"`
	KRIs            []RiskReportKriDTO      `json:"kris,omitempty"`
}

// RiskIdentificationDTO represents an identified risk
type RiskIdentificationDTO struct {
	ReferenceNumber        string    `json:"reference_number"`
	RiskCategoryTitleEn    string    `json:"risk_category_title_en"`
	RiskCategoryTitleAr    string    `json:"risk_category_title_ar"`
	Title                  string    `json:"title"`
	Description            *string   `json:"description"`
	RiskSourceTitleEn      *string   `json:"risk_source_title_en"`
	RiskSourceTitleAr      *string   `json:"risk_source_title_en_ar"`
	ReportedAt             time.Time `json:"reported_at"`
	ResponsibleUnitTitleEn string    `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr string    `json:"responsible_unit_title_ar"`
	RiskStatusTitleEn      string    `json:"risk_status_title_en"`
	RiskStatusTitleAr      string    `json:"risk_status_title_ar"`
	RiskStatusColorCode    string    `json:"risk_status_color_code"`
}

// RiskActionPlanDTO represents an action plan under mitigation
type RiskActionPlanDTO struct {
	Title                 string     `json:"title"`
	Description           *string    `json:"description"`
	SeverityLevelTitleEn  *string    `json:"severity_level_title_en"`
	SeverityLevelTitleAr  *string    `json:"severity_level_title_ar"`
	ResponsibleUserNameEn string     `json:"responsible_user_name_en"`
	ResponsibleUserNameAr string     `json:"responsible_user_name_ar"`
	TargetDate            *time.Time `json:"target_date"`
	Status                string     `json:"status"`
}

// RiskReportKriDTO represents a Key Risk Indicator (KRI)
type RiskReportKriDTO struct {
	RiskReferenceNumber string  `json:"risk_reference_number"`
	RiskTitle           string  `json:"risk_title"`
	Title               string  `json:"title"`
	Description         *string `json:"description"`
	LowThreshold        *int    `json:"low_threshold"`
	MediumThreshold     *int    `json:"medium_threshold"`
	UnitTitleEn         *string `json:"unit_title_en"`
	UnitTitleAr         *string `json:"unit_title_ar"`
	CurrentValue        *int    `json:"current_value"`
}

type RiskReportFormDTO struct {
	Title               string    `json:"title" binding:"required"`
	StartDate           time.Time `json:"start_date" binding:"required"`
	EndDate             time.Time `json:"end_date" binding:"required"`
	OrganizationUnitIDs []uint    `json:"organization_unit_ids"`
}

type RiskRegisterDTO struct {
	ID                     uint      `json:"id"`
	ReferenceNumber        string    `json:"reference_number"`
	Title                  string    `json:"title"`
	Description            *string   `json:"description"`
	ReportedAt             time.Time `json:"reported_at"`
	RiskCategoryTitleEn    string    `json:"risk_category_title_en"`
	RiskCategoryTitleAr    string    `json:"risk_category_title_ar"`
	RiskSourceTitleEn      *string   `json:"risk_source_title_en"`
	RiskSourceTitleAr      *string   `json:"risk_source_title_ar"`
	ResponsibleUnitTitleEn string    `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr string    `json:"responsible_unit_title_ar"`

	InherentRiskRatingTitleEn   *string `json:"inherent_risk_rating_title_en"`
	InherentRiskRatingTitleAr   *string `json:"inherent_risk_rating_title_ar"`
	InherentRiskRatingColorCode *string `json:"inherent_risk_rating_color_code"`

	InherentLikelihoodTitleEn *string `json:"inherent_likelihood_title_en"`
	InherentLikelihoodTitleAr *string `json:"inherent_likelihood_title_ar"`

	InherentImpactTitleEn *string `json:"inherent_impact_title_en"`
	InherentImpactTitleAr *string `json:"inherent_impact_title_ar"`

	ResidualRiskRatingTitleEn   *string `json:"residual_risk_rating_title_en"`
	ResidualRiskRatingTitleAr   *string `json:"residual_risk_rating_title_ar"`
	ResidualRiskRatingColorCode *string `json:"residual_risk_rating_color_code"`

	ResidualLikelihoodTitleEn *string `json:"residual_likelihood_title_en"`
	ResidualLikelihoodTitleAr *string `json:"residual_likelihood_title_ar"`

	ResidualImpactTitleEn *string `json:"residual_impact_title_en"`
	ResidualImpactTitleAr *string `json:"residual_impact_title_ar"`

	ExistingControls *string `json:"existing_controls"`
	ProposedControls *string `json:"proposed_controls"`

	RiskStatusTitleEn   string `json:"risk_status_title_en"`
	RiskStatusTitleAr   string `json:"risk_status_title_ar"`
	RiskStatusColorCode string `json:"risk_status_color_code"`
}

type RiskReportActionPlanFormDTO struct {
	ID                         uint      `json:"id"`
	RiskReportID               uint      `json:"risk_report_id"`
	RiskReportIdentificationID uint      `json:"risk_report_identification_id"`
	ActionPlanTitle            string    `json:"action_plan_title"`
	ActionPlanDescription      *string   `json:"action_plan_description"`
	SeverityLevelTitleEn       string    `json:"severity_level_title_en"`
	SeverityLevelTitleAr       string    `json:"severity_level_title_ar"`
	ResponsibleUserNameEn      string    `json:"responsible_user_name_en"`
	ResponsibleUserNameAr      string    `json:"responsible_user_name_ar"`
	TargetDate                 time.Time `json:"target_date"`
	ActionPlanStatusTitleEn    string    `json:"action_plan_status_title_en"`
	ActionPlanStatusTitleAr    string    `json:"action_plan_status_title_ar"`
	ActionPlanStatusColorCode  string    `json:"action_plan_status_color_code"`
}

type KRIDTO struct {
	Title           string  `json:"title"`
	Description     *string `json:"description,omitempty"`
	LowThreshold    int     `json:"low_threshold"`
	MediumThreshold int     `json:"medium_threshold"`
	UnitTitleEn     string  `json:"unit_title_en"`
	UnitTitleAr     string  `json:"unit_title_ar"`
	CurrentValue    int     `json:"current_value"`
}
