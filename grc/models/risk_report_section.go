package models

type RiskReportSection struct {
	ID            uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	RiskReportID  uint   `json:"risk_report_id" gorm:"not null"`
	Type          string `json:"type" gorm:"type:varchar(50);not null;check:type IN ('executive-summary', 'context', 'identification', 'assessment', 'mitigation', 'kri', 'communication', 'recommendation', 'documents')"`
	TitleEn       string `json:"title_en" gorm:"type:varchar(255);not null"`
	TitleAr       string `json:"title_ar" gorm:"type:varchar(255);not null"`
	DescriptionEn string `json:"description_en,omitempty" gorm:"type:text"`
	DescriptionAr string `json:"description_ar,omitempty" gorm:"type:text"`
	OrderNum      int    `json:"order_num" gorm:"not null"`

	RiskReport *RiskReport `gorm:"foreignKey:RiskReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (RiskReportSection) TableName() string {
	return "risk_report_sections"
}
