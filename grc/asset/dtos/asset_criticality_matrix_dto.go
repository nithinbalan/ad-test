package dtos

type AssetCriticalityMatrixResponseDTO struct {
	ID              uint                                          `json:"id"`
	TitleEn         string                                        `json:"title_en"`
	TitleAr         string                                        `json:"title_ar"`
	Status          string                                        `json:"status"`
	Factors         []*AssetCriticalityMatrixFactorDTO            `json:"factors"`
	FactorRatings   []*AssetCriticalityMatrixFactorRatingDTO      `json:"factor_ratings"`
	RatingScales    []*AssetCriticalityMatrixRatingScaleDTO       `json:"rating_scales"`
	AssetCategories []*AssetCriticalityMatrixAssetCategoryListDTO `json:"asset_categories"`
}

type AssetCriticalityMatrixDetailDTO struct {
	ID              uint                                          `json:"id"`
	TitleEn         string                                        `json:"title_en"`
	TitleAr         string                                        `json:"title_ar"`
	Status          string                                        `json:"status"`
	Factors         []*AssetCriticalityMatrixFactorDTO            `json:"factors"`
	FactorRatings   []*AssetCriticalityMatrixFactorRatingDTO      `json:"factor_ratings"`
	RatingScales    []*AssetCriticalityMatrixRatingScaleDTO       `json:"rating_scales"`
	AssetCategories []*AssetCriticalityMatrixAssetCategoryListDTO `json:"asset_categories"`
}

type AssetCriticalityMatrixFormDTO struct {
	TitleEn          string `json:"title_en" binding:"required"`
	TitleAr          string `json:"title_ar" binding:"required"`
	Status           string `json:"status" binding:"required"`
	AssetCategoryIDs []uint `json:"asset_category_ids" binding:"required"`
}

type AssetCriticalityMatrixAssetCategoryListDTO struct {
	ID                       uint             `json:"id"`
	AssetCriticalityMatrixID uint             `json:"asset_criticality_matrix_id"`
	AssetCategoryID          uint             `json:"asset_category_id"`
	AssetCategory            AssetCategoryDTO `json:"asset_category"`
}
