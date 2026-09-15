package dtos

type DocumentDashboardQueryParamsDTO struct {
	Filters map[string][]string
}

type DocumentDashboardSummaryDTO struct {
	TotalDocuments                      int64    `json:"total_documents"`
	TotalDocumentsChangePercentage      *float64 `json:"total_documents_change_percentage"`
	ApprovedDocuments                   int64    `json:"approved_documents"`
	ApprovedDocumentsChangePercentage   *float64 `json:"approved_documents_change_percentage"`
	PublishedDocuments                  int64    `json:"published_documents"`
	PublishedDocumentsChangePercentage  *float64 `json:"published_documents_change_percentage"`
	ReviewedOnTimePercentage            float64  `json:"reviewed_on_time_percentage"`
	ReviewCount                         int64    `json:"review_count"`
	AverageApprovalTimeDays             float64  `json:"average_approval_time_days"`
	AverageApprovalTimeChangePercentage *float64 `json:"average_approval_time_change_percentage"`
	TotalECCControls                    int64    `json:"total_ecc_controls"`
	ECCControlsCount                    int64    `json:"ecc_controls_count"`
}

type DocumentStatusByTypeDTO struct {
	DocumentTypeID uint   `json:"document_type_id"`
	Type           string `json:"type"`
	TitleEn        string `json:"title_en"`
	TitleAr        string `json:"title_ar"`
	Draft          int64  `json:"draft"`
	InReview       int64  `json:"in_review"`
	Rejected       int64  `json:"rejected"`
	Approved       int64  `json:"approved"`
	Published      int64  `json:"published"`
	Archived       int64  `json:"archived"`
	Reverted       int64  `json:"reverted"`
	Total          int64  `json:"total"`
}

type DocumentConfidentialityClassificationDTO struct {
	ConfidentialityClassificationID uint    `json:"confidentiality_classification_id"`
	TitleEn                         string  `json:"title_en"`
	TitleAr                         string  `json:"title_ar"`
	Count                           int64   `json:"count"`
	Percentage                      float64 `json:"percentage"`
}

type DocumentConfidentialityClassificationResponseDTO struct {
	Total                          int64                                      `json:"total"`
	ConfidentialityClassifications []DocumentConfidentialityClassificationDTO `json:"confidentiality_classifications"`
}

type DocumentResponsibleUnitDTO struct {
	ResponsibleUnitID uint    `json:"responsible_unit_id"`
	TitleEn           string  `json:"title_en"`
	TitleAr           string  `json:"title_ar"`
	Count             int64   `json:"count"`
	Percentage        float64 `json:"percentage"`
}

type DocumentResponsibleUnitResponseDTO struct {
	Total            int64                        `json:"total"`
	ResponsibleUnits []DocumentResponsibleUnitDTO `json:"responsible_units"`
}

type DocumentReviewCountTrendDTO struct {
	Month      string `json:"month"`
	MonthLabel string `json:"month_label"`
	Count      int64  `json:"count"`
}
