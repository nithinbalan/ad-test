package models

import "time"

type PenetrationTestExecution struct {
	ID                uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	PenetrationTestID uint       `json:"penetration_test_id" gorm:"not null;unique"`
	TestStartDate     time.Time  `json:"test_start_date" gorm:"type:date;not null" binding:"required"`
	TestEndDate       time.Time  `json:"test_end_date" gorm:"type:date;not null" binding:"required"`
	ExecutionStatusID uint       `json:"execution_status_id" gorm:"not null" binding:"required"`
	SecurityResultID  uint       `json:"security_result_id" gorm:"not null" binding:"required"`
	CreatedBYID       uint       `json:"created_by_id" gorm:"not null"`
	CreatedAt         time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedBYID       *uint      `json:"updated_by_id,omitempty" gorm:"default:null"`
	UpdatedAt         *time.Time `json:"updated_at,omitempty" gorm:"type:datetime;default:null"`

	ExecutionStatus *PenetrationTestExecutionStatus `json:"execution_status,omitempty" gorm:"foreignKey:ExecutionStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	SecurityResult  *PenetrationTestSecurityResult  `json:"security_result,omitempty" gorm:"foreignKey:SecurityResultID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser   *User                           `json:"created_by_user,omitempty" gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser   *User                           `json:"updated_by_user,omitempty" gorm:"foreignKey:UpdatedBYID"`
}

func (PenetrationTestExecution) TableName() string {
	return "penetration_test_executions"
}
