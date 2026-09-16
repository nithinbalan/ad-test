package models

import "time"

type PenetrationTestReportRetest struct {
	ID                                     uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	PenetrationTestReportID                uint       `json:"penetration_test_report_id" gorm:"not null"`
	PenetrationTestReportPenetrationTestID uint       `json:"penetration_test_report_penetration_test_id" gorm:"not null"`
	PenetrationTestRetestID                uint       `json:"penetration_test_retest_id" gorm:"not null"`
	ReferenceNumber                        *string    `json:"reference_number" gorm:"type:nvarchar(50)"`
	StartDate                              time.Time  `json:"start_date" gorm:"type:date;not null"`
	EndDate                                *time.Time `json:"end_date" gorm:"type:date"`
	ExecutorTitleEn                        *string    `json:"executor_title_en" gorm:"type:nvarchar(255)"`
	ExecutorTitleAr                        *string    `json:"executor_title_ar" gorm:"type:nvarchar(255)"`
	StatusTitleEn                          *string    `json:"status_title_en" gorm:"type:nvarchar(255)"`
	StatusTitleAr                          *string    `json:"status_title_ar" gorm:"type:nvarchar(255)"`
	StatusColorCode                        *string    `json:"status_color_code" gorm:"type:nvarchar(255)"`
	Notes                                  *string    `json:"notes" gorm:"type:nvarchar(max)"`
	CreatedAt                              time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                              *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	PenetrationTestReport                *PenetrationTestReport                `gorm:"foreignKey:PenetrationTestReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	PenetrationTestReportPenetrationTest *PenetrationTestReportPenetrationTest `gorm:"foreignKey:PenetrationTestReportPenetrationTestID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	PenetrationTestRetest                *PenetrationTestRetest                `gorm:"foreignKey:PenetrationTestRetestID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	RetestFindings                       []*PenetrationTestReportRetestFinding `gorm:"foreignKey:PenetrationTestReportRetestID"`
}

func (PenetrationTestReportRetest) TableName() string {
	return "penetration_test_report_retests"
}
