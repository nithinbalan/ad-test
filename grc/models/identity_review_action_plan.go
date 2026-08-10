package models

import "time"

type IdentityReviewActionPlan struct {
	ID                   uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	IdentityReviewUserID uint       `json:"identity_review_user_id" gorm:"not null"`
	Title                string     `json:"title" gorm:"type:nvarchar(500);not null"`
	Description          *string    `json:"description" gorm:"type:nvarchar(max);default:null"`
	SeverityLevelID      *uint      `json:"severity_level_id" gorm:"default:null"`
	TargetDate           *time.Time `json:"target_date" gorm:"type:date;default:null"`
	ResponsibleUserID    uint       `json:"responsible_user_id" gorm:"not null"`
	ReviewUserID         *uint      `json:"review_user_id" gorm:"default:null"`
	ActionPlanStatusID   uint       `json:"action_plan_status_id" gorm:"not null"`
	Percentage           uint       `json:"percentage" gorm:"not null;default:0"`
	Comments             *string    `json:"comments" gorm:"type:nvarchar(max);default:null"`
	DelegateeUserID      *uint      `json:"delegatee_user_id" gorm:"default:null"`
	CreatedAt            time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt            *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedByID          uint       `json:"created_by_id" gorm:"not null"`
	UpdatedByID          *uint      `json:"updated_by_id" gorm:"default:null"`

	IdentityReviewUser      *IdentityReviewUser                     `json:"identity_review_user" gorm:"foreignKey:IdentityReviewUserID"`
	SeverityLevel           *SeverityLevel                          `json:"severity_level" gorm:"foreignKey:SeverityLevelID"`
	ResponsibleUser         *User                                   `json:"responsible_user" gorm:"foreignKey:ResponsibleUserID"`
	ReviewUser              *User                                   `json:"review_user" gorm:"foreignKey:ReviewUserID"`
	ActionPlanStatus        *ActionPlanStatus                       `json:"action_plan_status" gorm:"foreignKey:ActionPlanStatusID"`
	DelegateeUser           *User                                   `json:"delegatee_user" gorm:"foreignKey:DelegateeUserID"`
	CreatedByUser           *User                                   `json:"created_by" gorm:"foreignKey:CreatedByID"`
	UpdatedByUser           *User                                   `json:"updated_by" gorm:"foreignKey:UpdatedByID"`
	ActionPlanStatusUpdates []*IdentityReviewActionPlanStatusUpdate `json:"action_plan_status_updates" gorm:"foreignKey:IdentityReviewActionPlanID"`
}

func (IdentityReviewActionPlan) TableName() string {
	return "identity_review_action_plans"
}
