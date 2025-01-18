package models

import (
	"time"
)

type Document struct {
	ID                    uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ReferenceNumber       string     `json:"reference_number" binding:"required" gorm:"type:varchar(45);not null"`
	DocumentTypeID        uint       `json:"document_type_id" binding:"required" gorm:"not null"`
	TitleEn               string     `json:"title_en" binding:"required" gorm:"type:nvarchar(255);not null"`
	TitleAr               string     `json:"title_ar" binding:"required" gorm:"type:nvarchar(255);not null"`
	DescriptionEn         *string    `json:"description_en" gorm:"type:nvarchar(max);default:null"`
	DescriptionAr         *string    `json:"description_ar" gorm:"type:nvarchar(max);default:null"`
	DocumentCategoryID    uint       `json:"document_category_id" binding:"required" gorm:"not null"`
	ResponsibleUserID     uint       `json:"responsible_user_id" binding:"required" gorm:"not null"`
	ResponsibleUnitID     uint       `json:"responsible_unit_id" binding:"required" gorm:"not null"`
	DocumentAccessTypeID  uint       `json:"document_access_type_id" binding:"required" gorm:"not null"`
	DocumentStatusID      uint       `json:"document_status_id" binding:"required" gorm:"not null"`
	MarkedAsAuditableItem bool       `json:"marked_as_auditable_item" gorm:"default:false" binding:"-"`
	IsLocked              bool       `json:"is_locked" binding:"required" gorm:"type:tinyint(1);not null;default:0"`
	LastCheckoutAt        *time.Time `json:"last_checkout_at" gorm:"type:date;default:null"`
	LastCheckoutBYID      *uint      `json:"last_checkout_by_id" gorm:"default:null"`
	LastCheckinAt         *time.Time `json:"last_checkin_at" gorm:"type:date;default:null"`
	LastCheckinBYID       *uint      `json:"last_checkin_by_id" gorm:"default:null"`
	Rto                   *int       `json:"rto" gorm:"default:null"`
	RtoUnitId             *uint      `json:"rto_uint_id" gorm:"default:null"`
	Mtpd                  *int       `json:"mtpd" gorm:"default:null"`
	MtpdUnitId            *uint      `json:"mtpd_uint_id" gorm:"default:null"`
	Rpo                   *int       `json:"rpo" gorm:"default:null"`
	RpoUnitId             *uint      `json:"rpo_uint_id" gorm:"default:null"`
	CreatedBYID           uint       `json:"created_by_id" binding:"required" gorm:"not null"`
	UpdatedBYID           *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt             time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	UpdatedAt             *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	PublishedBYID         *uint      `json:"published_by_ID" gorm:"default:null"`
	PublishedAt           *time.Time `json:"published_at" gorm:"type:datetime;default:null"`
	ReviewFrequencyID     *uint      `json:"review_frequency_id,omitempty" gorm:"default:null" binding:"omitempty"`
	NextReviewDate        *time.Time `json:"next_review_date,omitempty" gorm:"type:date;default:null" binding:"omitempty"`

	// Relations
	DocumentType       *DocumentType       `gorm:"foreignKey:DocumentTypeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	DocumentCategory   *DocumentCategory   `gorm:"foreignKey:DocumentCategoryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResponsibleUser    *User               `gorm:"foreignKey:ResponsibleUserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResponsibleUnit    *OrganizationUnit   `gorm:"foreignKey:ResponsibleUnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	DocumentAccessType *DocumentAccessType `gorm:"foreignKey:DocumentAccessTypeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ReviewFrequency    *ReviewFrequency    `gorm:"foreignKey:ReviewFrequencyID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	DocumentStatus     *DocumentStatus     `gorm:"foreignKey:DocumentStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser      *User               `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	UpdatedByUser      *User               `gorm:"foreignKey:UpdatedBYID"`
	LastCheckoutByUser *User               `gorm:"foreignKey:LastCheckoutBYID"`
	LastCheckinByUser  *User               `gorm:"foreignKey:LastCheckinBYID"`
	PublishedByUser    *User               `gorm:"foreignKey:PublishedBYID"`

	LatestDocumentVersion *DocumentVersion      `gorm:"foreignKey:DocumentID"`
	DocumentVersions      []*DocumentVersion    `gorm:"foreignKey:DocumentID"`
	DocumentSharedUsers   []*DocumentSharedUser `gorm:"foreignKey:DocumentID"`
	DocumentSharedUnits   []*DocumentSharedUnit `gorm:"foreignKey:DocumentID"`
	DocumentAuditLogs     []*DocumentAuditLog   `gorm:"foreignKey:DocumentID"`

	// New Field
	Acknowledged bool `json:"acknowledged,omitempty" gorm:"-"`
}

func (Document) TableName() string {
	return "documents"
}
