package dtos

import (
	complianceDTOs "jethur/grc/compliance/dtos"
	documentsDTOs "jethur/grc/documents/dtos"

	"time"
)

// AssetPolicyChecklistDTO represents the data structure for API responses
type AssetPolicyChecklistDTO struct {
	ID                         uint                                      `json:"id"`
	AssetPolicyID              uint                                      `json:"asset_policy_id"`
	DocumentRequirementGroupID uint                                      `json:"document_requirement_group_id"`
	DocumentRequirementGroup   documentsDTOs.DocumentRequirementGroupDTO `json:"document_requirement_group"`
	QuestionGroupTitle         string                                    `json:"question_group_title"`
	QuestionTitle              string                                    `json:"question_title"`
	Justification              *string                                   `json:"justification"`
	ComplianceStatusID         *uint                                     `json:"compliance_status_id"`
	ComplianceStatus           *complianceDTOs.ComplianceStatusDTO       `json:"compliance_status"`
	AnsweredBYID               *uint                                     `json:"answered_by_id"`
	AnsweredAt                 *time.Time                                `json:"answered_at"`
	CreatedAt                  time.Time                                 `json:"created_at"`
	UpdatedAt                  *time.Time                                `json:"updated_at"`
	CreatedBYID                uint                                      `json:"created_by_id"`
	UpdatedBYID                *uint                                     `json:"updated_by_id"`
}

// AssetPolicyChecklistFormDTO is used for creating or updating a checklist
type AssetPolicyChecklistAnswerFormDTO struct {
	Justification      *string `json:"justification,omitempty"`
	ComplianceStatusID uint    `json:"compliance_status_id" binding:"required"`
}

type AssetPolicyChecklistFormDTO struct {
	DocumentRequirementGroupID uint   `json:"document_requirement_group_id" binding:"required"`
	QuestionTitle              string `json:"question_title" binding:"required,max=500"`
}

type AssetPolicyChecklistUpdateDTO struct {
	QuestionTitle string `json:"question_title" binding:"required,max=500"`
}

type AssetPolicyChecklistSummaryDTO struct {
	Summary          ChecklistSummaryDTO                 `json:"summary"`
	ComplianceStatus *complianceDTOs.ComplianceStatusDTO `json:"compliance_status"`
}

type ChecklistSummaryDTO struct {
	TotalChecklist     int `json:"total_checklist"`
	Compliant          int `json:"compliant"`
	PartiallyCompliant int `json:"partially_compliant"`
	NotApplicable      int `json:"not_applicable"`
	NonCompliant       int `json:"non_compliant"`
}
