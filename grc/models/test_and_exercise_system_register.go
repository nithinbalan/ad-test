package models

import (
	"time"
)

type TestAndExerciseSystemRegister struct {
	ID                uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TestAndExerciseID uint       `json:"test_and_exercise_id" binding:"required"`
	SystemRegisterID  uint       `json:"system_register_id" binding:"required"`
	PlannedRto        *int       `json:"planned_rto"`
	PlannedRtoUnitId  *uint      `json:"planned_rto_unit_id"`
	PlannedMtpd       *int       `json:"planned_mtpd"`
	PlannedMtpdUnitId *uint      `json:"planned_mtpd_unit_id"`
	PlannedRpo        *int       `json:"planned_rpo"`
	PlannedRpoUnitId  *uint      `json:"planned_rpo_unit_id"`
	ActualRto         *int       `json:"actual_rto"`
	ActualRtoUnitId   *uint      `json:"actual_rto_unit_id"`
	ActualMtpd        *int       `json:"actual_mtpd"`
	ActualMtpdUnitId  *uint      `json:"actual_mtpd_unit_id"`
	ActualRpo         *int       `json:"actual_rpo"`
	ActualRpoUnitId   *uint      `json:"actual_rpo_unit_id"`
	Comments          *string    `json:"comments"`
	TimelinessID      *uint      `json:"timeliness_id"`
	IsRecovered       bool       `json:"is_recovered"`
	CreatedBYID       uint       `json:"created_by_id" binding:"required"`
	CreatedAt         time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedBYID       *uint      `json:"updated_by_id"`
	UpdatedAt         *time.Time `json:"updated_at"`

	// Relationships
	TestAndExercise *TestAndExercise                `gorm:"foreignKey:TestAndExerciseID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	SystemRegister  *SystemRegister                 `gorm:"foreignKey:SystemRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	PlannedRtoUnit  *BusinessContinuityDurationUnit `gorm:"foreignKey:PlannedRtoUnitId;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	PlannedMtpdUnit *BusinessContinuityDurationUnit `gorm:"foreignKey:PlannedMtpdUnitId;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	PlannedRpoUnit  *BusinessContinuityDurationUnit `gorm:"foreignKey:PlannedRpoUnitId;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ActualRtoUnit   *BusinessContinuityDurationUnit `gorm:"foreignKey:ActualRtoUnitId;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ActualMtpdUnit  *BusinessContinuityDurationUnit `gorm:"foreignKey:ActualMtpdUnitId;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ActualRpoUnit   *BusinessContinuityDurationUnit `gorm:"foreignKey:ActualRpoUnitId;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Timeliness      *Timeliness                     `gorm:"foreignKey:TimelinessID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser   *User                           `gorm:"foreignKey:CreatedBYID"`
}

func (TestAndExerciseSystemRegister) TableName() string {
	return "test_and_exercise_system_registers"
}
