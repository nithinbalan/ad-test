package repositories

import (
	"jethur/grc/models"
	riskDTOs "jethur/grc/risk/dtos"

	"gorm.io/gorm"
)

const ismsRiskMatrixType = "isms"

type AssetRiskMatrixRepositoryInterface interface {
	GetMatrix() (*models.RiskMatrix, error)
	GetMatrixDetails() (*models.RiskMatrix, error)
	GetAllImpacts() ([]*riskDTOs.RiskMatrixImpactDTO, error)
	GetAllLikelihoods() ([]*riskDTOs.RiskMatrixLikelihoodDTO, error)
	GetAllRiskRatings() ([]*riskDTOs.RiskMatrixRiskRatingDTO, error)
	GetAllRiskRatingValues() ([]*riskDTOs.RiskMatrixRiskRatingValueListResponseDTO, error)
}

type AssetRiskMatrixRepository struct {
	db *gorm.DB
}

func NewAssetRiskMatrixRepository(db *gorm.DB) *AssetRiskMatrixRepository {
	return &AssetRiskMatrixRepository{
		db: db,
	}
}

func (repo *AssetRiskMatrixRepository) GetMatrix() (*models.RiskMatrix, error) {
	var riskMatrix models.RiskMatrix

	result := repo.db.Where("type = ?", ismsRiskMatrixType).First(&riskMatrix)
	if result.Error != nil {
		return nil, result.Error
	}

	return &riskMatrix, nil
}

func (repo *AssetRiskMatrixRepository) GetMatrixDetails() (*models.RiskMatrix, error) {
	var riskMatrix *models.RiskMatrix

	result := repo.db.
		Preload("RiskMatrixImpact").
		Preload("RiskMatrixLikelihood").
		Preload("RiskMatrixRiskRating.RiskMatrixRiskRatingValue").
		Where("type = ?", ismsRiskMatrixType).
		First(&riskMatrix)
	if result.Error != nil {
		return nil, result.Error
	}

	return riskMatrix, nil
}

func (repo *AssetRiskMatrixRepository) GetAllImpacts() ([]*riskDTOs.RiskMatrixImpactDTO, error) {
	var impacts []*riskDTOs.RiskMatrixImpactDTO

	result := repo.db.Table("risk_matrix_impacts").
		Select("risk_matrix_impacts.*").
		Joins("LEFT JOIN risk_matrix ON risk_matrix.id = risk_matrix_impacts.risk_matrix_id").
		Where("risk_matrix.type = ?", ismsRiskMatrixType).
		Order("risk_matrix_impacts.score asc").
		Find(&impacts)

	return impacts, result.Error
}

func (repo *AssetRiskMatrixRepository) GetAllLikelihoods() ([]*riskDTOs.RiskMatrixLikelihoodDTO, error) {
	var likelihoods []*riskDTOs.RiskMatrixLikelihoodDTO

	result := repo.db.Table("risk_matrix_likelihoods").
		Select("risk_matrix_likelihoods.*").
		Joins("LEFT JOIN risk_matrix ON risk_matrix.id = risk_matrix_likelihoods.risk_matrix_id").
		Where("risk_matrix.type = ?", ismsRiskMatrixType).
		Order("risk_matrix_likelihoods.score asc").
		Find(&likelihoods)

	return likelihoods, result.Error
}

func (repo *AssetRiskMatrixRepository) GetAllRiskRatings() ([]*riskDTOs.RiskMatrixRiskRatingDTO, error) {
	var riskRatings []*riskDTOs.RiskMatrixRiskRatingDTO

	result := repo.db.Table("risk_matrix_risk_ratings").
		Select("risk_matrix_risk_ratings.*").
		Joins("LEFT JOIN risk_matrix ON risk_matrix.id = risk_matrix_risk_ratings.risk_matrix_id").
		Where("risk_matrix.type = ?", ismsRiskMatrixType).
		Order("risk_matrix_risk_ratings.title_en asc").
		Find(&riskRatings)

	return riskRatings, result.Error
}

func (repo *AssetRiskMatrixRepository) GetAllRiskRatingValues() ([]*riskDTOs.RiskMatrixRiskRatingValueListResponseDTO, error) {
	var riskRatingValues []*riskDTOs.RiskMatrixRiskRatingValueListResponseDTO

	result := repo.db.Table("risk_matrix_risk_rating_values").
		Select("risk_matrix_risk_rating_values.*").
		Joins("LEFT JOIN risk_matrix_risk_ratings ON risk_matrix_risk_ratings.id = risk_matrix_risk_rating_values.risk_matrix_risk_rating_id").
		Joins("LEFT JOIN risk_matrix ON risk_matrix.id = risk_matrix_risk_ratings.risk_matrix_id").
		Where("risk_matrix.type = ?", ismsRiskMatrixType).
		Find(&riskRatingValues)

	return riskRatingValues, result.Error
}
