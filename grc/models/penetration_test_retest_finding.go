package models

import "time"

type PenetrationTestRetestFinding struct {
	ID                       uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	PenetrationTestRetestID  uint       `json:"penetration_test_retest_id" gorm:"not null" binding:"required"`
	PenetrationTestFindingID uint       `json:"penetration_test_finding_id" gorm:"not null" binding:"required"`
	Comment                  *string    `json:"comment,omitempty" gorm:"type:nvarchar(MAX);default:null"`
	Result                   *string    `json:"result,omitempty" gorm:"type:varchar(20);default:null"`
	ResultAt                 *time.Time `json:"result_at,omitempty" gorm:"type:datetime;default:null"`
	CreatedBYID              uint       `json:"created_by_id" gorm:"not null"`
	CreatedAt                time.Time  `json:"created_at" gorm:"type:datetime;not null"`

	PenetrationTestRetest  *PenetrationTestRetest                  `json:"penetration_test_retest,omitempty" gorm:"foreignKey:PenetrationTestRetestID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	PenetrationTestFinding *PenetrationTestFinding                 `json:"penetration_test_finding,omitempty" gorm:"foreignKey:PenetrationTestFindingID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser          *User                                   `json:"created_by_user,omitempty" gorm:"foreignKey:CreatedBYID"`
	Documents              []*PenetrationTestRetestFindingDocument `json:"documents,omitempty" gorm:"foreignKey:PenetrationTestRetestFindingID"`
}

func (PenetrationTestRetestFinding) TableName() string {
	return "penetration_test_retest_findings"
}
