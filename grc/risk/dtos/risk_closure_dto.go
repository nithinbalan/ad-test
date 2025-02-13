package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type RiskClosureDetailsDTO struct {
	ID             uint                             `json:"id"`
	RiskRegisterID uint                             `json:"risk_register_id"`
	Justification  string                           `json:"justification"`
	LessonLearned  *string                          `json:"lesson_learned"`
	WhatWorkedWell *string                          `json:"what_worked_well"`
	Improvements   *string                          `json:"improvements"`
	Takeaways      *string                          `json:"takeaways"`
	CreatedAt      time.Time                        `json:"created_at"`
	CreatedByUser  organizationDTOs.UserMinimalDTO  `json:"created_by"`
	UpdatedAt      *time.Time                       `json:"updated_at,omitempty"`
	UpdatedByUser  *organizationDTOs.UserMinimalDTO `json:"updated_by"`
	Documents      []*RiskClosureDocumentDTO        `json:"documents"`
}

type RiskClosureDocumentDTO struct {
	ID            uint   `json:"id"`
	RiskClosureID uint   `json:"risk_closure_id"`
	URL           string `json:"url"`
	Name          string `json:"name"`
	Size          int    `json:"size"`
	Token         string `json:"token"`
	OriginalName  string `json:"original_name"`
	Type          string `json:"type"`
}

type RiskClosureFormDTO struct {
	Justification  string                        `json:"justification" binding:"required"`
	LessonLearned  *string                       `json:"lesson_learned"`
	WhatWorkedWell *string                       `json:"what_worked_well"`
	Improvements   *string                       `json:"improvements"`
	Takeaways      *string                       `json:"takeaways"`
	Documents      []*commonDTOs.DocumentFormDTO `json:"documents"`
}
