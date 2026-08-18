package models

import (
	"time"
)

type AuditFindingReassessment struct {
	ID             uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	AuditFindingID uint       `json:"audit_finding_id" binding:"required" gorm:"not null"`
	Outcome        *string    `json:"outcome" gorm:"type:varchar(20);default:null"`
	Comment        *string    `json:"comment" gorm:"type:nvarchar(max);default:null"`
	CreatedBYID    uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID    *uint      `json:"updated_by_id,omitempty" gorm:"default:null"`
	CreatedAt      time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty" gorm:"type:datetime;default:null"`

	// Relationships
	AuditFinding            *AuditFinding                       `gorm:"foreignKey:AuditFindingID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser           *User                               `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser           *User                               `gorm:"foreignKey:UpdatedBYID"`
	Documents               []*AuditFindingReassessmentDocument `gorm:"foreignKey:AuditFindingReassessmentID"`
	AuditFindingActionPlans []*AuditFindingActionPlan           `gorm:"foreignKey:AuditFindingReassessmentID"`
}

func (AuditFindingReassessment) TableName() string {
	return "audit_finding_reassessments"
}
