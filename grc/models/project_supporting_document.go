package models

import (
	"time"
)

type ProjectSupportingDocument struct {
	ID          uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ProjectID   uint       `json:"project_id" gorm:"not null"`
	TitleEn     *string    `json:"title_en" gorm:"type:nvarchar(255);null"`
	TitleAr     *string    `json:"title_ar" gorm:"type:nvarchar(255);null"`
	CreatedByID *uint      `json:"created_by_id" gorm:"null"`
	UpdatedByID *uint      `json:"updated_by_id" gorm:"null"`
	CreatedAt   time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt   *time.Time `json:"updated_at" gorm:"type:datetime;null"`

	// Relationships
	Project       *Project                             `gorm:"foreignKey:ProjectID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser *User                                `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:NO ACTION,OnDelete:SET NULL"`
	UpdatedByUser *User                                `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:NO ACTION,OnDelete:SET NULL"`
	Documents     []*ProjectSupportingDocumentDocument `json:"documents" gorm:"foreignKey:ProjectSupportingDocumentID"`
}

func (ProjectSupportingDocument) TableName() string {
	return "project_supporting_documents"
}
