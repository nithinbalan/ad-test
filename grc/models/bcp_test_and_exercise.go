package models

import (
	"time"
)

type BcpTestAndExercise struct {
	ID                       uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	BusinessContinuityPlanID uint       `json:"business_continuity_id" binding:"required"`
	Title                    string     `json:"title"`
	Objective                *string    `json:"objective"`
	Scenario                 *string    `json:"scenario"`
	SuccessCriteria          *string    `json:"success_criteria"`
	ResponsibleUserID        *uint      `json:"responsible_user_id" gorm:"not null" binding:"required"`
	ResponsibleUnitID        *uint      `json:"responsible_unit_id" gorm:"not null" binding:"required"`
	ReviewFrequencyID        *uint      `json:"review_frequency_id"`
	CreatedBYID              uint       `json:"created_by_id" binding:"required"`
	UpdatedBYID              *uint      `json:"updated_by_id"`
	CreatedAt                time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	BusinessContinuityPlan         *BusinessContinuityPlan         `gorm:"foreignKey:BusinessContinuityPlanID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResponsibleUser                *User                           `gorm:"foreignKey:ResponsibleUserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResponsibleUnit                *OrganizationUnit               `gorm:"foreignKey:ResponsibleUnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ReviewFrequency                *ReviewFrequency                `gorm:"foreignKey:ReviewFrequencyID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	BcpTestAndExerciseBcpProcesses []*BcpTestAndExerciseBcpProcess `gorm:"foreignKey:BcpTestAndExerciseID"`
	CreatedByUser                  *User                           `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser                  *User                           `gorm:"foreignKey:UpdatedBYID"`
}

func (BcpTestAndExercise) TableName() string {
	return "bcp_test_and_exercises"
}
