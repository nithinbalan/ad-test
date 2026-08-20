package dtos

import "time"

type IdentityAccessReviewSummaryDTO struct {
	OverallCompletionPercentage float64 `json:"overall_completion_percentage"`
	OverDueCampaignsCount       int64   `json:"overdue_campaigns_count"`
	ApprovedCampaignsCount      int64   `json:"approved_campaigns_count"`
	TotalCampaignsCount         int64   `json:"total_campaigns_count"`
	OpenCorrectiveActionsCount  int64   `json:"open_corrective_actions_count"`
}

type IdentityAccessReviewCampaignByLifecycleDTO struct {
	CampaignStatusID uint   `json:"campaign_status_id"`
	TitleEn          string `json:"title_en"`
	TitleAr          string `json:"title_ar"`
	ColorCode        string `json:"color_code"`
	Type             string `json:"type"`
	Count            int64  `json:"count"`
}

type IdentityAccessReviewCompletionBySystemDTO struct {
	AssetRegisterID      uint    `json:"asset_register_id"`
	AssetTitleEn         string  `json:"asset_title_en"`
	AssetTitleAr         string  `json:"asset_title_ar"`
	CompletionPercentage float64 `json:"completion_percentage"`
}

type IdentityAccessReviewRecentCampaignDTO struct {
	ID                      uint      `json:"id"`
	ReferenceNumber         string    `json:"reference_number"`
	CampaignName            string    `json:"campaign_name"`
	AssetRegisterID         uint      `json:"asset_register_id"`
	AssetTitleEn            string    `json:"asset_title_en"`
	AssetTitleAr            string    `json:"asset_title_ar"`
	ReviewTypeID            uint      `json:"review_type_id"`
	ReviewTypeTitleEn       string    `json:"review_type_title_en"`
	ReviewTypeTitleAr       string    `json:"review_type_title_ar"`
	ReviewType              string    `json:"review_type"`
	StartDate               string    `json:"start_date"`
	DueDate                 string    `json:"due_date"`
	ReviewerID              uint      `json:"reviewer_id"`
	ReviewerFirstNameEn     string    `json:"reviewer_first_name_en"`
	ReviewerFirstNameAr     string    `json:"reviewer_first_name_ar"`
	ReviewerLastNameEn      string    `json:"reviewer_last_name_en"`
	ReviewerLastNameAr      string    `json:"reviewer_last_name_ar"`
	ReviewerInitial         string    `json:"reviewer_initial"`
	CampaignStatusID        uint      `json:"campaign_status_id"`
	CampaignStatusTitleEn   string    `json:"campaign_status_title_en"`
	CampaignStatusTitleAr   string    `json:"campaign_status_title_ar"`
	CampaignStatusColorCode string    `json:"campaign_status_color_code"`
	CampaignStatusType      string    `json:"campaign_status_type"`
	CompletionPercentage    int       `json:"completion_percentage"`
	CreatedAt               time.Time `json:"created_at"`
}
