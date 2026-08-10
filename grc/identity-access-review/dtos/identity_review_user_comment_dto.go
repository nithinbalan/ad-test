package dtos

type IdentityReviewUserCommentFormDTO struct {
	AccountDecisionID *uint  `json:"account_decision_id" binding:"omitempty"`
	Comment           string `json:"comment" binding:"required"`
}
