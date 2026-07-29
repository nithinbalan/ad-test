package models

import "time"

type DocumentReport struct {
	ID                 uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Title              string    `json:"title" gorm:"type:nvarchar(255);not null"`
	DocumentCategoryID uint      `json:"document_category_id" gorm:"not null"`
	DocumentTypeID     uint      `json:"document_type_id" gorm:"not null"`
	CreatedByID        uint      `json:"created_by_id" gorm:"not null"`
	CreatedAt          time.Time `json:"created_at" gorm:"type:datetime;not null;index"`

	DocumentCategory *DocumentCategory         `gorm:"foreignKey:DocumentCategoryID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	DocumentType     *DocumentType             `gorm:"foreignKey:DocumentTypeID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser    *User                     `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Sections         []*DocumentReportSection  `gorm:"foreignKey:DocumentReportID"`
	Documents        []*DocumentReportDocument `gorm:"foreignKey:DocumentReportID"`
}

func (DocumentReport) TableName() string {
	return "document_reports"
}
