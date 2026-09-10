package models

import "time"

type PenetrationTestRetestFinding struct {
	ID                       uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	PenetrationTestRetestID  uint       `json:"penetration_test_retest_id" gorm:"not null" binding:"required"`
	PenetrationTestFindingID uint       `json:"penetration_test_finding_id" gorm:"not null" binding:"required"`
	Comment                  *string    `json:"comment" gorm:"type:nvarchar(MAX);default:null"`
	RetestResultID           *uint      `json:"retest_result_id" gorm:"default:null"`
	ResultAt                 *time.Time `json:"result_at" gorm:"type:datetime;default:null"`
	CommentBYID              *uint      `json:"comment_by_id" gorm:"default:null"`
	CommentAt                *time.Time `json:"comment_at" gorm:"type:datetime;default:null"`
	ResultBYID               *uint      `json:"result_by_id" gorm:"default:null"`

	PenetrationTestRetest  *PenetrationTestRetest                  `json:"penetration_test_retest,omitempty" gorm:"foreignKey:PenetrationTestRetestID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	PenetrationTestFinding *PenetrationTestFinding                 `json:"penetration_test_finding,omitempty" gorm:"foreignKey:PenetrationTestFindingID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RetestResult           *PenetrationTestRetestResult            `json:"retest_result,omitempty" gorm:"foreignKey:RetestResultID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CommentByUser          *User                                   `json:"comment_by_user,omitempty" gorm:"foreignKey:CommentBYID"`
	ResultByUser           *User                                   `json:"result_by_user,omitempty" gorm:"foreignKey:ResultBYID"`
	Documents              []*PenetrationTestRetestFindingDocument `json:"documents,omitempty" gorm:"foreignKey:PenetrationTestRetestFindingID"`
}

func (PenetrationTestRetestFinding) TableName() string {
	return "penetration_test_retest_findings"
}
