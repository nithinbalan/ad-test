package models

import "time"

type AssetOtherDocument struct {
	ID                      uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	AssetRegisterID         uint       `json:"asset_register_id" gorm:"not null"`
	AssetDocumentCategoryID uint       `json:"asset_document_category_id" gorm:"not null"`
	URL                     *string    `json:"url,omitempty" gorm:"type:varchar(1000);default:null"`
	Name                    *string    `json:"name,omitempty" gorm:"type:nvarchar(255);default:null"`
	Size                    *int       `json:"size,omitempty" gorm:"default:null"`
	Token                   *string    `json:"token,omitempty" gorm:"type:varchar(255);default:null"`
	OriginalName            *string    `json:"original_name,omitempty" gorm:"type:nvarchar(255);default:null"`
	Type                    *string    `json:"type,omitempty" gorm:"type:varchar(255);default:null"`
	CreatedBYID             uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID             *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt               time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt               *time.Time `json:"updated_at,omitempty" gorm:"type:datetime;default:null"`

	// Relationships
	AssetRegister         *AssetRegister         `gorm:"foreignKey:AssetRegisterID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	AssetDocumentCategory *AssetDocumentCategory `gorm:"foreignKey:AssetDocumentCategoryID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser         *User                  `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedByUser         *User                  `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (AssetOtherDocument) TableName() string {
	return "asset_other_documents"
}
