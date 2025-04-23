package models

import (
	"time"
)

type RiskClosure struct {
	ID             uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	RiskRegisterID uint       `json:"risk_register_id" binding:"required"`
	Justification  string     `json:"justification" gorm:"type:nvarchar(max)" binding:"required"`
	LessonLearned  *string    `json:"lesson_learned" gorm:"type:nvarchar(max)"`
	WhatWorkedWell *string    `json:"what_worked_well" gorm:"type:nvarchar(max)"`
	Improvements   *string    `json:"improvements" gorm:"type:nvarchar(max)"`
	Takeaways      *string    `json:"takeaways" gorm:"type:nvarchar(max)"`
	CreatedAt      time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	CreatedBYID    uint       `json:"created_by_id" binding:"required"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty" gorm:"type:datetime;default:null"`
	UpdatedBYID    *uint      `json:"updated_by_id"`

	// Relationships
	RiskRegister  *RiskRegister          `gorm:"foreignKey:RiskRegisterID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser *User                  `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedByUser *User                  `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Documents     []*RiskClosureDocument `gorm:"foreignKey:RiskClosureID"`
}

// TableName overrides the table name used by RiskClosure to `risk_closures`
func (RiskClosure) TableName() string {
	return "risk_closures"
}
