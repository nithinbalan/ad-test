package dtos

import (
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type ComplianceAssessmentRequirementDetailDTO struct {
	ID                    int                              `json:"id"`
	ComplianceAssessment  *ComplianceAssessmentDetailDTO   `json:"compliance_assessment"`
	ComplianceRequirement ComplianceRequirementDetailDTO   `json:"compliance_requirement"`
	IsCompleted           bool                             `json:"is_completed"`
	Score                 *int                             `json:"score"`
	Comments              *string                          `json:"comments"`
	ComplianceStatusID    *uint                            `json:"compliance_status_id"`
	ComplianceStatus      *ComplianceStatusDTO             `json:"compliance_status"`
	CreatedAt             time.Time                        `json:"created_at"`
	UpdatedAt             *time.Time                       `json:"updated_at"`
	CreatedByUser         organizationDTOs.UserMinimalDTO  `json:"created_by_user"`
	UpdatedByUser         *organizationDTOs.UserMinimalDTO `json:"updated_by_user"`
}

type ComplianceAssessmentRequirementUpdateFormDTO struct {
	Score              int     `json:"score" binding:"required"`
	ComplianceStatusID uint    `json:"compliance_status_id" binding:"required"`
	Comments           *string `json:"comments"`
}
