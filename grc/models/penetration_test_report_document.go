package models

import "time"

type PenetrationTestReportDocument struct {
	ID                                          uint       `json:"id" gorm:"primary_key;auto_increment"`
	PenetrationTestReportPenetrationTestID      uint       `json:"penetration_test_report_penetration_test_id" gorm:"not null"`
	PenetrationTestSupportingDocumentDocumentID *uint      `json:"penetration_test_supporting_document_document_id"`
	PenetrationTestRetestFindingDocumentID      *uint      `json:"penetration_test_retest_finding_document_id"`
	PenetrationTestReportRetestFindingID        *uint      `json:"penetration_test_report_retest_finding_id"`
	URL                                         *string    `json:"url" gorm:"type:varchar(1000)"`
	Name                                        *string    `json:"name" gorm:"type:nvarchar(255)"`
	Size                                        *int       `json:"size" gorm:"type:int"`
	Token                                       *string    `json:"token" gorm:"type:varchar(255)"`
	OriginalName                                *string    `json:"original_name" gorm:"type:nvarchar(255)"`
	Type                                        *string    `json:"type" gorm:"type:varchar(255)"`
	CreatedBYID                                 *uint      `json:"created_by_id" gorm:"type:int"`
	UpdatedBYID                                 *uint      `json:"updated_by_id" gorm:"type:int"`
	CreatedAt                                   time.Time  `json:"created_at" gorm:"type:datetime;not null;default:getdate()"`
	UpdatedAt                                   *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	PenetrationTestReportPenetrationTest      *PenetrationTestReportPenetrationTest      `gorm:"foreignKey:PenetrationTestReportPenetrationTestID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	PenetrationTestSupportingDocumentDocument *PenetrationTestSupportingDocumentDocument `gorm:"foreignKey:PenetrationTestSupportingDocumentDocumentID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	PenetrationTestRetestFindingDocument      *PenetrationTestRetestFindingDocument      `gorm:"foreignKey:PenetrationTestRetestFindingDocumentID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	PenetrationTestReportRetestFinding        *PenetrationTestReportRetestFinding        `gorm:"foreignKey:PenetrationTestReportRetestFindingID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser                             *User                                      `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedByUser                             *User                                      `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (PenetrationTestReportDocument) TableName() string {
	return "penetration_test_report_documents"
}
