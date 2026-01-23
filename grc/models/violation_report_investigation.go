package models

import "time"

type ViolationReportInvestigation struct {
	ID                                 uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	ViolationReportID                  uint       `json:"violation_report_id" gorm:"not null"`
	StepsTaken                         *string    `json:"steps_taken" gorm:"type:nvarchar(max)"`
	InvestigationDate                  time.Time  `json:"investigation_date" gorm:"type:datetime;not null"`
	SuggestedRecommendation            *string    `json:"suggested_recommendation" gorm:"type:nvarchar(max)"`
	FindingSummary                     *string    `json:"finding_summary" gorm:"type:nvarchar(max)"`
	DisciplinaryProcedure              *string    `json:"disciplinary_procedure" gorm:"type:nvarchar(max)"`
	IsDisciplinaryRegulatoryCompliance *bool      `json:"is_disciplinary_regulatory_compliance"`
	ActualRecommendation               *string    `json:"actual_recommendation" gorm:"type:nvarchar(max)"`
	NoteAndObservations                *string    `json:"note_and_observations" gorm:"type:nvarchar(max)"`
	CreatedBYID                        uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID                        *uint      `json:"updated_by_id"`
	CreatedAt                          time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                          *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	ViolationReport ViolationReport                           `gorm:"foreignKey:ViolationReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser   User                                      `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedByUser   *User                                     `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Departments     []*ViolationReportInvestigationDepartment `gorm:"foreignKey:ViolationReportInvestigationID"`
	Documents       []*ViolationReportInvestigationDocument   `gorm:"foreignKey:ViolationReportInvestigationID"`
}

// TableName overrides the table name
func (ViolationReportInvestigation) TableName() string {
	return "violation_report_investigations"
}
