package dtos

type ProjectSecurityGateDashboardQueryParamsDTO struct {
	DateFrom string
	DateTo   string
	Year     int
}

type ProjectSecurityGateDashboardSummaryDTO struct {
	TotalSecurityGates               int64   `json:"total_security_gates"`
	TotalChecklists                  int64   `json:"total_checklists"`
	CompletedSecurityGates           int64   `json:"completed_security_gates"`
	CompletedSecurityGatesPercentage float64 `json:"completed_security_gates_percentage"`
	CompletedChecklists              int64   `json:"completed_checklists"`
	CompletedChecklistsPercentage    float64 `json:"completed_checklists_percentage"`
	RisksRaised                      int64   `json:"risks_raised"`
	RisksRaisedPercentage            float64 `json:"risks_raised_percentage"`
	CriticalChecklists               int64   `json:"critical_checklists"`
	CriticalChecklistsPercentage     float64 `json:"critical_checklists_percentage"`
}

type ProjectSecurityGateExceptionSummaryDTO struct {
	TotalChecklists              int64   `json:"total_checklists"`
	TotalExceptions              int64   `json:"total_exceptions"`
	TotalExceptionsPercentage    float64 `json:"total_exceptions_percentage"`
	ExceptionsApproved           int64   `json:"exceptions_approved"`
	ExceptionsApprovedPercentage float64 `json:"exceptions_approved_percentage"`
	ExceptionsPending            int64   `json:"exceptions_pending"`
	ExceptionsPendingPercentage  float64 `json:"exceptions_pending_percentage"`
}

type ProjectSecurityGateBreakdownItemDTO struct {
	ID         uint    `json:"id"`
	TitleEn    string  `json:"title_en"`
	TitleAr    string  `json:"title_ar"`
	Type       string  `json:"type"`
	ColorCode  string  `json:"color_code"`
	Count      int64   `json:"count"`
	Percentage float64 `json:"percentage"`
}

type ProjectSecurityGateChecklistsByResponseStatusDTO struct {
	TotalChecklists int64                                 `json:"total_checklists"`
	Statuses        []ProjectSecurityGateBreakdownItemDTO `json:"statuses"`
}

type ProjectSecurityGatesByReadinessStatusDTO struct {
	TotalSecurityGates int64                                 `json:"total_security_gates"`
	Statuses           []ProjectSecurityGateBreakdownItemDTO `json:"statuses"`
}

type ProjectSecurityGatesByClassificationDTO struct {
	TotalSecurityGates   int64                                 `json:"total_security_gates"`
	TotalClassifications int64                                 `json:"total_classifications"`
	Classifications      []ProjectSecurityGateBreakdownItemDTO `json:"classifications"`
}
