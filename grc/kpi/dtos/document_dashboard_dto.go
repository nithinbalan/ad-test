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

type DocumentConfidentialityClassificationDTO struct {
	ConfidentialityClassificationID uint    `json:"confidentiality_classification_id"`
	TitleEn                         string  `json:"title_en"`
	TitleAr                         string  `json:"title_ar"`
	Count                           int64   `json:"count"`
	Percentage                      float64 `json:"percentage"`
}

type DocumentResponsibleUnitDTO struct {
	ResponsibleUnitID uint    `json:"responsible_unit_id"`
	TitleEn           string  `json:"title_en"`
	TitleAr           string  `json:"title_ar"`
	Count             int64   `json:"count"`
	Percentage        float64 `json:"percentage"`
}

type DocumentReviewCountTrendDTO struct {
	Month      string `json:"month"`
	MonthLabel string `json:"month_label"`
	Count      int64  `json:"count"`
}
