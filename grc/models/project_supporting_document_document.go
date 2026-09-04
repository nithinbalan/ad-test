package models

import (
	"time"
)

type ProjectSupportingDocumentDocument struct {
	ID                          uint       `json:"id" gorm:"primary_key;auto_increment"`
	ProjectSupportingDocumentID uint       `json:"project_supporting_document_id" gorm:"not null"`
	URL                         string     `json:"url" gorm:"type:varchar(255)"`
	Name                        string     `json:"name" binding:"required" gorm:"type:varchar(100);not null"`
	Size                        int        `json:"size" gorm:"not null"`
	Token                       string     `json:"token" gorm:"type:varchar(max);not null"`
	OriginalName                *string    `json:"original_name" gorm:"type:nvarchar(255)"`
	Type                        string     `json:"type" gorm:"type:varchar(100);not null"`
	CreatedByID                 *uint      `json:"created_by_id" gorm:"null"`
	UpdatedByID                 *uint      `json:"updated_by_id" gorm:"null"`
	CreatedAt                   time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                   *time.Time `json:"updated_at" gorm:"type:datetime;null"`

	// Relationships
	ProjectSupportingDocument *ProjectSupportingDocument `gorm:"foreignKey:ProjectSupportingDocumentID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser             *User                      `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser             *User                      `gorm:"foreignKey:UpdatedByID"`
}

func (ProjectSupportingDocumentDocument) TableName() string {
	return "project_supporting_document_documents"
}
