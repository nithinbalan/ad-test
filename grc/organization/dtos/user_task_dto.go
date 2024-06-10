package dtos

type UserTaskStatusCountResponseDTO struct {
	ID        int    `json:"id"`
	Type      string `json:"type"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	ColorCode string `json:"color_code"`
	Count     int    `json:"count"`
}

type UserTaskSummaryResponseDTO struct {
	StatusCounts []UserTaskStatusCountResponseDTO `json:"status_counts"`
	TotalCount   int                              `json:"total_count"`
}