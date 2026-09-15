package models

import "time"

type PenetrationTestReportFinding struct {
	ID                                     uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	PenetrationTestReportID                uint       `json:"penetration_test_report_id" gorm:"not null"`
	PenetrationTestReportPenetrationTestID uint       `json:"penetration_test_report_penetration_test_id" gorm:"not null"`
	PenetrationTestFindingID               uint       `json:"penetration_test_finding_id" gorm:"not null"`
	ReferenceNumber                        *string    `json:"reference_number" gorm:"type:nvarchar(50)"`
	Title                                  string     `json:"title" gorm:"type:nvarchar(255);not null"`
	Description                            *string    `json:"description" gorm:"type:nvarchar(max)"`
	Recommendation                         *string    `json:"recommendation" gorm:"type:nvarchar(max)"`
	SeverityLevelTitleEn                   *string    `json:"severity_level_title_en" gorm:"type:nvarchar(255)"`
	SeverityLevelTitleAr                   *string    `json:"severity_level_title_ar" gorm:"type:nvarchar(255)"`
	SeverityLevelColorCode                 *string    `json:"severity_level_color_code" gorm:"type:nvarchar(255)"`
	OwnerNameEn                            *string    `json:"owner_name_en" gorm:"type:nvarchar(255)"`
	OwnerNameAr                            *string    `json:"owner_name_ar" gorm:"type:nvarchar(255)"`
	ActionPlanStatusTitleEn                *string    `json:"action_plan_status_title_en" gorm:"type:nvarchar(255)"`
	ActionPlanStatusTitleAr                *string    `json:"action_plan_status_title_ar" gorm:"type:nvarchar(255)"`
	ActionPlanStatusColorCode              *string    `json:"action_plan_status_color_code" gorm:"type:nvarchar(255)"`
	Percentage                             *int       `json:"percentage"`
	TargetRemediationDate                  *time.Time `json:"target_remediation_date" gorm:"type:date"`
	CreatedAt                              *time.Time `json:"created_at" gorm:"type:datetime"`
	UpdatedAt                              *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	PenetrationTestReport                *PenetrationTestReport                `gorm:"foreignKey:PenetrationTestReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	PenetrationTestReportPenetrationTest *PenetrationTestReportPenetrationTest `gorm:"foreignKey:PenetrationTestReportPenetrationTestID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	PenetrationTestFinding               *PenetrationTestFinding               `gorm:"foreignKey:PenetrationTestFindingID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (PenetrationTestReportFinding) TableName() string {
	return "penetration_test_report_findings"
}
