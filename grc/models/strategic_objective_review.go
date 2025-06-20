package models

import "time"

type StrategicObjectiveReview struct {
	ID                   uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	StrategicObjectiveID uint       `json:"strategic_objective_id" gorm:"not null"`
	ReviewDate           time.Time  `json:"review_date" gorm:"type:date;not null"`
	Description          *string    `json:"description" gorm:"type:nvarchar(max);default:null"`
	CreatedBYID          uint       `json:"created_by_id" gorm:"not null"`
	CreatedAt            time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedByID          *uint      `json:"updated_by_id" gorm:"default:null"`
	UpdatedAt            *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	StrategicObjective *StrategicObjective                 `gorm:"foreignKey:StrategicObjectiveID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedByUser      *User                               `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser      *User                               `gorm:"foreignKey:UpdatedByID"`
	Documents          []*StrategicObjectiveReviewDocument `gorm:"foreignKey:StrategicObjectiveReviewID"`
}

func (StrategicObjectiveReview) TableName() string {
	return "strategic_objective_reviews"
}
