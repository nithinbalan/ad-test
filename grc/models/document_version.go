package models

import (
	"time"
)

type DocumentVersion struct {
	ID              uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	DocumentID      uint       `json:"document_id" binding:"required" gorm:"not null;index;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Version         float64    `json:"version" binding:"required" gorm:"type:decimal(4,2);not null"`
	IsLatestVersion bool       `json:"is_latest_version" binding:"required" gorm:"type:tinyint(1);not null;default:1"`
	URL             string     `json:"url" binding:"required" gorm:"type:varchar(1000);not null"`
	Name            string     `json:"name" binding:"required" gorm:"type:nvarchar(255);not null"`
	Size            int        `json:"size" binding:"required" gorm:"not null"`
	Token           string     `json:"token" binding:"required" gorm:"type:varchar(255);not null"`
	OriginalName    *string    `json:"original_name" binding:"omitempty" gorm:"type:nvarchar(255)"`
	Type            string     `json:"type" binding:"required" gorm:"type:varchar(255);not null"`
	SubmittedBYID   *uint      `json:"submitted_by_id" gorm:"default:null"`
	SubmittedAt     *time.Time `json:"submitted_at" gorm:"type:date;default:null"`
	PublishedBYID   *uint      `json:"published_by_id" gorm:"default:null"`
	PublishedAt     *time.Time `json:"published_at" gorm:"type:date;default:null"`
	NextReviewLevel *int       `json:"next_review_level" gorm:"default:null"`
	CreatedBYID     uint       `json:"created_by_id" binding:"required" gorm:"not null"`
	UpdatedBYID     *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt       time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	UpdatedAt       *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	//Relations
	Document        *Document `gorm:"foreignKey:DocumentID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	SubmittedByUser *User     `gorm:"foreignKey:SubmittedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser   *User     `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	UpdatedByUser   *User     `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	DocumentEditHistory []*DocumentEditHistory `gorm:"foreignKey:DocumentVersionID"`
}

func (DocumentVersion) TableName() string {
	return "document_versions"
}
