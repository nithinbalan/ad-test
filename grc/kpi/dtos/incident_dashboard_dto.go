package dtos

type IncidentDashboardQueryParamsDTO struct {
	DateFrom string
	DateTo   string
	Year     int
}

type IncidentStatusCountDTO struct {
	ID        uint   `json:"id"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	ColorCode string `json:"color_code"`
	Type      string `json:"type"`
	Count     int64  `json:"count"`
}

type IncidentStatusSummaryDTO struct {
	TotalIncidents  int64                    `json:"total_incidents"`
	OpenIncidents   int64                    `json:"open_incidents"`
	ClosedIncidents int64                    `json:"closed_incidents"`
	Statuses        []IncidentStatusCountDTO `json:"statuses"`
}

type IncidentSeverityCountDTO struct {
	ID         uint    `json:"id"`
	TitleEn    string  `json:"title_en"`
	TitleAr    string  `json:"title_ar"`
	ColorCode  string  `json:"color_code"`
	Type       string  `json:"type"`
	Count      int64   `json:"count"`
	Percentage float64 `json:"percentage"`
}

type IncidentSeverityBreakdownDTO struct {
	TotalIncidents int64                      `json:"total_incidents"`
	SeverityLevels []IncidentSeverityCountDTO `json:"severity_levels"`
}

type IncidentAcknowledgementSummaryDTO struct {
	Total           int64 `json:"total"`
	Acknowledged    int64 `json:"acknowledged"`
	NotAcknowledged int64 `json:"not_acknowledged"`
}

type IncidentSLASummaryDTO struct {
	SLAAchieved      int64   `json:"sla_achieved"`
	SLAExceeded      int64   `json:"sla_exceeded"`
	SLACompliancePct float64 `json:"sla_compliance_pct"`
}

type IncidentCategoryCountDTO struct {
	CategoryID uint   `json:"category_id"`
	TitleEn    string `json:"title_en"`
	TitleAr    string `json:"title_ar"`
	Count      int64  `json:"count"`
}

type IncidentTypeCountDTO struct {
	TypeID  uint   `json:"type_id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Count   int64  `json:"count"`
}

type IncidentTrendDTO struct {
	Month  string `json:"month"`
	Opened int64  `json:"opened"`
	Closed int64  `json:"closed"`
}

type IncidentActionPlanSummaryDTO struct {
	Total     int64 `json:"total"`
	Completed int64 `json:"completed"`
	Open      int64 `json:"open"`
	Overdue   int64 `json:"overdue"`
}

type OverdueActionPlanDTO struct {
	ActionPlanID                    uint    `json:"action_plan_id"`
	Title                           string  `json:"title"`
	IncidentRegisterID              uint    `json:"incident_register_id"`
	IncidentReferenceNumber         string  `json:"incident_reference_number"`
	IncidentTitle                   string  `json:"incident_title"`
	SeverityLevelTitleEn            string  `json:"severity_level_title_en"`
	SeverityLevelTitleAr            string  `json:"severity_level_title_ar"`
	SeverityLevelColorCode          string  `json:"severity_level_color_code"`
	SeverityLevelType               string  `json:"severity_level_type"`
	ResponsibleUserID               uint    `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn      string  `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr      string  `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn       string  `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr       string  `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          string  `json:"responsible_user_initial"`
	ResponsibleUserColor            string  `json:"responsible_user_color"`
	ResponsibleUserEmail            string  `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath string  `json:"responsible_user_profile_image_path"`
	TargetDate                      *string `json:"target_date"`
	OverdueDays                     int     `json:"overdue_days"`
}
