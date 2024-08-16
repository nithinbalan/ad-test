package dtos

import (
	"time"
)

type RiskTimelineListDTO struct {
	ID                          uint      `json:"id"`
	RiskRegisterID              uint      `json:"risk_register_id"`
	RiskTimelineTypeTitleEn     string    `json:"risk_timeline_type_title_en"`
	RiskTimelineTypeTitleAr     string    `json:"risk_timeline_type_title_ar"`
	RiskTimelineTypeColorCode   string    `json:"risk_timeline_type_color_code"`
	TimelineStatusTitleEn       string    `json:"timeline_status_title_en"`
	TimelineStatusTitleAr       string    `json:"timeline_status_title_ar"`
	TimelineStatusType          string    `json:"timeline_status_type"`
	TimelineStatusColorCode     string    `json:"timeline_status_color_code"`
	PerformedByID               uint      `json:"performed_by_id"`
	PerformedByFirstNameEn      string    `json:"performed_by_first_name_en"`
	PerformedByFirstNameAr      string    `json:"performed_by_first_name_ar"`
	PerformedByLastNameEn       string    `json:"performed_by_last_name_en"`
	PerformedByLastNameAr       string    `json:"performed_by_last_name_ar"`
	PerformedByEmail            string    `json:"performed_by_email"`
	PerformedByInitial          string    `json:"performed_by_initial"`
	PerformedByColor            string    `json:"performed_by_color"`
	PerformedByProfileImagePath string    `json:"performed_by_profile_image_path"`
	PerformedAt                 time.Time `json:"performed_at"`
}

type RiskTimelineListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	QuickFilter []string
}

type RiskTimelineListResponseDTO struct {
	Items []*RiskTimelineListDTO `json:"items"`
}

type RiskTimelineFormDTO struct {
	ID                 uint       `json:"id,omitempty"`
	RiskRegisterID     uint       `json:"risk_register_id" binding:"required"`
	RiskTimelineTypeID uint       `json:"risk_timeline_type_id" binding:"required"`
	TimelineStatusID   uint       `json:"timeline_status_id" binding:"required"`
	PerformedByID      *uint      `json:"performed_by_id"`
	PerformedAt        *time.Time `json:"performed_at"`
}
