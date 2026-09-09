package dtos

type ProjectDashboardQueryParamsDTO struct {
	DateFrom string
	DateTo   string
	Year     int
	Limit    int
}

type ProjectDashboardSummaryDTO struct {
	TotalProjects      int64   `json:"total_projects"`
	ActiveProjects     int64   `json:"active_projects"`
	CompletedProjects  int64   `json:"completed_projects"`
	AvgCompletion      float64 `json:"avg_completion"`
	OnTrack            int64   `json:"on_track"`
	AtRisk             int64   `json:"at_risk"`
	Delayed            int64   `json:"delayed"`
	OpenIssues         int64   `json:"open_issues"`
	OpenChangeRequests int64   `json:"open_change_requests"`
	ProjectsDueSoon    int64   `json:"projects_due_soon"`
	OverdueProjects    int64   `json:"overdue_projects"`
}

type ProjectsByPhaseItemDTO struct {
	ID         uint    `json:"id"`
	TitleEn    string  `json:"title_en"`
	TitleAr    string  `json:"title_ar"`
	Type       string  `json:"type"`
	Count      int64   `json:"count"`
	Percentage float64 `json:"percentage"`
}

type ProjectsByPhaseDTO struct {
	Total  int64                    `json:"total"`
	Phases []ProjectsByPhaseItemDTO `json:"phases"`
}

type ProjectsByClassificationItemDTO struct {
	ID         uint    `json:"id"`
	TitleEn    string  `json:"title_en"`
	TitleAr    string  `json:"title_ar"`
	Count      int64   `json:"count"`
	Percentage float64 `json:"percentage"`
}

type ProjectsByClassificationDTO struct {
	Total           int64                             `json:"total"`
	Classifications []ProjectsByClassificationItemDTO `json:"classifications"`
}

type ProjectsByStatusItemDTO struct {
	ID         uint    `json:"id"`
	TitleEn    string  `json:"title_en"`
	TitleAr    string  `json:"title_ar"`
	Type       string  `json:"type"`
	ColorCode  string  `json:"color_code"`
	Count      int64   `json:"count"`
	Percentage float64 `json:"percentage"`
}

type ProjectsByStatusDTO struct {
	Total    int64                     `json:"total"`
	Statuses []ProjectsByStatusItemDTO `json:"statuses"`
}

type ProjectCompletionListDTO struct {
	ID                           uint    `json:"id"`
	ReferenceNumber              string  `json:"reference_number"`
	TitleEn                      string  `json:"title_en"`
	TitleAr                      string  `json:"title_ar"`
	ProjectPhaseID               *uint   `json:"project_phase_id"`
	ProjectPhaseTitleEn          *string `json:"project_phase_title_en"`
	ProjectPhaseTitleAr          *string `json:"project_phase_title_ar"`
	ProjectClassificationID      *uint   `json:"project_classification_id"`
	ProjectClassificationTitleEn *string `json:"project_classification_title_en"`
	ProjectClassificationTitleAr *string `json:"project_classification_title_ar"`
	StartDate                    string  `json:"start_date"`
	StartDateFormatted           string  `json:"start_date_formatted"`
	EndDate                      string  `json:"end_date"`
	EndDateFormatted             string  `json:"end_date_formatted"`
	CompletionPercentage         float64 `json:"completion_percentage"`
	ProjectStatusID              uint    `json:"project_status_id"`
	ProjectStatusTitleEn         string  `json:"project_status_title_en"`
	ProjectStatusTitleAr         string  `json:"project_status_title_ar"`
	ProjectStatusType            string  `json:"project_status_type"`
	ProjectStatusColorCode       string  `json:"project_status_color_code"`
}

type ProjectDueSoonOverdueListDTO struct {
	ID                           uint    `json:"id"`
	ReferenceNumber              string  `json:"reference_number"`
	TitleEn                      string  `json:"title_en"`
	TitleAr                      string  `json:"title_ar"`
	ProjectPhaseID               *uint   `json:"project_phase_id"`
	ProjectPhaseTitleEn          *string `json:"project_phase_title_en"`
	ProjectPhaseTitleAr          *string `json:"project_phase_title_ar"`
	ProjectPhaseType             *string `json:"project_phase_type"`
	ProjectClassificationID      *uint   `json:"project_classification_id"`
	ProjectClassificationTitleEn *string `json:"project_classification_title_en"`
	ProjectClassificationTitleAr *string `json:"project_classification_title_ar"`
	StartDate                    string  `json:"start_date"`
	StartDateFormatted           string  `json:"start_date_formatted"`
	EndDate                      string  `json:"end_date"`
	EndDateFormatted             string  `json:"end_date_formatted"`
	PlannedEndDate               string  `json:"planned_end_date"`
	PlannedEndDateFormatted      string  `json:"planned_end_date_formatted"`
	IsOverdue                    bool    `json:"is_overdue"`
	OverdueDays                  int     `json:"overdue_days"`
	RemainingDays                int     `json:"remaining_days"`
	DaysRemaining                int     `json:"days_remaining"`
	ProjectStatusID              uint    `json:"project_status_id"`
	ProjectStatusTitleEn         string  `json:"project_status_title_en"`
	ProjectStatusTitleAr         string  `json:"project_status_title_ar"`
	ProjectStatusType            string  `json:"project_status_type"`
	ProjectStatusColorCode       string  `json:"project_status_color_code"`
	StatusTitleEn                string  `json:"status_title_en"`
	StatusTitleAr                string  `json:"status_title_ar"`
	StatusType                   string  `json:"status_type"`
	StatusColorCode              string  `json:"status_color_code"`
}
