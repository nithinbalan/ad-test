package dtos

import "time"

type RiskReportFormDTO struct {
	Title              string    `json:"title" binding:"required"`
	StartDate          time.Time `json:"start_date" binding:"required"`
	EndDate            time.Time `json:"end_date" binding:"required"`
	ResponsibleUnitIDs []*uint   `json:"responsible_unit_ids"`
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
	ResponsibleUnitTitleEn *string   `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr *string   `json:"responsible_unit_title_ar"`

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

	RiskStatusTitleEn   *string `json:"risk_status_title_en"`
	RiskStatusTitleAr   *string `json:"risk_status_title_ar"`
	RiskStatusColorCode *string `json:"risk_status_color_code"`
}

type RiskReportActionPlanDTO struct {
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
