package models

import (
	"time"
)

type ProjectDecisionLog struct {
	ID                uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ReferenceNumber   string     `json:"reference_number" gorm:"type:varchar(50);not null"`
	ProjectID         uint       `json:"project_id" gorm:"not null"`
	Decision          string     `json:"decision" gorm:"type:nvarchar(max);not null"`
	DecisionNotes     *string    `json:"decision_notes" gorm:"type:nvarchar(max);null"`
	Reason            *string    `json:"reason" gorm:"type:nvarchar(max);null"`
	RequestedByUserID uint       `json:"requested_by_user_id" gorm:"not null"`
	DecisionDate      time.Time  `json:"decision_date" gorm:"type:date;not null"`
	CreatedByID       uint       `json:"created_by_id" gorm:"not null"`
	UpdatedByID       *uint      `json:"updated_by_id" gorm:"null"`
	CreatedAt         time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt         *time.Time `json:"updated_at" gorm:"type:datetime;null"`

	// Relationships
	Project         *Project                      `gorm:"foreignKey:ProjectID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	RequestedByUser *User                         `gorm:"foreignKey:RequestedByUserID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser   *User                         `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:NO ACTION,OnDelete:SET NULL"`
	UpdatedByUser   *User                         `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:NO ACTION,OnDelete:SET NULL"`
	Documents       []*ProjectDecisionLogDocument `json:"documents" gorm:"foreignKey:ProjectDecisionLogID"`
}

func (ProjectDecisionLog) TableName() string {
	return "project_decision_logs"
}
