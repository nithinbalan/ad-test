package dtos

type RiskCommentFormDTO struct {
	RiskTimelineTypeID uint   `json:"risk_timeline_type_id"`
	Comment            string `json:"comment" binding:"required"`
}
