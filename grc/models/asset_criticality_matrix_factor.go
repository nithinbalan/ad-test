package models

type AssetCriticalityMatrixFactor struct {
	ID                       uint    `json:"id" gorm:"primaryKey;autoIncrement"`
	AssetCriticalityMatrixID uint    `json:"asset_criticality_matrix_id" gorm:"not null"`
	TitleEn                  string  `json:"title_en" gorm:"type:nvarchar(100);not null" binding:"required"`
	TitleAr                  string  `json:"title_ar" gorm:"type:nvarchar(100);default:null" binding:"omitempty"`
	DescriptionEn            *string `json:"description_en" gorm:"type:nvarchar(MAX);default:null" binding:"omitempty"`
	DescriptionAr            *string `json:"description_ar" gorm:"type:nvarchar(MAX);default:null" binding:"omitempty"`
	Weightage                float64 `json:"weightage" gorm:"not null" binding:"required"`

	// Relationships
	Matrix *AssetCriticalityMatrix `gorm:"foreignKey:AssetCriticalityMatrixID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (AssetCriticalityMatrixFactor) TableName() string {
	return "asset_criticality_matrix_factors"
}
