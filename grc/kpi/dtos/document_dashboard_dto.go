package dtos

type DocumentDashboardQueryParamsDTO struct {
	Filters map[string][]string
}

type DocumentDashboardSummaryDTO struct {
	TotalDocuments                  int64   `json:"total_documents"`
	ApprovedDocuments               int64   `json:"approved_documents"`
	PublishedDocuments              int64   `json:"published_documents"`
	ReviewedOnTimePercentage        float64 `json:"reviewed_on_time_percentage"`
	AverageApprovalTimeDays         float64 `json:"average_approval_time_days"`
	ECCControlsDocumentedPercentage float64 `json:"ecc_controls_documented_percentage"`
}

type DocumentStatusByTypeDTO struct {
	DocumentTypeID uint   `json:"document_type_id"`
	Type           string `json:"type"`
	TitleEn        string `json:"title_en"`
	TitleAr        string `json:"title_ar"`
	Approved       int64  `json:"approved"`
	InReview       int64  `json:"in_review"`
	Draft          int64  `json:"draft"`
	Overdue        int64  `json:"overdue"`
	Total          int64  `json:"total"`
}

type DocumentOnTimeReviewTrendDTO struct {
	Month              string   `json:"month"`
	MonthLabel         string   `json:"month_label"`
	OnTimePercentage   *float64 `json:"on_time_percentage"`
	TargetPercentage   float64  `json:"target_percentage"`
	VarianceFromTarget *float64 `json:"variance_from_target"`
}
