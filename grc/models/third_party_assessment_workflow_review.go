package models

import (
	"time"
)

type ThirdPartyAssessmentWorkflowReview struct {
	ID                       uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TprmAssessmentWorkflowID uint       `json:"tprm_assessment_workflow_id" binding:"required"`
	WorkflowStatusID         uint       `json:"workflow_status_id" binding:"required" gorm:"not null"`
	Comment                  string     `json:"comment" binding:"-" gorm:"type:text"`
	CreatedBYID              uint       `json:"created_by_id" binding:"required" gorm:"not null"`
	UpdatedBYID              *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt                time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	UpdatedAt                *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relations
	TprmAssessmentWorkflow *ThirdPartyAssessmentWorkflow `gorm:"foreignKey:TprmAssessmentWorkflowID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	WorkflowStatus         *WorkflowStatus               `gorm:"foreignKey:WorkflowStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser          *User                         `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	UpdatedByUser          *User                         `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (ThirdPartyAssessmentWorkflowReview) TableName() string {
	return "tprm_assessment_workflow_reviews"
}
