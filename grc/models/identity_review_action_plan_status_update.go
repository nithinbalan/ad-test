package models

import "time"

type IdentityReviewActionPlanStatusUpdate struct {
	ID                         uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	IdentityReviewActionPlanID uint       `json:"identity_review_action_plan_id" gorm:"not null"`
	ActionPlanStatusID         uint       `json:"action_plan_status_id" gorm:"not null"`
	Comment                    string     `json:"comment" gorm:"type:nvarchar(max);not null"`
	Percentage                 *uint      `json:"percentage" gorm:"default:null"`
	CreatedAt                  time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                  *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedBYID                uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID                *uint      `json:"updated_by_id" gorm:"default:null"`

	IdentityReviewActionPlan *IdentityReviewActionPlan                       `json:"identity_review_action_plan,omitempty" gorm:"foreignKey:IdentityReviewActionPlanID"`
	ActionPlanStatus         *ActionPlanStatus                               `json:"action_plan_status" gorm:"foreignKey:ActionPlanStatusID"`
	CreatedByUser            *User                                           `json:"created_by" gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser            *User                                           `json:"updated_by" gorm:"foreignKey:UpdatedBYID"`
	Documents                []*IdentityReviewActionPlanStatusUpdateDocument `json:"documents" gorm:"foreignKey:IdentityReviewActionPlanStatusUpdateID"`
}

func (IdentityReviewActionPlanStatusUpdate) TableName() string {
	return "identity_review_action_plan_status_updates"
}
