package models

import "time"

type RiskUpdatePlanReviewComment struct {
	ID                       uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	RiskRegisterUpdatePlanID uint      `json:"risk_register_update_plan_id" gorm:"not null"`
	Comment                  string    `json:"comment" binding:"required" gorm:"type:nvarchar(max)"`
	Level                    *int      `json:"level"`
	CommentStatusID          uint      `json:"comment_status_id" gorm:"not null"`
	IsEdited                 bool      `json:"is_edited" gorm:"not null;default:false"`
	ParentCommentID          *uint     `json:"risk_update_plan_review_comment_id"`
	CreatedBYID              uint      `json:"created_by_id" binding:"required" gorm:"not null"`
	UpdatedBYID              *uint     `json:"updated_by_id"`
	CreatedAt                time.Time `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	RiskRegisterUpdatePlan *RiskRegisterUpdatePlan               `gorm:"foreignKey:RiskRegisterUpdatePlanID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ParentComment           *RiskUpdatePlanReviewComment          `gorm:"foreignKey:ParentCommentID"`
	ChildComments           []*RiskUpdatePlanReviewComment        `gorm:"foreignKey:ParentCommentID"`
	CommentStatus           *CommentStatus                        `gorm:"foreignKey:CommentStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser           *User                                 `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser           *User                                 `gorm:"foreignKey:UpdatedBYID"`
	Documents               []*RiskUpdatePlanReviewCommentDocument `gorm:"foreignKey:RiskUpdatePlanReviewCommentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (RiskUpdatePlanReviewComment) TableName() string {
	return "risk_update_plan_review_comments"
}
