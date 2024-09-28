package dtos

type AuditReportWorkflowSubmitFormDTO struct {
	UserIDs []*uint `json:"user_ids" binding:"required"`
	Comment *string `json:"comment"`
}

type AuditReportWorkflowApproveFormDTO struct {
	Comment *string `json:"comment"`
}