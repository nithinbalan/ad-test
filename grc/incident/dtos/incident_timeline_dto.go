package dtos

import (
	"time"
)

type IncidentTimelineListDTO struct {
	ID                            uint      `json:"id"`
	IncidentRegisterID            uint      `json:"incident_register_id"`
	IncidentTimelineTypeTitleEn   string    `json:"incident_timeline_type_title_en"`
	IncidentTimelineTypeTitleAr   string    `json:"incident_timeline_type_title_ar"`
	IncidentTimelineTypeColorCode string    `json:"incident_timeline_type_color_code"`
	TimelineStatusTitleEn         string    `json:"timeline_status_title_en"`
	TimelineStatusTitleAr         string    `json:"timeline_status_title_ar"`
	TimelineStatusType            string    `json:"timeline_status_type"`
	TimelineStatusColorCode       string    `json:"timeline_status_color_code"`
	PerformedByFirstNameEn        string    `json:"performed_by_first_name_en"`
	PerformedByFirstNameAr        string    `json:"performed_by_first_name_ar"`
	PerformedByLastNameEn         string    `json:"performed_by_last_name_en"`
	PerformedByLastNameAr         string    `json:"performed_by_last_name_ar"`
	PerformedByEmail              string    `json:"performed_by_email"`
	PerformedByProfileImagePath   string    `json:"performed_by_profile_image_path"`
	PerformedAt                   time.Time `json:"performed_at"`
}

type IncidentTimelineListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	QuickFilter []string
}

type IncidentTimelineListResponseDTO struct {
	Items []*IncidentTimelineListDTO `json:"items"`
}

type IncidentTimelineFormDTO struct {
	ID                     uint       `json:"id,omitempty"`
	IncidentRegisterID     uint       `json:"incident_register_id" binding:"required"`
	IncidentTimelineTypeID uint       `json:"incident_timeline_type_id" binding:"required"`
	TimelineStatusID       uint       `json:"timeline_status_id" binding:"required"`
	PerformedByID          *uint      `json:"performed_by_id"`
	PerformedAt            *time.Time `json:"performed_at"`
}
