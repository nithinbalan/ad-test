package models

import "time"

type IdentityReviewActionPlanStatusUpdateDocument struct {
	ID                                     uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	IdentityReviewActionPlanStatusUpdateID uint       `json:"identity_review_action_plan_status_update_id" gorm:"not null"`
	URL                                    string     `json:"url" gorm:"type:varchar(1000);not null"`
	Name                                   string     `json:"name" gorm:"type:nvarchar(255);not null"`
	Size                                   int        `json:"size" gorm:"not null"`
	Token                                  string     `json:"token" gorm:"type:varchar(255);not null"`
	OriginalName                           string     `json:"original_name" gorm:"type:nvarchar(255)"`
	Type                                   string     `json:"type" gorm:"type:varchar(255);not null"`
	CreatedByID                            uint       `json:"created_by_id" gorm:"not null"`
	UpdatedByID                            *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt                              time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                              *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
}

func (IdentityReviewActionPlanStatusUpdateDocument) TableName() string {
	return "identity_review_action_plan_status_update_documents"
}
