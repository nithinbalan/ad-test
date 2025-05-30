package models

import "time"

type UserDocument struct {
	ID                    uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	UserID                uint       `json:"user_id" gorm:"not null"`
	UserDocumentTypeID    uint       `json:"user_document_type_id"`
	UserDocumentSubTypeID uint       `json:"user_document_sub_type_id"`
	IssuedBy              *string    `json:"issued_by"`
	ExpiryDate            *time.Time `json:"expiry_date"`
	Remarks               *string    `json:"remarks"`
	URL                   string     `json:"url"`
	Name                  string     `json:"name" binding:"required"`
	OriginalName          string     `json:"original_name"`
	Size                  int        `json:"size"`
	Token                 string     `json:"token"`
	Type                  string     `json:"type"`
	CreatedBYID           uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID           *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt             time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt             *time.Time `json:"updated_at" gorm:"type:datetime;not null"`

	// Relationships
	User                *User                `gorm:"foreignKey:UserID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	UserDocumentType    *UserDocumentType    `gorm:"foreignKey:UserDocumentTypeID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	UserDocumentSubType *UserDocumentSubType `gorm:"foreignKey:UserDocumentSubTypeID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	CreatedByUser       *User                `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser       *User                `gorm:"foreignKey:UpdatedBYID"`
}

// TableName overrides the table name used by UserDocument to `user_documents`
func (UserDocument) TableName() string {
	return "user_documents"
}
