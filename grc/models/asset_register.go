package models

import (
	"time"
)

type AssetRegister struct {
	ID                uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	ReferenceNumber   string     `json:"reference_number" gorm:"type:varchar(45);not null;unique" binding:"required"`
	TitleEn           string     `json:"title_en" gorm:"type:nvarchar(500);not null" binding:"required"`
	TitleAr           string     `json:"title_ar" gorm:"type:nvarchar(500);not null" binding:"required"`
	AssetTypeID       uint       `json:"asset_type_id" gorm:"not null" binding:"required"`
	AssetCategoryID   uint       `json:"asset_category_id" gorm:"not null" binding:"required"`
	AssetStatusID     uint       `json:"asset_status_id" gorm:"not null" binding:"required"`
	ResponsibleUserID uint       `json:"responsible_user_id" gorm:"not null" binding:"required"`
	Location          *string    `json:"location,omitempty" gorm:"type:nvarchar(255);default:null"`
	Description       *string    `json:"description,omitempty" gorm:"type:nvarchar(MAX);default:null"`
	AcquisitionDate   *time.Time `json:"acquisition_date,omitempty" gorm:"type:date;default:null"`
	CreatedByID       uint       `json:"created_by_id" gorm:"not null" binding:"required"`
	CreatedAt         time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedByID       *uint      `json:"updated_by_id" gorm:"default:null"`
	UpdatedAt         *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	AssetType       *AssetType     `gorm:"foreignKey:AssetTypeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	AssetCategory   *AssetCategory `gorm:"foreignKey:AssetCategoryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	AssetStatus     *AssetStatus   `gorm:"foreignKey:AssetStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResponsibleUser *User          `gorm:"foreignKey:ResponsibleUserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser   *User          `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	UpdatedByUser   *User          `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// TableName overrides the default table name.
func (AssetRegister) TableName() string {
	return "asset_register"
}
