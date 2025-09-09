package dtos

type RAGSettingDTO struct {
	ID                   uint   `json:"id"`
	ThresholdDays        int    `json:"threshold_days"`
	OnTrackColorCode     string `json:"on_track_color_code"`
	ApproachingColorCode string `json:"approaching_color_code"`
	OverdueColorCode     string `json:"overdue_color_code"`
}
