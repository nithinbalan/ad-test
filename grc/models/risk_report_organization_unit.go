package models

type RiskReportOrganizationUnit struct {
	ID                    uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	RiskReportID          uint   `json:"risk_report_id" gorm:"not null"`
	OrganizationUnitID    uint   `json:"organization_unit_id" gorm:"not null"`
	OrganizationUnitTitle string `json:"organization_unit_title"`

	// Relationships
	RiskReport       *RiskReport       `gorm:"foreignKey:RiskReportID"`
	OrganizationUnit *OrganizationUnit `gorm:"foreignKey:OrganizationUnitID"`
}

func (RiskReportOrganizationUnit) TableName() string {
	return "risk_report_organization_unit"
}
