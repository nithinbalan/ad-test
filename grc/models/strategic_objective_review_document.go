package models

import "time"

type StrategicObjectiveReviewDocument struct {
	ID                         uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	StrategicObjectiveReviewID uint       `json:"strategic_objective_review_id" gorm:"not null"`
	URL                        string     `json:"url" gorm:"type:varchar(1000);not null"`
	Name                       string     `json:"name" gorm:"type:nvarchar(255);not null"`
	Size                       int        `json:"size" gorm:"not null"`
	Token                      string     `json:"token" gorm:"type:varchar(255);not null"`
	OriginalName               *string    `json:"original_name" gorm:"type:nvarchar(255);default:null"`
	Type                       string     `json:"type" gorm:"type:varchar(255);not null"`
	CreatedAt                  time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                  *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	StrategicObjectiveReview *StrategicObjectiveReview `gorm:"foreignKey:StrategicObjectiveReviewID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (StrategicObjectiveReviewDocument) TableName() string {
	return "strategic_objective_review_documents"
}
