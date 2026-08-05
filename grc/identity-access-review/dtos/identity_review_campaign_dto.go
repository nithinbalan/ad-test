package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	"time"
)

type IdentityReviewCampaignListDTO struct {
	ID                      uint       `json:"id"`
	ReferenceNumber         string     `json:"reference_number"`
	CampaignName            string     `json:"campaign_name"`
	Description             *string    `json:"description"`
	StartDate               string     `json:"start_date"`
	DueDate                 string     `json:"due_date"`
	ReminderDaysOffset      *int       `json:"reminder_days_offset"`
	CompletionPercentage    int        `json:"completion_percentage"`
	ReviewTypeID            uint       `json:"review_type_id"`
	ReviewTypeTitleEn       string     `json:"review_type_title_en"`
	ReviewTypeTitleAr       string     `json:"review_type_title_ar"`
	CampaignStatusID        uint       `json:"campaign_status_id"`
	CampaignStatusTitleEn   string     `json:"campaign_status_title_en"`
	CampaignStatusTitleAr   string     `json:"campaign_status_title_ar"`
	CampaignStatusColorCode string     `json:"campaign_status_color_code"`
	ReviewerID              uint       `json:"reviewer_id"`
	ReviewerFirstNameEn     string     `json:"reviewer_first_name_en"`
	ReviewerFirstNameAr     string     `json:"reviewer_first_name_ar"`
	ReviewerLastNameEn      string     `json:"reviewer_last_name_en"`
	ReviewerLastNameAr      string     `json:"reviewer_last_name_ar"`
	ReviewerInitial         string     `json:"reviewer_initial"`
	ApproverID              uint       `json:"approver_id"`
	ApproverFirstNameEn     string     `json:"approver_first_name_en"`
	ApproverFirstNameAr     string     `json:"approver_first_name_ar"`
	ApproverLastNameEn      string     `json:"approver_last_name_en"`
	ApproverLastNameAr      string     `json:"approver_last_name_ar"`
	ApproverInitial         string     `json:"approver_initial"`
	CreatedAt               time.Time  `json:"created_at"`
	UpdatedAt               *time.Time `json:"updated_at"`
}

type IdentityReviewCampaignListQueryParamsDTO struct {
	Sort     string
	SortBy   string
	Keyword  string
	Page     int
	PageSize int
	Filters  map[string][]string
}

type IdentityReviewCampaignListResponseDTO struct {
	Items      []*IdentityReviewCampaignListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO     `json:"pagination"`
}

type IdentityReviewCampaignFormDTO struct {
	CampaignName       string                        `json:"campaign_name" binding:"required,max=250"`
	ReviewTypeID       uint                          `json:"review_type_id" binding:"required"`
	Description        *string                       `json:"description"`
	StartDate          string                        `json:"start_date" binding:"required,datetime=2006-01-02"`
	DueDate            string                        `json:"due_date" binding:"required,datetime=2006-01-02"`
	ReviewerID         uint                          `json:"reviewer_id" binding:"required"`
	ApproverID         uint                          `json:"approver_id" binding:"required"`
	ReminderDaysOffset *int                          `json:"reminder_days_offset"`
	AssetIDs           []uint                        `json:"asset_ids"`
	Documents          []*commonDTOs.DocumentFormDTO `json:"documents" binding:"omitempty"`
}

type IdentityReviewCampaignUserDTO struct {
	ID               uint    `json:"id"`
	FirstNameEn      string  `json:"first_name_en"`
	FirstNameAr      string  `json:"first_name_ar"`
	LastNameEn       string  `json:"last_name_en"`
	LastNameAr       string  `json:"last_name_ar"`
	Email            string  `json:"email"`
	Initial          string  `json:"initial"`
	Color            string  `json:"color"`
	ProfileImagePath *string `json:"profile_image_path"`
}

type IdentityReviewCampaignSystemRegisterDTO struct {
	ID      uint   `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
}

type IdentityReviewCampaignSystemDTO struct {
	ID               uint                                    `json:"id"`
	SystemRegisterID uint                                    `json:"system_register_id"`
	SystemRegister   IdentityReviewCampaignSystemRegisterDTO `json:"system_register"`
}

type IdentityReviewCampaignDocumentDTO struct {
	ID           uint   `json:"id"`
	URL          string `json:"url"`
	Name         string `json:"name"`
	Size         int    `json:"size"`
	Token        string `json:"token"`
	OriginalName string `json:"original_name"`
	Type         string `json:"type"`
}

type IdentityReviewCampaignDetailDTO struct {
	ID                   uint                                 `json:"id"`
	ReferenceNumber      string                               `json:"reference_number"`
	CampaignName         string                               `json:"campaign_name"`
	Description          *string                              `json:"description"`
	StartDate            string                               `json:"start_date" mapstructure:"-"`
	DueDate              string                               `json:"due_date" mapstructure:"-"`
	ReminderDaysOffset   *int                                 `json:"reminder_days_offset"`
	CompletionPercentage int                                  `json:"completion_percentage"`
	ReviewType           ReviewTypeDTO                        `json:"review_type"`
	CampaignStatus       ReviewCampaignStatusDTO              `json:"campaign_status"`
	Reviewer             IdentityReviewCampaignUserDTO        `json:"reviewer"`
	Approver             IdentityReviewCampaignUserDTO        `json:"approver"`
	Systems              []*IdentityReviewCampaignSystemDTO   `json:"systems"`
	Documents            []*IdentityReviewCampaignDocumentDTO `json:"documents"`
	CreatedBy            IdentityReviewCampaignUserDTO        `json:"created_by"`
	UpdatedBy            *IdentityReviewCampaignUserDTO       `json:"updated_by"`
	CreatedAt            time.Time                            `json:"created_at"`
	UpdatedAt            *time.Time                           `json:"updated_at"`
}

type IdentityReviewCampaignParsedDates struct {
	StartDate *time.Time
	DueDate   *time.Time
}
