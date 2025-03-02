package dtos

type IncidentByCategoryDTO struct {
	IncidentCategoryID      uint   `json:"id"`
	IncidentCategoryTitleEn string `json:"title_en"`
	IncidentCategoryTitleAr string `json:"title_ar"`
	Count                   uint   `json:"count"`
}
type IncidentByStatusDTO struct {
	ID        uint   `json:"id"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	ColorCode string `json:"color_code"`
	Count     uint   `json:"count"`
}

type IncidentBySeverityLevelDTO struct {
	ID        uint   `json:"id"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	ColorCode string `json:"color_code"`
	Count     uint   `json:"count"`
}

type IncidentByResponsibleUnitDTO struct {
	ResponsibleUnitID        uint   `json:"id"`
	ResponsibleUnitTitleEn   string `json:"title_en"`
	ResponsibleUnitTitleAr   string `json:"title_ar"`
	ResponsibleUnitColorCode string `json:"color_code"`
	Count                    uint   `json:"count"`
}

type IncidentByReportedUserDTO struct {
	ReportedUserID               uint   `json:"id"`
	ReportedUserFirstNameEn      string `json:"first_name_en"`
	ReportedUserFirstNameAr      string `json:"first_name_ar"`
	ReportedUserLastNameEn       string `json:"last_name_en"`
	ReportedUserLastNameAr       string `json:"last_name_ar"`
	ReportedUserTitleEn          string `json:"title_en"`
	ReportedUserTitleAr          string `json:"title_ar"`
	ReportedUserInitial          string `json:"initial"`
	ReportedUserColor            string `json:"color"`
	ReportedUserEmail            string `json:"email"`
	ReportedUserProfileImagePath string `json:"profile_image"`
	Count                        int64  `json:"count"`
}

type IncidentByReportedYearDTO struct {
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Count   uint   `json:"count"`
}

type IncidentByAssetDTO struct {
	AssetID uint   `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Count   uint   `json:"count"`
}

type IncidentBySLADTO struct {
	SLAID   uint   `json:"sla_id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Count   uint   `json:"count"`
}

type IncidentInvestigationSummaryDTO struct {
	InvestigationID         uint   `json:"investigation_id"`
	IncidentID              uint   `json:"incident_id"`
	IncidentTitle           string `json:"incident_title"`
	InvestigatorFirstNameEn string `json:"investigator_first_name_en"`
	InvestigatorLastNameEn  string `json:"investigator_last_name_en"`
	InvestigatorFirstNameAr string `json:"investigator_first_name_ar"`
	InvestigatorLastNameAr  string `json:"investigator_last_name_ar"`
	InvestigationTeamsEn    string `json:"investigation_teams_en"`
	InvestigationTeamsAr    string `json:"investigation_teams_ar"`
	IncidentStatusTitleEn   string `json:"incident_status_title_en"`
	IncidentStatusTitleAr   string `json:"incident_status_title_ar"`
	FindingsSummary         string `json:"findings_summary"`
	RootCauseAnalysis       string `json:"root_cause_analysis"`
}

type IncidentRCADTO struct {
	IncidentID          uint   `json:"incident_id"`
	IncidentTitle       string `json:"incident_title"`
	RootCauseIdentified string `json:"root_cause_identified"`
}

type IncidentClosureDTO struct {
	IncidentID            uint   `json:"incident_id"`
	IncidentTitle         string `json:"incident_title"`
	ClosureDate           string `json:"closure_date"` // Use `string` to handle formatting (e.g., ISO8601)
	ResolutionSummary     string `json:"resolution_summary"`
	FinalStatus           string `json:"final_status"`
	AssociatedActionPlans string `json:"associated_action_plans"`
}

type OpenIncidentWithRiskDTO struct {
	IncidentID      uint   `json:"incident_id"`
	IncidentTitle   string `json:"incident_title"`
	AssociatedRisks string `json:"associated_risks"`
	Severity        string `json:"severity"`
	CurrentStatus   string `json:"current_status"`
	ResponsibleUser string `json:"responsible_user"`
}

type IncidentPendingActionPlanDTO struct {
	IncidentID      uint   `json:"incident_id"`
	IncidentTitle   string `json:"incident_title"`
	ActionPlanID    uint   `json:"action_plan_id"`
	ActionPlanTitle string `json:"action_plan_title"`
	Description     string `json:"description"`
	DueDate         string `json:"due_date"`
	Status          string `json:"status"`
}

type IncidentActionPlanByStatusDTO struct {
	ID        uint   `json:"id"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	ColorCode string `json:"color_code"`
	Count     uint   `json:"count"`
}

type IncidentActionPlanByResponsibleUnitDTO struct {
	ResponsibleUnitID      uint   `json:"id"`
	ResponsibleUnitTitleEn string `json:"title_en"`
	ResponsibleUnitTitleAr string `json:"title_ar"`
	ColorCode              string `json:"color_code"`
	Count                  uint   `json:"count"`
}

type IncidentActionPlanByResponsibleUserDTO struct {
	ResponsibleUserID               uint   `json:"id"`
	ResponsibleUserFirstNameEn      string `json:"first_name_en"`
	ResponsibleUserFirstNameAr      string `json:"first_name_ar"`
	ResponsibleUserLastNameEn       string `json:"last_name_en"`
	ResponsibleUserLastNameAr       string `json:"last_name_ar"`
	ResponsibleUserTitleEn          string `json:"title_en"`
	ResponsibleUserTitleAr          string `json:"title_ar"`
	ResponsibleUserInitial          string `json:"initial"`
	ResponsibleUserColor            string `json:"color"`
	ResponsibleUserEmail            string `json:"email"`
	ResponsibleUserProfileImagePath string `json:"profile_image"`
	Count                           int64  `json:"count"`
}
