package models

import (
	"time"
)

type PenetrationTestSupportingDocument struct {
	ID                uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	PenetrationTestID uint       `json:"penetration_test_id" gorm:"not null"`
	TitleEn           *string    `json:"title_en" gorm:"type:nvarchar(255);null"`
	TitleAr           *string    `json:"title_ar" gorm:"type:nvarchar(255);null"`
	CreatedByID       *uint      `json:"created_by_id" gorm:"null"`
	UpdatedByID       *uint      `json:"updated_by_id" gorm:"null"`
	CreatedAt         time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt         *time.Time `json:"updated_at" gorm:"type:datetime;null"`

	PenetrationTest *PenetrationTest                             `gorm:"foreignKey:PenetrationTestID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser   *User                                        `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:NO ACTION,OnDelete:SET NULL"`
	UpdatedByUser   *User                                        `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:NO ACTION,OnDelete:SET NULL"`
	Documents       []*PenetrationTestSupportingDocumentDocument `json:"documents" gorm:"foreignKey:PenetrationTestSupportingDocumentID"`
}

func (PenetrationTestSupportingDocument) TableName() string {
	return "penetration_test_supporting_documents"
}
