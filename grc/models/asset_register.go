package models

import (
	"time"
)

type AssetRegister struct {
	ID                                  uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	ReferenceNumber                     string     `json:"reference_number" gorm:"type:varchar(45);not null;unique" binding:"required"`
	TitleEn                             string     `json:"title_en" gorm:"type:nvarchar(500);not null" binding:"required"`
	TitleAr                             string     `json:"title_ar" gorm:"type:nvarchar(500);not null" binding:"required"`
	AssetCategoryID                     uint       `json:"asset_category_id" gorm:"not null" binding:"required"`
	AssetSubCategoryID                  uint       `json:"asset_sub_category_id" gorm:"not null" binding:"required"`
	AssetStatusID                       uint       `json:"asset_status_id" gorm:"not null" binding:"required"`
	ResponsibleUserID                   uint       `json:"responsible_user_id" gorm:"not null" binding:"required"`
	ResponsibleUnitID                   uint       `json:"responsible_unit_id" gorm:"not null" binding:"required"`
	Location                            *string    `json:"location,omitempty" gorm:"type:nvarchar(255);default:null"`
	Description                         *string    `json:"description,omitempty" gorm:"type:nvarchar(MAX);default:null"`
	AcquisitionDate                     *time.Time `json:"acquisition_date,omitempty" gorm:"type:date;default:null"`
	AssignmentDate                      *time.Time `json:"assignment_date,omitempty" gorm:"type:date;default:null"`
	WarrantyExpiryDate                  *time.Time `json:"warranty_expiry_date,omitempty" gorm:"type:date;default:null"`
	Manufacturer                        *string    `json:"manufacturer,omitempty" gorm:"type:nvarchar(255);default:null"`
	Model                               *string    `json:"model,omitempty" gorm:"type:nvarchar(255);default:null"`
	SerialNumber                        *string    `json:"serial_number,omitempty" gorm:"type:nvarchar(255);default:null"`
	AssetCriticalityMatrixID            uint       `json:"asset_criticality_matrix_id"`
	AssetCriticalityMatrixRatingScaleID *uint      `json:"asset_criticality_matrix_rating_scale_id,omitempty" gorm:"default:null"`
	TotalScore                          *float64   `json:"total_score,omitempty" gorm:"default:null"`
	Rto                                 *int       `json:"rto" gorm:"default:null"`
	RtoUnitId                           *uint      `json:"rto_uint_id" gorm:"default:null"`
	Mtpd                                *int       `json:"mtpd" gorm:"default:null"`
	MtpdUnitId                          *uint      `json:"mtpd_uint_id" gorm:"default:null"`
	Rpo                                 *int       `json:"rpo" gorm:"default:null"`
	RpoUnitId                           *uint      `json:"rpo_uint_id" gorm:"default:null"`
	CreatedBYID                         uint       `json:"created_by_id" gorm:"not null" binding:"required"`
	CreatedAt                           time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedBYID                         *uint      `json:"updated_by_id,omitempty" gorm:"default:null"`
	UpdatedAt                           *time.Time `json:"updated_at,omitempty" gorm:"type:datetime;default:null"`

	// Relationships
	AssetCategory                     *AssetCategory                     `gorm:"foreignKey:AssetCategoryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	AssetSubCategory                  *AssetSubCategory                  `gorm:"foreignKey:AssetSubCategoryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	AssetStatus                       *AssetStatus                       `gorm:"foreignKey:AssetStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResponsibleUser                   *User                              `gorm:"foreignKey:ResponsibleUserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResponsibleUnit                   *OrganizationUnit                  `gorm:"foreignKey:ResponsibleUnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	AssetCriticalityMatrix            *AssetCriticalityMatrix            `gorm:"foreignKey:AssetCriticalityMatrixID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"` // New relationship
	AssetCriticalityMatrixRatingScale *AssetCriticalityMatrixRatingScale `gorm:"foreignKey:AssetCriticalityMatrixRatingScaleID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Controls                          []*AssetControl                    `gorm:"foreignKey:AssetRegisterID"`
	Contributions                     []*AssetCriticalityContribution    `gorm:"foreignKey:AssetRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Documents                         *[]AssetDocument                   `gorm:"foreignKey:AssetRegisterID"`

	CreatedByUser *User `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	UpdatedByUser *User `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// TableName overrides the default table name.
func (AssetRegister) TableName() string {
	return "asset_register"
}
