package models

type AssetCriticalityMatrixRatingScale struct {
	ID                       uint    `json:"id" gorm:"primaryKey;autoIncrement"`
	AssetCriticalityMatrixID uint    `json:"asset_criticality_matrix_id" gorm:"not null"`
	TitleEn                  string  `json:"title_en" gorm:"type:nvarchar(100);not null" binding:"required"`
	TitleAr                  string  `json:"title_ar" gorm:"type:nvarchar(100);default:null" binding:"omitempty"`
	DescriptionEn            *string `json:"description_en" gorm:"type:nvarchar(MAX);default:null" binding:"omitempty"`
	DescriptionAr            *string `json:"description_ar" gorm:"type:nvarchar(MAX);default:null" binding:"omitempty"`
	ScoreFrom                float64 `json:"score_from" gorm:"not null" binding:"required"`
	ScoreTo                  float64 `json:"score_to" gorm:"not null" binding:"required"`
	Type                     string  `json:"type" gorm:"type:nvarchar(50);default:null" binding:"omitempty"`
	ColorCode                string  `json:"color_code" gorm:"type:nvarchar(25);default:null" binding:"omitempty"`

	// Relationships
	Matrix *AssetCriticalityMatrix `gorm:"foreignKey:AssetCriticalityMatrixID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (AssetCriticalityMatrixRatingScale) TableName() string {
	return "asset_criticality_matrix_rating_scales"
}
