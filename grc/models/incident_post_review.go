package models

import (
	"time"
)

type IncidentPostReview struct {
	ID                                uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	IncidentRegisterID                uint       `json:"incident_register_id" binding:"required" gorm:"not null"`
	IncidentPostReviewEffectivenessID int        `json:"incident_post_review_effectiveness_id" binding:"required" gorm:"not null"`
	PostReviewDate                    time.Time  `json:"post_review_date" binding:"required" gorm:"type:date;not null"`
	Summary                           *string    `json:"summary" gorm:"type:nvarchar(max)"`
	CreatedBy                         uint       `json:"created_by_id" binding:"required"`
	UpdatedBy                         *uint      `json:"updated_by_id"`
	CreatedAt                         time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                         *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	IncidentRegister                IncidentRegister                    `gorm:"foreignkey:IncidentRegisterID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	IncidentPostReviewEffectiveness IncidentPostReviewEffectiveness     `gorm:"foreignkey:IncidentPostReviewEffectivenessID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Users                           []*IncidentPostReviewUser           `gorm:"foreignKey:IncidentPostReviewID"`
	Lessons                         []*IncidentPostReviewLesson         `gorm:"foreignKey:IncidentPostReviewID"`
	Recommendations                 []*IncidentPostReviewRecommendation `gorm:"foreignKey:IncidentPostReviewID"`
	CreatedByUser                   User                                `gorm:"foreignkey:CreatedBy"`
	UpdatedByUser                   *User                               `gorm:"foreignkey:UpdatedBy"`
}

// TableName overrides the table name used by IncidentPostReview to `incident_post_reviews`
func (IncidentPostReview) TableName() string {
	return "incident_post_reviews"
}
