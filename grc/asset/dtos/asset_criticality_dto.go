package dtos

import "time"

type AssetCriticalityContributionListDTO struct {
	ID                                        uint       `json:"id"`
	AssetRegisterID                           uint       `json:"asset_register_id"`
	AssetRegisterTitle                        string     `json:"asset_register_title"`
	AssetCriticalityMatrixFactorID            uint       `json:"asset_criticality_matrix_factor_id"`
	AssetCriticalityMatrixFactorTitleEn       string     `json:"asset_criticality_matrix_factor_title_en"`
	AssetCriticalityMatrixFactorTitleAr       string     `json:"asset_criticality_matrix_factor_title_ar"`
	AssetCriticalityMatrixFactorDescriptionEn *string    `json:"asset_criticality_matrix_factor_description_en"`
	AssetCriticalityMatrixFactorDescriptionAr *string    `json:"asset_criticality_matrix_factor_description_ar"`
	AssetCriticalityMatrixFactorWeightage     float64    `json:"asset_criticality_matrix_factor_weightage"`
	AssetCriticalityMatrixID                  uint       `json:"asset_criticality_matrix_id"`
	AssetCriticalityMatrixFactorRatingID      uint       `json:"asset_criticality_matrix_factor_rating_id"`
	AssetCriticalityMatrixFactorRatingTitleEn string     `json:"asset_criticality_matrix_factor_rating_title_en"`
	AssetCriticalityMatrixFactorRatingTitleAr string     `json:"asset_criticality_matrix_factor_rating_title_ar"`
	AssetCriticalityMatrixFactorRatingScore   int        `json:"asset_criticality_matrix_factor_rating_score"`
	AssetCriticalityMatrixFactorRatingType    string     `json:"asset_criticality_matrix_factor_rating_type"`
	AssetCriticalityMatrixFactorRatingColor   string     `json:"asset_criticality_matrix_factor_rating_color"`
	Score                                     float64    `json:"score"`
	Comments                                  *string    `json:"comments"`
	CreatedByID                               uint       `json:"created_by_id"`
	CreatedByFirstNameEn                      string     `json:"created_by_first_name_en"`
	CreatedByFirstNameAr                      string     `json:"created_by_first_name_ar"`
	CreatedByLastNameEn                       string     `json:"created_by_last_name_en"`
	CreatedByLastNameAr                       string     `json:"created_by_last_name_ar"`
	CreatedByInitial                          string     `json:"created_by_initial"`
	CreatedByColor                            string     `json:"created_by_color"`
	CreatedByEmail                            string     `json:"created_by_email"`
	CreatedByProfileImagePath                 string     `json:"created_by_profile_image_path"`
	CreatedAt                                 time.Time  `json:"created_at"`
	UpdatedByID                               *uint      `json:"updated_by_id"`
	UpdatedByFirstNameEn                      *string    `json:"updated_by_first_name_en"`
	UpdatedByFirstNameAr                      *string    `json:"updated_by_first_name_ar"`
	UpdatedByLastNameEn                       *string    `json:"updated_by_last_name_en"`
	UpdatedByLastNameAr                       *string    `json:"updated_by_last_name_ar"`
	UpdatedByInitial                          *string    `json:"updated_by_initial"`
	UpdatedByColor                            *string    `json:"updated_by_color"`
	UpdatedByEmail                            *string    `json:"updated_by_email"`
	UpdatedByProfileImagePath                 *string    `json:"updated_by_profile_image_path"`
	UpdatedAt                                 *time.Time `json:"updated_at"`
}

type AssetCriticalityFormDTO struct {
	AssetCriticalityMatrixFactorID       uint    `json:"asset_criticality_matrix_factor_id" binding:"required"`
	AssetCriticalityMatrixFactorRatingID uint    `json:"asset_criticality_matrix_factor_rating_id" binding:"required"`
	Comments                             *string `json:"comments,omitempty"`
}

type AssetCriticalityContributionListQueryParamsDTO struct {
	Sort    string `json:"sort_order"`
	SortBy  string `json:"sort_by"`
	Keyword string `json:"keyword"`
}
