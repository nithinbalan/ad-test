package models

import "time"

type DocumentReportDocument struct {
	ID                                   uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	DocumentReportID                     uint       `json:"document_report_id" gorm:"not null;index"`
	SourceDocumentID                     uint       `json:"source_document_id" gorm:"not null"`
	ReferenceNumber                      string     `json:"reference_number" gorm:"type:varchar(45);not null"`
	TitleEn                              string     `json:"title_en" gorm:"type:nvarchar(255);not null"`
	TitleAr                              string     `json:"title_ar" gorm:"type:nvarchar(255);not null"`
	DescriptionEn                        *string    `json:"description_en" gorm:"type:nvarchar(max)"`
	DescriptionAr                        *string    `json:"description_ar" gorm:"type:nvarchar(max)"`
	DocumentCategoryTitleEn              string     `json:"document_category_title_en" gorm:"type:nvarchar(255);not null"`
	DocumentCategoryTitleAr              string     `json:"document_category_title_ar" gorm:"type:nvarchar(255);not null"`
	DocumentTypeTitleEn                  string     `json:"document_type_title_en" gorm:"type:nvarchar(255);not null"`
	DocumentTypeTitleAr                  string     `json:"document_type_title_ar" gorm:"type:nvarchar(255);not null"`
	DocumentStatusTitleEn                string     `json:"document_status_title_en" gorm:"type:nvarchar(100);not null"`
	DocumentStatusTitleAr                string     `json:"document_status_title_ar" gorm:"type:nvarchar(100);not null"`
	DocumentStatusColorCode              string     `json:"document_status_color_code" gorm:"type:varchar(25);not null"`
	DocumentStatusType                   string     `json:"document_status_type" gorm:"type:varchar(10);not null"`
	DocumentAccessTypeTitleEn            string     `json:"document_access_type_title_en" gorm:"type:nvarchar(100);not null"`
	DocumentAccessTypeTitleAr            string     `json:"document_access_type_title_ar" gorm:"type:nvarchar(100);not null"`
	DocumentAccessTypeColorCode          string     `json:"document_access_type_color_code" gorm:"type:varchar(25);not null"`
	ConfidentialityClassificationTitleEn string     `json:"confidentiality_classification_title_en" gorm:"type:nvarchar(100);not null"`
	ConfidentialityClassificationTitleAr string     `json:"confidentiality_classification_title_ar" gorm:"type:nvarchar(100);not null"`
	ConfidentialityClassificationType    string     `json:"confidentiality_classification_type" gorm:"type:varchar(150);not null"`
	ReviewFrequencyTitleEn               *string    `json:"review_frequency_title_en" gorm:"type:nvarchar(100)"`
	ReviewFrequencyTitleAr               *string    `json:"review_frequency_title_ar" gorm:"type:nvarchar(100)"`
	ReviewFrequencyType                  *string    `json:"review_frequency_type" gorm:"type:varchar(20)"`
	ResponsibleUnitTitleEn               *string    `json:"responsible_unit_title_en" gorm:"type:nvarchar(255)"`
	ResponsibleUnitTitleAr               *string    `json:"responsible_unit_title_ar" gorm:"type:nvarchar(255)"`
	ResponsibleUserNameEn                *string    `json:"responsible_user_name_en" gorm:"type:nvarchar(511)"`
	ResponsibleUserNameAr                *string    `json:"responsible_user_name_ar" gorm:"type:nvarchar(511)"`
	Version                              *float64   `json:"version" gorm:"type:decimal(4,2)"`
	FileName                             *string    `json:"file_name" gorm:"type:nvarchar(255)"`
	FileType                             *string    `json:"file_type" gorm:"type:varchar(255)"`
	FileSize                             *int       `json:"file_size"`
	DocumentCreatedAt                    time.Time  `json:"document_created_at" gorm:"type:datetime;not null"`
	NextReviewDate                       *time.Time `json:"next_review_date" gorm:"type:date"`
	PublishedAt                          *time.Time `json:"published_at" gorm:"type:datetime"`
	ClausesJSON                          string     `json:"-" gorm:"column:clauses_json;type:nvarchar(max);not null"`
	ControlsJSON                         string     `json:"-" gorm:"column:controls_json;type:nvarchar(max);not null"`
	AssociatedItemsJSON                  string     `json:"-" gorm:"column:associated_items_json;type:nvarchar(max);not null"`
	CreatedByID                          uint       `json:"created_by_id" gorm:"not null"`
	CreatedAt                            time.Time  `json:"created_at" gorm:"type:datetime;not null"`

	DocumentReport *DocumentReport `gorm:"foreignKey:DocumentReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser  *User           `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (DocumentReportDocument) TableName() string {
	return "document_report_documents"
}
