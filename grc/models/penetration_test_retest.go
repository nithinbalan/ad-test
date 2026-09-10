package models

import "time"

type PenetrationTestRetest struct {
	ID                        uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	ReferenceNumber           string     `json:"reference_number" gorm:"type:varchar(45);not null;unique"`
	PenetrationTestID         uint       `json:"penetration_test_id" gorm:"not null" binding:"required"`
	StartDate                 time.Time  `json:"start_date" gorm:"type:date;not null" binding:"required"`
	EndDate                   *time.Time `json:"end_date,omitempty" gorm:"type:date;default:null"`
	PenetrationTestExecutorID uint       `json:"penetration_test_executor_id" gorm:"not null" binding:"required"`
	PenetrationTestRetestStatusID uint   `json:"penetration_test_retest_status_id" gorm:"not null"`
	Notes                     *string    `json:"notes,omitempty" gorm:"type:nvarchar(max);default:null"`
	CreatedBYID               uint       `json:"created_by_id" gorm:"not null"`
	CreatedAt                 time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedBYID               *uint      `json:"updated_by_id,omitempty" gorm:"default:null"`
	UpdatedAt                 *time.Time `json:"updated_at,omitempty" gorm:"type:datetime;default:null"`

	// Relationships
	PenetrationTest         *PenetrationTest                `json:"penetration_test,omitempty" gorm:"foreignKey:PenetrationTestID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	PenetrationTestExecutor *PenetrationTestExecutor        `json:"penetration_test_executor,omitempty" gorm:"foreignKey:PenetrationTestExecutorID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	PenetrationTestRetestStatus *PenetrationTestRetestStatus `json:"penetration_test_retest_status,omitempty" gorm:"foreignKey:PenetrationTestRetestStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RetestFindings          []*PenetrationTestRetestFinding `json:"retest_findings,omitempty" gorm:"foreignKey:PenetrationTestRetestID"`
	CreatedByUser           *User                           `json:"created_by_user,omitempty" gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser           *User                           `json:"updated_by_user,omitempty" gorm:"foreignKey:UpdatedBYID"`
}

func (PenetrationTestRetest) TableName() string {
	return "penetration_test_retests"
}
