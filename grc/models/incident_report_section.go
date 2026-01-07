package models

type IncidentReportSection struct {
	ID               uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	IncidentReportID uint   `json:"incident_report_id" gorm:"not null"`
	Type             string `json:"type" gorm:"type:varchar(50);not null;check:type IN ('executive-summary', 'scope', 'objectives', 'overview', 'details', 'containment', 'recovery', 'post-incident')"`
	TitleEn          string `json:"title_en" gorm:"type:varchar(255);not null"`
	TitleAr          string `json:"title_ar" gorm:"type:varchar(255);not null"`
	DescriptionEn    string `json:"description_en" gorm:"type:text"`
	DescriptionAr    string `json:"description_ar" gorm:"type:text"`
	OrderNum         int    `json:"order_num" gorm:"not null"`

	IncidentReport *IncidentReport `gorm:"foreignKey:IncidentReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (IncidentReportSection) TableName() string {
	return "incident_report_sections"
}