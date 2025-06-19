package models

import "time"

type GRCReportSection struct {
	ID            uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	GRCReportID   uint       `json:"grc_report_id" gorm:"not null"`
	Type          string     `json:"type" gorm:"type:varchar(100);not null"`
	TitleEn       string     `json:"title_en" gorm:"type:nvarchar(255);not null"`
	TitleAr       string     `json:"title_ar" gorm:"type:nvarchar(255);not null"`
	DescriptionEn string     `json:"description_en" gorm:"type:nvarchar(max)"`
	DescriptionAr string     `json:"description_ar" gorm:"type:nvarchar(max)"`
	OrderNum      int        `json:"order_num" gorm:"not null"`
	CreatedBYID   uint       `json:"created_by_id" gorm:"not null"`
	CreatedAt     time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedBYID   *uint      `json:"updated_by_id"`
	UpdatedAt     *time.Time `json:"updated_at"`

	CreatedByUser *User      `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser *User      `gorm:"foreignKey:UpdatedBYID"`
	GRCReport     *GRCReport `gorm:"foreignKey:GRCReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (GRCReportSection) TableName() string {
	return "grc_report_sections"
}
