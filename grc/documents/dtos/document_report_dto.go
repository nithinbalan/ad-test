package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type DocumentReportListQueryParamsDTO struct {
	Sort     string
	SortBy   string
	Keyword  string
	Page     int
	PageSize int
}

type DocumentReportListDTO struct {
	ID                      uint      `json:"id"`
	Title                   string    `json:"title"`
	DocumentCategoryID      *uint     `json:"document_category_id"`
	DocumentCategoryTitleEn *string   `json:"document_category_title_en"`
	DocumentCategoryTitleAr *string   `json:"document_category_title_ar"`
	DocumentTypeID          *uint     `json:"document_type_id"`
	DocumentTypeTitleEn     *string   `json:"document_type_title_en"`
	DocumentTypeTitleAr     *string   `json:"document_type_title_ar"`
	CreatedAt               time.Time `json:"created_at"`
	CreatedUserID           uint      `json:"created_user_id"`
	CreatedUserFirstNameEn  string    `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr  string    `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn   string    `json:"created_user_last_name_en"`
	CreatedUserLastNameAr   string    `json:"created_user_last_name_ar"`
	CreatedUserInitial      string    `json:"created_user_initial"`
	CreatedUserColor        string    `json:"created_user_color"`
	CreatedUserEmail        string    `json:"created_user_email"`
	CreatedUserProfileImage *string   `json:"created_user_profile_image_path"`
	CreatedUserMobileNumber *string   `json:"created_user_mobile_number"`
}

type DocumentReportListResponseDTO struct {
	Items      []*DocumentReportListDTO     `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type DocumentReportFormDTO struct {
	Title              string `json:"title" binding:"required"`
	DocumentCategoryID *uint  `json:"document_category_id"`
	DocumentTypeID     *uint  `json:"document_type_id"`
}

type DocumentReportSectionResponseDTO struct {
	ID            uint    `json:"id"`
	Type          string  `json:"type"`
	TitleEn       string  `json:"title_en"`
	TitleAr       string  `json:"title_ar"`
	DescriptionEn *string `json:"description_en"`
	DescriptionAr *string `json:"description_ar"`
	OrderNumber   int     `json:"order_number"`
}

type DocumentReportClauseDTO struct {
	ID            uint                       `json:"id"`
	ClauseNumber  *string                    `json:"clause_number"`
	ClauseTitleEn string                     `json:"clause_title_en"`
	ClauseTitleAr string                     `json:"clause_title_ar"`
	Description   *string                    `json:"description"`
	Controls      []DocumentReportControlDTO `json:"controls" gorm:"-"`
}

type DocumentReportControlDTO struct {
	ID              uint    `json:"id"`
	ClauseID        *uint   `json:"-"`
	ReferenceNumber string  `json:"reference_number"`
	Code            *string `json:"code"`
	TitleEn         string  `json:"title_en"`
	TitleAr         string  `json:"title_ar"`
	Description     *string `json:"description"`
}

type DocumentReportDocumentDTO struct {
	ID                                   uint                      `json:"id"`
	ReferenceNumber                      string                    `json:"reference_number"`
	TitleEn                              string                    `json:"title_en"`
	TitleAr                              string                    `json:"title_ar"`
	DescriptionEn                        *string                   `json:"description_en"`
	DescriptionAr                        *string                   `json:"description_ar"`
	DocumentCategoryTitleEn              string                    `json:"document_category_title_en"`
	DocumentCategoryTitleAr              string                    `json:"document_category_title_ar"`
	DocumentTypeTitleEn                  string                    `json:"document_type_title_en"`
	DocumentTypeTitleAr                  string                    `json:"document_type_title_ar"`
	DocumentStatusTitleEn                string                    `json:"document_status_title_en"`
	DocumentStatusTitleAr                string                    `json:"document_status_title_ar"`
	DocumentStatusColorCode              string                    `json:"document_status_color_code"`
	DocumentStatusType                   string                    `json:"document_status_type"`
	DocumentAccessTypeTitleEn            string                    `json:"document_access_type_title_en"`
	DocumentAccessTypeTitleAr            string                    `json:"document_access_type_title_ar"`
	DocumentAccessTypeColorCode          string                    `json:"document_access_type_color_code"`
	ConfidentialityClassificationTitleEn string                    `json:"confidentiality_classification_title_en"`
	ConfidentialityClassificationTitleAr string                    `json:"confidentiality_classification_title_ar"`
	ConfidentialityClassificationType    string                    `json:"confidentiality_classification_type"`
	ReviewFrequencyTitleEn               *string                   `json:"review_frequency_title_en"`
	ReviewFrequencyTitleAr               *string                   `json:"review_frequency_title_ar"`
	ReviewFrequencyType                  *string                   `json:"review_frequency_type"`
	ResponsibleUnitTitleEn               *string                   `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr               *string                   `json:"responsible_unit_title_ar"`
	ResponsibleUserNameEn                *string                   `json:"responsible_user_name_en"`
	ResponsibleUserNameAr                *string                   `json:"responsible_user_name_ar"`
	Version                              *float64                  `json:"version"`
	FileName                             *string                   `json:"file_name"`
	FileType                             *string                   `json:"file_type"`
	FileSize                             *int                      `json:"file_size"`
	DocumentCreatedAt                    time.Time                 `json:"document_created_at"`
	NextReviewDate                       *time.Time                `json:"next_review_date"`
	PublishedAt                          *string                   `json:"published_at" mapstructure:"-"`
	Clauses                              []DocumentReportClauseDTO `json:"clauses"`
}

type DocumentReportDetailDTO struct {
	ID                 uint                               `json:"id"`
	Title              string                             `json:"title"`
	DocumentCategoryID *uint                              `json:"document_category_id"`
	DocumentTypeID     *uint                              `json:"document_type_id"`
	CreatedAt          time.Time                          `json:"created_at"`
	CreatedBy          organizationDTOs.UserMinimalDTO    `json:"created_by"`
	Sections           []DocumentReportSectionResponseDTO `json:"sections"`
	Documents          []*DocumentReportDocumentDTO       `json:"documents"`
}

type DocumentReportPDFDetailDTO struct {
	ID        uint                            `json:"id"`
	Title     string                          `json:"title"`
	CreatedAt string                          `json:"created_at"`
	CreatedBy organizationDTOs.UserMinimalDTO `json:"created_by"`
	Documents []*DocumentReportDocumentDTO    `json:"documents"`
}

// DocumentReportSourceDocumentDTO is used only while creating an immutable report snapshot.
type DocumentReportSourceDocumentDTO struct {
	SourceDocumentID                                                           uint
	ReferenceNumber, TitleEn, TitleAr                                          string
	DescriptionEn, DescriptionAr                                               *string
	DocumentCategoryTitleEn, DocumentCategoryTitleAr                           string
	DocumentTypeTitleEn, DocumentTypeTitleAr                                   string
	DocumentStatusTitleEn, DocumentStatusTitleAr                               string
	DocumentStatusColorCode, DocumentStatusType                                string
	DocumentAccessTypeTitleEn, DocumentAccessTypeTitleAr                       string
	DocumentAccessTypeColorCode                                                string
	ConfidentialityClassificationTitleEn, ConfidentialityClassificationTitleAr string
	ConfidentialityClassificationType                                          string
	ReviewFrequencyTitleEn, ReviewFrequencyTitleAr                             *string
	ReviewFrequencyType                                                        *string
	ResponsibleUnitTitleEn, ResponsibleUnitTitleAr                             *string
	ResponsibleUserNameEn, ResponsibleUserNameAr                               *string
	Version                                                                    *float64
	FileName                                                                   *string
	FileType                                                                   *string
	FileSize                                                                   *int
	DocumentCreatedAt                                                          time.Time
	NextReviewDate                                                             *time.Time
	PublishedAt                                                                *time.Time
}
