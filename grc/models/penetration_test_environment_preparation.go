package models

import "time"

type PenetrationTestEnvironmentPreparation struct {
	ID                     uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	PenetrationTestID      uint       `json:"penetration_test_id" gorm:"not null;unique"`
	EnvironmentNotifyDate  time.Time  `json:"environment_notify_date" gorm:"type:date;not null" binding:"required"`
	ExpectedReadyDate      time.Time  `json:"expected_ready_date" gorm:"type:date;not null" binding:"required"`
	EnvironmentReadinessID uint       `json:"environment_readiness_id" gorm:"not null" binding:"required"`
	CreatedBYID            uint       `json:"created_by_id" gorm:"not null"`
	CreatedAt              time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedBYID            *uint      `json:"updated_by_id,omitempty" gorm:"default:null"`
	UpdatedAt              *time.Time `json:"updated_at,omitempty" gorm:"type:datetime;default:null"`

	EnvironmentReadiness *PenetrationTestEnvironmentReadiness `json:"environment_readiness,omitempty" gorm:"foreignKey:EnvironmentReadinessID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser        *User                                `json:"created_by_user,omitempty" gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser        *User                                `json:"updated_by_user,omitempty" gorm:"foreignKey:UpdatedBYID"`
}

func (PenetrationTestEnvironmentPreparation) TableName() string {
	return "penetration_test_environment_preparations"
}
