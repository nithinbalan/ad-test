package dtos

import commonDTOs "jethur/grc/common/dtos"

// RiskRegisterExportPDFDTO represents the data structure for Risk Register PDF export
type RiskRegisterExportPDFDTO struct {
	Title                    string  `json:"title"`
	CreatedAt                string  `json:"created_at"`
	GeneratedAt              string  `json:"generated_at"`
	ReferenceNumber          string  `json:"reference_number"`
	Description              *string `json:"description"`
	CreatedUserFirstNameEn   string  `json:"created_user_first_name_en"`
	CreatedUserLastNameEn    string  `json:"created_user_last_name_en"`
	OrganizationUnitTitlesEn string  `json:"organization_unit_titles_en"`
	StartDate                string  `json:"start_date"`
	EndDate                  string  `json:"end_date"`

	// Flattened risk details for template convenience
	RiskStatusTitleEn   string  `json:"risk_status_title_en"`
	RiskStatusTitleAr   string  `json:"risk_status_title_ar"`
	RiskStatusColorCode string  `json:"risk_status_color_code"`
	RiskCategoryTitleEn string  `json:"risk_category_title_en"`
	RiskCategoryTitleAr string  `json:"risk_category_title_ar"`
	RiskSourceTitleEn   *string `json:"risk_source_title_en"`
	RiskSourceTitleAr   *string `json:"risk_source_title_ar"`

	ResponsibleUnitTitleEn     string `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr     string `json:"responsible_unit_title_ar"`
	ResponsibleUserFirstNameEn string `json:"responsible_user_first_name_en"`
	ResponsibleUserLastNameEn  string `json:"responsible_user_last_name_en"`
	ResponsibleUserInitial     string `json:"responsible_user_initial"`

	RiskCode         *string `json:"risk_code"`
	DivisionTitleEn  string  `json:"division_title_en"`
	DivisionTitleAr  string  `json:"division_title_ar"`
	SectionTitleEn   string  `json:"section_title_en"`
	SectionTitleAr   string  `json:"section_title_ar"`
	RegistrationDate *string `json:"registration_date"`
	IsInternal       bool    `json:"is_internal"`
	IsExternal       bool    `json:"is_external"`
	IsStrategic      bool    `json:"is_strategic"`
	IsOperational    bool    `json:"is_operational"`
	ExpectedCause    *string `json:"expected_cause"`
	Consequence      *string `json:"consequence"`

	ReviewFrequencyTitleEn *string `json:"review_frequency_title_en"`
	ReviewFrequencyTitleAr *string `json:"review_frequency_title_ar"`
	NextReviewDateString   *string `json:"next_review_date"`
	LastReviewNote         *string `json:"last_review_note"`
	NextReviewLevel        *int    `json:"next_review_level"`

	SeverityLevelTitleEn *string `json:"severity_level_title_en"`
	SeverityLevelTitleAr *string `json:"severity_level_title_ar"`

	SubmittedByUserFirstNameEn *string `json:"submitted_by_user_first_name_en"`
	SubmittedByUserLastNameEn  *string `json:"submitted_by_user_last_name_en"`
	SubmittedAtString          *string `json:"submitted_at"`
	UpdatedByUserFirstNameEn   *string `json:"updated_by_user_first_name_en"`
	UpdatedByUserLastNameEn    *string `json:"updated_by_user_last_name_en"`
	UpdatedAtString            *string `json:"updated_at"`

	RiskResponseStrategyTitleEn           *string `json:"risk_response_strategy_title_en"`
	RiskResponseStrategyTitleAr           *string `json:"risk_response_strategy_title_ar"`
	RiskResponseStrategyDescription       *string `json:"risk_response_strategy_description"`
	RiskResponseStrategyTreatmentGuidance *string `json:"risk_response_strategy_treatment_guidance"`

	RootCauseAnalyses []RiskRegisterExportRootCauseAnalysisDTO `json:"root_cause_analyses,omitempty"`

	// Current / Inherent / Residual ratings (flattened)
	CurrentRiskRatingTitleEn       *string `json:"current_risk_rating_title_en"`
	CurrentRiskRatingTitleAr       *string `json:"current_risk_rating_title_ar"`
	CurrentRiskRatingColorCode     *string `json:"current_risk_rating_color_code"`
	CurrentRiskRatingValue         *uint   `json:"current_risk_rating_value"`
	CurrentRiskPointerPosition     *string `json:"current_risk_pointer_position"`
	CurrentLikelihoodTitleEn       *string `json:"current_likelihood_title_en"`
	CurrentLikelihoodTitleAr       *string `json:"current_likelihood_title_ar"`
	CurrentLikelihoodScore         *uint   `json:"current_likelihood_score"`
	CurrentLikelihoodDescriptionEn *string `json:"current_likelihood_description_en"`
	CurrentImpactTitleEn           *string `json:"current_impact_title_en"`
	CurrentImpactTitleAr           *string `json:"current_impact_title_ar"`
	CurrentImpactScore             *uint   `json:"current_impact_score"`
	CurrentImpactDescriptionEn     *string `json:"current_impact_description_en"`

	InherentRiskRatingTitleEn       *string `json:"inherent_risk_rating_title_en"`
	InherentRiskRatingTitleAr       *string `json:"inherent_risk_rating_title_ar"`
	InherentRiskRatingColorCode     *string `json:"inherent_risk_rating_color_code"`
	InherentRiskRatingValue         *uint   `json:"inherent_risk_rating_value"`
	InherentRiskPointerPosition     *string `json:"inherent_risk_pointer_position"`
	InherentLikelihoodTitleEn       *string `json:"inherent_likelihood_title_en"`
	InherentLikelihoodTitleAr       *string `json:"inherent_likelihood_title_ar"`
	InherentLikelihoodScore         *uint   `json:"inherent_likelihood_score"`
	InherentLikelihoodDescriptionEn *string `json:"inherent_likelihood_description_en"`
	InherentImpactTitleEn           *string `json:"inherent_impact_title_en"`
	InherentImpactTitleAr           *string `json:"inherent_impact_title_ar"`
	InherentImpactScore             *uint   `json:"inherent_impact_score"`
	InherentImpactDescriptionEn     *string `json:"inherent_impact_description_en"`
	InherentImpactCategoryTitleEn   *string `json:"inherent_impact_category_title_en"`
	InherentImpactCategoryTitleAr   *string `json:"inherent_impact_category_title_ar"`
	InherentImpactGuideline         *string `json:"inherent_impact_guideline"`

	ResidualRiskRatingTitleEn       *string `json:"residual_risk_rating_title_en"`
	ResidualRiskRatingTitleAr       *string `json:"residual_risk_rating_title_ar"`
	ResidualRiskRatingColorCode     *string `json:"residual_risk_rating_color_code"`
	ResidualRiskRatingValue         *uint   `json:"residual_risk_rating_value"`
	ResidualRiskPointerPosition     *string `json:"residual_risk_pointer_position"`
	ResidualLikelihoodTitleEn       *string `json:"residual_likelihood_title_en"`
	ResidualLikelihoodTitleAr       *string `json:"residual_likelihood_title_ar"`
	ResidualLikelihoodScore         *uint   `json:"residual_likelihood_score"`
	ResidualLikelihoodDescriptionEn *string `json:"residual_likelihood_description_en"`
	ResidualImpactTitleEn           *string `json:"residual_impact_title_en"`
	ResidualImpactTitleAr           *string `json:"residual_impact_title_ar"`
	ResidualImpactScore             *uint   `json:"residual_impact_score"`
	ResidualImpactDescriptionEn     *string `json:"residual_impact_description_en"`

	// Control Effectiveness
	AverageControlEffectivenessScore *float64                              `json:"average_control_effectiveness_score"`
	ControlEffectivenessRatings      []ControlEffectivenessRatingExportDTO `json:"control_effectiveness_ratings"`

	// Flattened collections
	ActionPlans            []RiskReportActionPlanResponseDTO                 `json:"action_plans"`
	Controls               []RiskRegisterExportControlDTO                    `json:"controls"`
	KeyRiskIdentifiers     []RiskReportKriResponseDTO                        `json:"key_risk_identifiers,omitempty"`
	Assets                 []RiskRegisterExportAssetDTO                      `json:"assets,omitempty"`
	Documents              []*commonDTOs.AssociatedDocumentListDTO           `json:"documents,omitempty"`
	Processes              []*commonDTOs.AssociatedProcessListDTO            `json:"processes,omitempty"`
	KPIs                   []*commonDTOs.AssociatedKPIRegisterListDTO        `json:"kpis,omitempty"`
	StrategicObjectives    []*commonDTOs.AssociatedStrategicObjectiveListDTO `json:"strategic_objectives,omitempty"`
	ComplianceRequirements []*commonDTOs.AssociatedClauseListDTO             `json:"compliance_requirements,omitempty"`
	Incidents              []*commonDTOs.AssociatedIncidentListDTO           `json:"incidents,omitempty"`

	// Risk Appetite Statements
	RiskAppetiteStatements         []RiskRegisterExportRiskAppetiteStatementDTO `json:"risk_appetite_statements,omitempty"`
	ResidualRiskAppetiteStatements []RiskRegisterExportRiskAppetiteStatementDTO `json:"residual_risk_appetite_statements,omitempty"`

	// Risk Closure
	ClosureJustification *string `json:"closure_justification"`
	LessonsLearned       *string `json:"lessons_learned"`
	WhatWorkedWell       *string `json:"what_worked_well"`
	WhatNeedsImprovement *string `json:"what_needs_improvement"`
	KeyTakeaways         *string `json:"key_takeaways"`

	// Control Statistics
	ControlCoveragePercent   *float64 `json:"control_coverage_percent"`
	IneffectiveControlsCount *uint    `json:"ineffective_controls_count"`
	ImplementedPercent       *float64 `json:"implemented_percent"`
	EffectiveControlsPercent *float64 `json:"effective_controls_percent"`
	EffectiveControlsCount   *uint    `json:"effective_controls_count"`
	ControlScoreAchieved     *uint    `json:"control_score_achieved"`
	ControlScoreRemaining    *uint    `json:"control_score_remaining"`

	// Reviews
	Reviews []RiskRegisterExportReviewDTO `json:"reviews,omitempty"`

	Sections []RiskRegisterExportSectionDTO `json:"sections"`
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
	InherentRiskPointerPosition     *string `json:"inherent_risk_pointer_position"`
	InherentLikelihoodTitleEn       *string `json:"inherent_likelihood_title_en"`
	InherentLikelihoodTitleAr       *string `json:"inherent_likelihood_title_ar"`
	InherentLikelihoodScore         *uint   `json:"inherent_likelihood_score"`
	InherentLikelihoodDescriptionEn *string `json:"inherent_likelihood_description_en"`
	InherentImpactTitleEn           *string `json:"inherent_impact_title_en"`
	InherentImpactTitleAr           *string `json:"inherent_impact_title_ar"`
	InherentImpactScore             *uint   `json:"inherent_impact_score"`
	InherentImpactDescriptionEn     *string `json:"inherent_impact_description_en"`
	InherentImpactCategoryTitleEn   *string `json:"inherent_impact_category_title_en"`
	InherentImpactCategoryTitleAr   *string `json:"inherent_impact_category_title_ar"`
	InherentImpactGuideline         *string `json:"inherent_impact_guideline"`

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
	ReferenceNumber                        string   `json:"reference_number"`
	TitleEn                                string   `json:"title_en"`
	Code                                   string   `json:"code"`
	ControlTypeTitleEn                     *string  `json:"control_type_title_en"`
	TestScore                              *uint    `json:"test_score"`
	ImplementationStatusTitleEn            *string  `json:"implementation_status_title_en"`
	ImplementationStatusTitleEnValue       string   `json:"implementation_status_title_en_value"`
	ControlEffectivenessRatingTitleEn      *string  `json:"control_effectiveness_rating_title_en"`
	ControlEffectivenessRatingTitleEnValue string   `json:"control_effectiveness_rating_title_en_value"`
	ControlEffectivenessRatingColorCode    string   `json:"control_effectiveness_rating_color_code"`
	ControlEffectivenessRatingType         string   `json:"control_effectiveness_rating_type"`
	ControlEffectivenessScore              *float64 `json:"control_effectiveness_score"`
	RiskReduction                          *float64 `json:"risk_reduction"`
	TestedAtString                         *string  `json:"tested_at"`
	TestedByUserFirstNameEn                string   `json:"tested_by_user_first_name_en"`
	TestedByUserLastNameEn                 string   `json:"tested_by_user_last_name_en"`
	TestedByUserInitial                    string   `json:"tested_by_user_initial"`
	ControlStatusTitleEn                   *string  `json:"control_status_title_en"`
	ControlStatusTitleEnValue              string   `json:"control_status_title_en_value"`
	ControlStatusTitleAr                   *string  `json:"control_status_title_ar"`
	ControlStatusColorCode                 *string  `json:"control_status_color_code"`
	ControlStatusType                      *string  `json:"control_status_type"`
	ControlStatusTypeValue                 string   `json:"control_status_type_value"`
	ControlTestDescription                 *string  `json:"control_test_description"`
}

// RiskRegisterExportAssetDTO represents asset details in the Risk Register PDF export
type RiskRegisterExportAssetDTO struct {
	ReferenceNumber            string  `json:"reference_number"`
	TitleEn                    string  `json:"title_en"`
	TitleAr                    string  `json:"title_ar"`
	AssetTypeTitleEn           *string `json:"asset_type_title_en"`
	CategoryTitleEn            *string `json:"category_title_en"`
	Description                *string `json:"description"`
	ResponsibleUserInitial     *string `json:"responsible_user_initial"`
	ResponsibleUserNameEn      *string `json:"responsible_user_name_en"`
	ResponsibleUserFirstNameEn string  `json:"responsible_user_first_name_en"`
	ResponsibleUserLastNameEn  string  `json:"responsible_user_last_name_en"`
}

// RiskRegisterExportRootCauseAnalysisDTO represents root cause analysis in the Risk Register PDF export
type RiskRegisterExportRootCauseAnalysisDTO struct {
	Why      uint   `json:"why"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type RiskRegisterExportRiskAppetiteStatementDTO struct {
	StatusTitleEn   string `json:"status_title_en"`
	StatusClass     string `json:"status_class"`
	StatusColorCode string `json:"status_color_code"`
	Description     string `json:"description"`
	ToleranceType   string `json:"tolerance_type"`
}

// ControlEffectivenessRatingExportDTO represents control effectiveness rating in the Risk Register PDF export
type ControlEffectivenessRatingExportDTO struct {
	TitleEn          string  `json:"title_en"`
	TitleAr          string  `json:"title_ar"`
	MinScore         float64 `json:"min_score"`
	MaxScore         float64 `json:"max_score"`
	ReductionPercent float64 `json:"reduction_percent"`
	ColorCode        string  `json:"color_code"`
}

// RiskRegisterExportReviewDTO represents review information in the Risk Register PDF export
type RiskRegisterExportReviewDTO struct {
	ReviewerInitial  *string `json:"reviewer_initial"`
	ReviewerNameEn   *string `json:"reviewer_name_en"`
	ReviewDateString *string `json:"review_date_string"`
	StatusTitleEn    *string `json:"status_title_en"`
}

// RiskRegisterExportClosureDTO represents risk closure information in the Risk Register PDF export
type RiskRegisterExportClosureDTO struct {
	Justification        *string `json:"justification"`
	LessonsLearned       *string `json:"lessons_learned"`
	WhatWorkedWell       *string `json:"what_worked_well"`
	WhatNeedsImprovement *string `json:"what_needs_improvement"`
	KeyTakeaways         *string `json:"key_takeaways"`
}
