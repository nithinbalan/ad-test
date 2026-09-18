package models

import "time"

type PenetrationTestReportRetestFinding struct {
	ID                             uint       `json:"id" gorm:"primary_key;auto_increment"`
	PenetrationTestReportRetestID  uint       `json:"penetration_test_report_retest_id" gorm:"not null"`
	PenetrationTestReportFindingID uint       `json:"penetration_test_report_finding_id" gorm:"not null"`
	PenetrationTestRetestFindingID uint       `json:"penetration_test_retest_finding_id" gorm:"not null"`
	Comment                        *string    `json:"comment" gorm:"type:nvarchar(max)"`
	CommentAt                      *time.Time `json:"comment_at" gorm:"type:datetime"`
	ResultTitleEn                  *string    `json:"result_title_en" gorm:"type:nvarchar(255)"`
	ResultTitleAr                  *string    `json:"result_title_ar" gorm:"type:nvarchar(255)"`
	ResultAt                       *time.Time `json:"result_at" gorm:"type:datetime"`

	// Relationships
	PenetrationTestReportRetest  *PenetrationTestReportRetest     `gorm:"foreignKey:PenetrationTestReportRetestID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	PenetrationTestReportFinding *PenetrationTestReportFinding    `gorm:"foreignKey:PenetrationTestReportFindingID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	PenetrationTestRetestFinding *PenetrationTestRetestFinding    `gorm:"foreignKey:PenetrationTestRetestFindingID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Documents                    []*PenetrationTestReportDocument `gorm:"foreignKey:PenetrationTestReportRetestFindingID"`
}

func (PenetrationTestReportRetestFinding) TableName() string {
	return "penetration_test_report_retest_findings"
}
