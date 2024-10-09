package dtos

import "time"

type AuditInsightsDraftFindingApproveFormDTO struct {
	Comment           string  `json:"comment" binding:"required"`
	Title             string  `json:"title" binding:"required"`
	Type              string  `json:"type" binding:"required,oneof=DR"`
	Description       *string `json:"description" binding:"omitempty"`
	ResponsibleUserID uint    `json:"responsible_user_id" binding:"required"`
	TargetDate        *string `json:"target_date" binding:"omitempty,datetime=2006-01-02"`
}

type AuditInsightsDraftFindingActionPlanParsedDates struct {
	TargetDate *time.Time
}