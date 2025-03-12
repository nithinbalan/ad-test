package dtos

import "time"

type RiskReportDetailDTO struct {
	ID              uint                    `json:"id"`
	Title           string                  `json:"title"`
	StartDate       time.Time               `json:"start_date"`
	EndDate         time.Time               `json:"end_date"`
	CreatedBy       string                  `json:"created_by"`
	CreatedAt       time.Time               `json:"created_at"`
	Sections        []RiskReportSectionDTO  `json:"sections"`
	Identifications []RiskIdentificationDTO `json:"identifications"`
	ActionPlans     []RiskActionPlanDTO     `json:"action_plans"`
	KRIs            []RiskReportKriDTO      `json:"kris"`
}

type RiskReportSectionDTO struct {
	Type          string `json:"type"`
	TitleEn       string `json:"title_en"`
	TitleAr       string `json:"title_ar"`
	DescriptionEn string `json:"description_en"`
	DescriptionAr string `json:"description_ar"`
	OrderNum      int    `json:"order_num"`
}

type RiskIdentificationDTO struct {
	ReferenceNumber string  `json:"reference_number"`
	Title           string  `json:"title"`
	Description     *string `json:"description"`
	RiskCategory    string  `json:"risk_category"`
	RiskStatus      string  `json:"risk_status"`
}

type RiskActionPlanDTO struct {
	ReferenceNumber  string  `json:"reference_number"`
	Title            string  `json:"title"`
	Description      *string `json:"description"`
	SeverityLevel    string  `json:"severity_level"`
	ResponsibleUser  string  `json:"responsible_user"`
	OrganizationUnit string  `json:"organization_unit"`
	TargetDate       string  `json:"target_date"`
	Status           string  `json:"status"`
}

type RiskReportKriDTO struct {
	ReferenceNumber string  `json:"reference_number"`
	Title           string  `json:"title"`
	Description     *string `json:"description"`
	Threshold       string  `json:"threshold"`
	CurrentValue    string  `json:"current_value"`
	Status          string  `json:"status"`
}

// type RiskReportDetailsDTO struct {
// 	ID                uint                      `json:"id"`
// 	Title             string                    `json:"title"`
// 	StartDateString   *string                   `json:"start_date"`
// 	EndDateString     *string                   `json:"end_date"`
// 	CreatedAt         time.Time                 `json:"created_at"`
// 	UpdatedAt         *time.Time                `json:"updated_at"`
// 	Sections          []RiskReportSectionDTO    `json:"sections"`
// 	Risks             []RiskIdentificationDTO   `json:"risks"`
// 	Assessments       []RiskAssessmentDTO       `json:"assessments"`
// 	ActionPlans       []RiskReportActionPlanDTO `json:"action_plans"`
// 	KRIs              []RiskReportKriDTO        `json:"kris"`
// 	OrganizationUnits []RiskReportOrgUnitDTO    `json:"organization_units"`
// 	CreatedByUser     UserDTO                   `json:"created_by_user"`
// 	UpdatedByUser     *UserDTO                  `json:"updated_by_user,omitempty"`
// }

// type RiskReportSectionDTO struct {
// 	ID            uint   `json:"id"`
// 	Type          string `json:"type"`
// 	TitleEn       string `json:"title_en"`
// 	TitleAr       string `json:"title_ar"`
// 	DescriptionEn string `json:"description_en"`
// 	DescriptionAr string `json:"description_ar"`
// }

// type RiskIdentificationDTO struct {
// 	ID                  uint   `json:"id"`
// 	ReferenceNumber     string `json:"reference_number"`
// 	Title               string `json:"title"`
// 	Description         string `json:"description"`
// 	RiskCategoryTitleEn string `json:"risk_category_title_en"`
// 	RiskCategoryTitleAr string `json:"risk_category_title_ar"`
// 	RiskStatusTitleEn   string `json:"risk_status_title_en"`
// 	RiskStatusTitleAr   string `json:"risk_status_title_ar"`
// }

// type RiskAssessmentDTO struct {
// 	ID              uint   `json:"id"`
// 	ReferenceNumber string `json:"reference_number"`
// 	Title           string `json:"title"`
// 	Impact          string `json:"impact"`
// 	Likelihood      string `json:"likelihood"`
// 	RiskRating      string `json:"risk_rating"`
// }

// type RiskReportActionPlanDTO struct {
// 	ActionPlanTitle       string `json:"action_plan_title"`
// 	ResponsibleUserNameEn string `json:"responsible_user_name_en"`
// }

// type RiskReportKriDTO struct {
// 	Title       string `json:"title"`
// 	UnitTitleEn string `json:"unit_title_en"`
// }

// type RiskReportOrgUnitDTO struct {
// 	OrganizationUnitID uint   `json:"organization_unit_id"`
// 	TitleEn            string `json:"title_en"`
// 	TitleAr            string `json:"title_ar"`
// }

// type UserDTO struct {
// 	ID        uint   `json:"id"`
// 	FirstName string `json:"first_name"`
// 	LastName  string `json:"last_name"`
// 	Email     string `json:"email"`
// }

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
