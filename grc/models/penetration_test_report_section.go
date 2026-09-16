package models

type PenetrationTestReportSection struct {
	ID                      uint    `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	PenetrationTestReportID uint    `json:"penetration_test_report_id" gorm:"not null"`
	Type                    string  `json:"type" gorm:"type:varchar(50);not null;check:type IN ('executive-summary', 'scope', 'objectives', 'overview', 'details', 'findings', 'retest-results', 'management-sharing', 'closure')"`
	TitleEn                 string  `json:"title_en" gorm:"type:nvarchar(255);not null"`
	TitleAr                 string  `json:"title_ar" gorm:"type:nvarchar(255);not null"`
	DescriptionEn           *string `json:"description_en" gorm:"type:nvarchar(max)"`
	DescriptionAr           *string `json:"description_ar" gorm:"type:nvarchar(max)"`
	OrderNum                int     `json:"order_num" gorm:"not null"`

	// Relationships
	PenetrationTestReport *PenetrationTestReport `gorm:"foreignKey:PenetrationTestReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (PenetrationTestReportSection) TableName() string {
	return "penetration_test_report_sections"
}
