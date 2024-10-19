package dtos

import (
	organizationDTOs "jethur/grc/organization/dtos"
)

type AuditReportWorkflowParticipantDTO struct {
	ID                    uint                            `json:"id"`
	AuditReportWorkflowID uint                            `json:"audit_report_workflow_id"`
	UserID                uint                            `json:"user_id"`
	User                  organizationDTOs.UserMinimalDTO `json:"user"`
}
