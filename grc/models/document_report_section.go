package models

type DocumentReportSection struct {
	ID               uint    `json:"id" gorm:"primaryKey;autoIncrement"`
	DocumentReportID uint    `json:"document_report_id" gorm:"not null;index"`
	Type             string  `json:"type" gorm:"type:varchar(50);not null;check:type IN ('documents', 'clauses-controls', 'associated-items')"`
	TitleEn          string  `json:"title_en" gorm:"type:nvarchar(255);not null"`
	TitleAr          string  `json:"title_ar" gorm:"type:nvarchar(255);not null"`
	DescriptionEn    *string `json:"description_en" gorm:"type:nvarchar(max)"`
	DescriptionAr    *string `json:"description_ar" gorm:"type:nvarchar(max)"`
	OrderNum         int     `json:"order_num" gorm:"not null"`

	DocumentReport *DocumentReport `gorm:"foreignKey:DocumentReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (DocumentReportSection) TableName() string {
	return "document_report_sections"
}
