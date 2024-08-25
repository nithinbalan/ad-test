package dtos

import (
	"time"
)

type AuditTimelineListDTO struct {
	ID                          uint      `json:"id"`
	AuditID                     uint      `json:"audit_id"`
	AuditTimelineTypeTitleEn    string    `json:"audit_timeline_type_title_en"`
	AuditTimelineTypeTitleAr    string    `json:"audit_timeline_type_title_ar"`
	AuditTimelineTypeColorCode  string    `json:"audit_timeline_type_color_code"`
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

type AuditTimelineListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	QuickFilter []string
}

type AuditTimelineListResponseDTO struct {
	Items []*AuditTimelineListDTO `json:"items"`
}

type AuditTimelineFormDTO struct {
	ID                  uint       `json:"id,omitempty"`
	AuditID             uint       `json:"audit_id" binding:"required"`
	AuditTimelineTypeID uint       `json:"audit_timeline_type_id" binding:"required"`
	TimelineStatusID    uint       `json:"timeline_status_id" binding:"required"`
	PerformedByID       *uint      `json:"performed_by_id"`
	PerformedAt         *time.Time `json:"performed_at"`
}
