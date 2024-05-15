package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type SLACheckAnswerDetailDTO struct {
	ID                      uint                            `json:"id"`
	SLACheckQuestionTitleEn string                          `json:"sla_check_question_title_en"`
	SLACheckQuestionTitleAr string                          `json:"sla_check_question_title_ar"`
	ComplianceStatus        ComplianceStatusDTO             `json:"compliance_status"`
	Comment                 string                          `json:"comment"`
	CreatedAt               time.Time                       `json:"created_at"`
	UpdatedAt               *time.Time                      `json:"updated_at"`
	CreatedByUser           organizationDTOs.UserMinimalDTO `json:"created_by"`
	UpdatedByUser           organizationDTOs.UserMinimalDTO `json:"updated_by"`
	Documents               []*SLACheckAnswerDocumentDTO    `json:"documents"`
}

type SLACheckAnswerDocumentDTO struct {
	ID    *uint  `json:"id,omitempty"`
	URL   string `json:"url"`
	Icon  string `json:"icon"`
	Name  string `json:"name"`
	Token string `json:"token"`
	Type  string `json:"type"`
	Size  int    `json:"size"`
}

type SLACheckAnswerFormDTO struct {
	ComplianceStatusID uint                          `json:"compliance_status_id" binding:"required"`
	Comment            string                        `json:"comment" binding:"required"`
	Documents          []*commonDTOs.DocumentFormDTO `json:"documents" binding:"omitempty"`
}

type AssessmentSummaryDTO struct {
	TotalQuestions int          `json:"totalQuestions"`
	TotalAnswers   int          `json:"totalAnswers"`
	AnswerCounts   map[uint]int `json:"answerCounts"`
}
