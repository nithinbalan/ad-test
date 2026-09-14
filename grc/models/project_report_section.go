package models

type ProjectReportSection struct {
	ID              uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ProjectReportID uint   `json:"project_report_id" gorm:"not null"`
	Type            string `json:"type" gorm:"type:varchar(50);not null"`
	TitleEn         string `json:"title_en" gorm:"type:varchar(255);not null"`
	TitleAr         string `json:"title_ar" gorm:"type:varchar(255);not null"`
	DescriptionEn   string `json:"description_en" gorm:"type:nvarchar(max)"`
	DescriptionAr   string `json:"description_ar" gorm:"type:nvarchar(max)"`
	OrderNum        int    `json:"order_num" gorm:"not null"`

	ProjectReport *ProjectReport `gorm:"foreignKey:ProjectReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (ProjectReportSection) TableName() string {
	return "project_report_sections"
}
