package models

type AssetCriticalityMatrix struct {
	ID      uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	TitleEn string `json:"title_en" gorm:"type:nvarchar(100);not null" binding:"required"`
	TitleAr string `json:"title_ar" gorm:"type:nvarchar(100);default:null" binding:"omitempty"`
	Status  string `json:"status" gorm:"type:varchar(10);default:'active'" binding:"required,oneof=active inactive"`

	// Relationships
	Factors       []*AssetCriticalityMatrixFactor       `gorm:"foreignKey:AssetCriticalityMatrixID"`
	FactorRatings []*AssetCriticalityMatrixFactorRating `gorm:"foreignKey:AssetCriticalityMatrixID"`
	RatingScales  []*AssetCriticalityMatrixRatingScale  `gorm:"foreignKey:AssetCriticalityMatrixID"`
}

func (AssetCriticalityMatrix) TableName() string {
	return "asset_criticality_matrices"
}
