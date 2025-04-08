package models

type GovernanceReportSection struct {
	ID                 uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	GovernanceReportID uint   `json:"governance_report_id" gorm:"not null"`
	Type               string `json:"type" gorm:"type:varchar(50);not null;check:type IN ('executive-summary', 'policy-acknowledgement', 'meetings', 'actions', 'conclusion')"`
	TitleEn            string `json:"title_en" gorm:"type:varchar(255);not null"`
	TitleAr            string `json:"title_ar" gorm:"type:varchar(255);not null"`
	DescriptionEn      string `json:"description_en" gorm:"type:text"`
	DescriptionAr      string `json:"description_ar" gorm:"type:text"`
	OrderNum           int    `json:"order_num" gorm:"not null"`

	GovernanceReport *GovernanceReport `gorm:"foreignKey:GovernanceReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (GovernanceReportSection) TableName() string {
	return "governance_report_sections"
}
