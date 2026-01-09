package dtos

// RiskRegisterExportPDFDTO represents the data structure for Risk Register PDF export
type RiskRegisterExportPDFDTO struct {
	Title                    string                         `json:"title"`
	CreatedAt                string                         `json:"created_at"`
	CreatedUserFirstNameEn   string                         `json:"created_user_first_name_en"`
	CreatedUserLastNameEn    string                         `json:"created_user_last_name_en"`
	OrganizationUnitTitlesEn string                         `json:"organization_unit_titles_en"`
	StartDate                string                         `json:"start_date"`
	EndDate                  string                         `json:"end_date"`
	Sections                 []RiskRegisterExportSectionDTO `json:"sections"`
}

// RiskRegisterExportSectionDTO represents a section in the Risk Register PDF export
type RiskRegisterExportSectionDTO struct {
	Type        string                      `json:"type"`
	TitleEn     string                      `json:"title_en"`
	Description string                      `json:"description"`
	Risks       []RiskRegisterExportRiskDTO `json:"risks"`
}

// RiskRegisterExportRiskDTO represents a risk item in the Risk Register PDF export
type RiskRegisterExportRiskDTO struct {
	ReferenceNumber            string  `json:"reference_number"`
	Title                      string  `json:"title"`
	Description                *string `json:"description"`
	ReportedAt                 string  `json:"reported_at"`
	RiskCategoryTitleEn        string  `json:"risk_category_title_en"`
	RiskCategoryTitleAr        string  `json:"risk_category_title_ar"`
	RiskSourceTitleEn          *string `json:"risk_source_title_en"`
	RiskSourceTitleAr          *string `json:"risk_source_title_ar"`
	ResponsibleUnitTitleEn     string  `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr     string  `json:"responsible_unit_title_ar"`
	ResponsibleUserFirstNameEn string  `json:"responsible_user_first_name_en"`
	ResponsibleUserLastNameEn  string  `json:"responsible_user_last_name_en"`
	ResponsibleUserInitial     string  `json:"responsible_user_initial"`

	// Current Risk Rating
	CurrentRiskRatingTitleEn       *string `json:"current_risk_rating_title_en"`
	CurrentRiskRatingTitleAr       *string `json:"current_risk_rating_title_ar"`
	CurrentRiskRatingColorCode     *string `json:"current_risk_rating_color_code"`
	CurrentRiskRatingValue         *uint   `json:"current_risk_rating_value"`
	CurrentLikelihoodTitleEn       *string `json:"current_likelihood_title_en"`
	CurrentLikelihoodTitleAr       *string `json:"current_likelihood_title_ar"`
	CurrentLikelihoodScore         *uint   `json:"current_likelihood_score"`
	CurrentLikelihoodDescriptionEn *string `json:"current_likelihood_description_en"`
	CurrentImpactTitleEn           *string `json:"current_impact_title_en"`
	CurrentImpactTitleAr           *string `json:"current_impact_title_ar"`
	CurrentImpactScore             *uint   `json:"current_impact_score"`
	CurrentImpactDescriptionEn     *string `json:"current_impact_description_en"`

	// Inherent Risk Rating
	InherentRiskRatingTitleEn       *string `json:"inherent_risk_rating_title_en"`
	InherentRiskRatingTitleAr       *string `json:"inherent_risk_rating_title_ar"`
	InherentRiskRatingColorCode     *string `json:"inherent_risk_rating_color_code"`
	InherentRiskRatingValue         *uint   `json:"inherent_risk_rating_value"`
	InherentLikelihoodTitleEn       *string `json:"inherent_likelihood_title_en"`
	InherentLikelihoodTitleAr       *string `json:"inherent_likelihood_title_ar"`
	InherentLikelihoodScore         *uint   `json:"inherent_likelihood_score"`
	InherentLikelihoodDescriptionEn *string `json:"inherent_likelihood_description_en"`
	InherentImpactTitleEn           *string `json:"inherent_impact_title_en"`
	InherentImpactTitleAr           *string `json:"inherent_impact_title_ar"`
	InherentImpactScore             *uint   `json:"inherent_impact_score"`
	InherentImpactDescriptionEn     *string `json:"inherent_impact_description_en"`

	// Residual Risk Rating
	ResidualRiskRatingTitleEn       *string `json:"residual_risk_rating_title_en"`
	ResidualRiskRatingTitleAr       *string `json:"residual_risk_rating_title_ar"`
	ResidualRiskRatingColorCode     *string `json:"residual_risk_rating_color_code"`
	ResidualRiskRatingValue         *uint   `json:"residual_risk_rating_value"`
	ResidualLikelihoodTitleEn       *string `json:"residual_likelihood_title_en"`
	ResidualLikelihoodTitleAr       *string `json:"residual_likelihood_title_ar"`
	ResidualLikelihoodScore         *uint   `json:"residual_likelihood_score"`
	ResidualLikelihoodDescriptionEn *string `json:"residual_likelihood_description_en"`
	ResidualImpactTitleEn           *string `json:"residual_impact_title_en"`
	ResidualImpactTitleAr           *string `json:"residual_impact_title_ar"`
	ResidualImpactScore             *uint   `json:"residual_impact_score"`
	ResidualImpactDescriptionEn     *string `json:"residual_impact_description_en"`

	// Severity Level
	SeverityLevelTitleEn *string `json:"severity_level_title_en"`
	SeverityLevelTitleAr *string `json:"severity_level_title_ar"`

	// Review Information
	ReviewFrequencyTitleEn *string `json:"review_frequency_title_en"`
	ReviewFrequencyTitleAr *string `json:"review_frequency_title_ar"`
	NextReviewDateString   *string `json:"next_review_date"`
	LastReviewNote         *string `json:"last_review_note"`
	NextReviewLevel        *int    `json:"next_review_level"`

	// Submission & Review History
	SubmittedByUserFirstNameEn *string `json:"submitted_by_user_first_name_en"`
	SubmittedByUserLastNameEn  *string `json:"submitted_by_user_last_name_en"`
	SubmittedAtString          *string `json:"submitted_at"`
	UpdatedByUserFirstNameEn   string  `json:"updated_by_user_first_name_en"`
	UpdatedByUserLastNameEn    string  `json:"updated_by_user_last_name_en"`
	UpdatedAtString            *string `json:"updated_at"`

	ExistingControls *string `json:"existing_controls"`
	ProposedControls *string `json:"proposed_controls"`

	RiskStatusTitleEn   string `json:"risk_status_title_en"`
	RiskStatusTitleAr   string `json:"risk_status_title_ar"`
	RiskStatusColorCode string `json:"risk_status_color_code"`

	RiskResponseStrategyTitleEn *string `json:"risk_response_strategy_title_en"`
	RiskResponseStrategyTitleAr *string `json:"risk_response_strategy_title_ar"`

	ActionPlans        []RiskReportActionPlanResponseDTO `json:"action_plans"`
	KeyRiskIdentifiers []RiskReportKriResponseDTO        `json:"key_risk_identifiers,omitempty"`

	// Detailed Control Information
	Controls []RiskRegisterExportControlDTO `json:"controls,omitempty"`

	// Associated Items
	Assets []RiskRegisterExportAssetDTO `json:"assets,omitempty"`

	// Root Cause Analysis
	RootCauseAnalyses []RiskRegisterExportRootCauseAnalysisDTO `json:"root_cause_analyses,omitempty"`
}

// RiskRegisterExportControlDTO represents control details in the Risk Register PDF export
type RiskRegisterExportControlDTO struct {
	ReferenceNumber             string  `json:"reference_number"`
	TitleEn                     string  `json:"title_en"`
	Code                        string  `json:"code"`
	TestScore                   *uint   `json:"test_score"`
	ImplementationStatusTitleEn *string `json:"implementation_status_title_en"`
	EfficiencyStatusTitleEn     *string `json:"efficiency_status_title_en"`
	TestedAtString              *string `json:"tested_at"`
	TestedByUserFirstNameEn     string  `json:"tested_by_user_first_name_en"`
	TestedByUserLastNameEn      string  `json:"tested_by_user_last_name_en"`
	TestedByUserInitial         string  `json:"tested_by_user_initial"`
}

// RiskRegisterExportAssetDTO represents asset details in the Risk Register PDF export
type RiskRegisterExportAssetDTO struct {
	ReferenceNumber  string  `json:"reference_number"`
	TitleEn          string  `json:"title_en"`
	TitleAr          string  `json:"title_ar"`
	AssetTypeTitleEn *string `json:"asset_type_title_en"`
	Description      *string `json:"description"`
}

// RiskRegisterExportRootCauseAnalysisDTO represents root cause analysis in the Risk Register PDF export
type RiskRegisterExportRootCauseAnalysisDTO struct {
	Why      uint   `json:"why"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}
