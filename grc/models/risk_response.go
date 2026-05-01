package models

import "time"

type RiskResponse struct {
	ID               uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	RiskRegisterID   uint       `json:"risk_register_id" binding:"required" gorm:"not null"`
	WorkflowStatusID *uint      `json:"workflow_status_id" gorm:"default:null"`
	NextReviewLevel  *int       `json:"next_review_level" gorm:"default:null"`
	SubmittedByID    *uint      `json:"submitted_by_id" gorm:"default:null"`
	SubmittedAt      *time.Time `json:"submitted_at" gorm:"type:datetime;default:null"`
	ApprovedByID     *uint      `json:"approved_by_id" gorm:"default:null"`
	ApprovedAt       *time.Time `json:"approved_at" gorm:"type:datetime;default:null"`
	CreatedBYID      *uint      `json:"created_by_id" gorm:"default:null"`
	UpdatedBYID      *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt        time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	UpdatedAt        *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relations
	RiskRegister               *RiskRegister                `gorm:"foreignKey:RiskRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	WorkflowStatus             *WorkflowStatus              `gorm:"foreignKey:WorkflowStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	SubmittedBy                *User                        `gorm:"foreignKey:SubmittedByID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ApprovedBy                 *User                        `gorm:"foreignKey:ApprovedByID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedBy                  *User                        `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	UpdatedBy                  *User                        `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RiskResponseWorkflows      []*RiskResponseWorkflow      `gorm:"foreignKey:RiskResponseID"`
	RiskResponseReviewComments []*RiskResponseReviewComment `gorm:"foreignKey:RiskResponseID"`
}

func (RiskResponse) TableName() string {
	return "risk_responses"
}
