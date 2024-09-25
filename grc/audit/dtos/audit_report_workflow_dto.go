package dtos

type AuditReportWorkflowFormDTO struct {
	UserIDs []*uint `json:"user_ids" binding:"required"`
	Comment *string `json:"comment"`
}
