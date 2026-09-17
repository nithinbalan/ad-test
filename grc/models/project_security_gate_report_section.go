package models

type ProjectSecurityGateReportSection struct {
	ID            uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ReportID      uint   `json:"report_id" gorm:"not null"`
	Type          string `json:"type" gorm:"type:varchar(100);not null"`
	TitleEn       string `json:"title_en" gorm:"type:nvarchar(255);not null"`
	TitleAr       string `json:"title_ar" gorm:"type:nvarchar(255);not null"`
	DescriptionEn string `json:"description_en" gorm:"type:nvarchar(max)"`
	DescriptionAr string `json:"description_ar" gorm:"type:nvarchar(max)"`
	OrderNum      int    `json:"order_num" gorm:"not null"`

	Report *ProjectSecurityGateReport `gorm:"foreignKey:ReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (ProjectSecurityGateReportSection) TableName() string {
	return "project_security_gate_report_sections"
}
