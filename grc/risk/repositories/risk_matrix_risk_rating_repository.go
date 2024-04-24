package repositories

import (
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/grc/risk/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type RiskMatrixRiskRatingRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*riskDTOs.RiskMatrixRiskRatingDTO, error)
	GetByID(id uint) (*models.RiskMatrixRiskRating, error)
	Create(riskRating *models.RiskMatrixRiskRating) (uint, error)
	Update(riskRating *models.RiskMatrixRiskRating) error
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string, RiskMatrixID uint) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string, RiskMatrixID uint) (bool, error)
	ExistsByType(ratingType string, RiskMatrixID uint) (bool, error)
	ExistsByTypeExcludingID(id uint, ratingType string, RiskMatrixID uint) (bool, error)
}

// implementation
type RiskMatrixRiskRatingRepository struct {
	db *gorm.DB
}

func NewRiskMatrixRiskRatingRepository(db *gorm.DB) *RiskMatrixRiskRatingRepository {
	return &RiskMatrixRiskRatingRepository{
		db: db,
	}
}

func (repo *RiskMatrixRiskRatingRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*riskDTOs.RiskMatrixRiskRatingDTO, error) {
	var riskRatings []*riskDTOs.RiskMatrixRiskRatingDTO

	query := repo.db.Table("risk_matrix_risk_ratings").
		Select(`
            risk_matrix_risk_ratings.*, risk_matrix.title as risk_matrix_title`).
		Joins(`
		    LEFT JOIN risk_matrix ON risk_matrix.id = risk_matrix_risk_ratings.risk_matrix_id
		`)

	if keyword != "" {
		query = query.Where(`
				(risk_matrix_risk_ratings.title_en LIKE ? OR
				risk_matrix_risk_ratings.title_ar LIKE ? OR
				risk_matrix_risk_ratings.type LIKE ?)`,
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	if sortBy == "" {
		sortBy = "risk_matrix_risk_ratings.title_en"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}
	query = query.Order(sortBy + " " + sort)
	result := query.Find(&riskRatings)

	return riskRatings, result.Error

}

func (repo *RiskMatrixRiskRatingRepository) GetByID(id uint) (*models.RiskMatrixRiskRating, error) {
	var riskRating *models.RiskMatrixRiskRating
	result := repo.db.Preload("RiskMatrix").First(&riskRating, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return riskRating, result.Error
}

func (repo *RiskMatrixRiskRatingRepository) Create(riskRating *models.RiskMatrixRiskRating) (uint, error) {
	if err := repo.db.Create(riskRating).Error; err != nil {
		return 0, err
	}

	return riskRating.ID, nil
}

func (repo *RiskMatrixRiskRatingRepository) Update(updatedRiskRating *models.RiskMatrixRiskRating) error {
	if err := repo.db.Save(updatedRiskRating).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskMatrixRiskRatingRepository) Delete(id uint) error {
	var riskRating *models.RiskMatrixRiskRating
	if err := repo.db.First(&riskRating, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.RiskMatrixRiskRating{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *RiskMatrixRiskRatingRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.RiskMatrixRiskRating{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.RiskMatrixRiskRating{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil
}

func (repo *RiskMatrixRiskRatingRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.RiskMatrixRiskRating{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.RiskMatrixRiskRating{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil

}

func (repo *RiskMatrixRiskRatingRepository) ExistsByTitle(titleEn, titleAr string, RiskMatrixID uint) (bool, error) {
	var count int64

	// Check for the existence of titleEn or titleAr in a single query
	err := repo.db.Model(&models.RiskMatrixRiskRating{}).
		Where("(title_en = ? OR title_ar = ?) AND risk_matrix_id = ?", titleEn, titleAr, RiskMatrixID).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)

	}
	// Return true if either title exists, otherwise false
	return count > 0, nil
}

func (repo *RiskMatrixRiskRatingRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string, RiskMatrixID uint) (bool, error) {
	var count int64

	err := repo.db.Model(&models.RiskMatrixRiskRating{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?) AND risk_matrix_id = ?", id, titleEn, titleAr, RiskMatrixID).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}

func (repo *RiskMatrixRiskRatingRepository) ExistsByType(ratingType string, RiskMatrixID uint) (bool, error) {
	var count int64

	// Check for the existence of titleEn or titleAr in a single query
	err := repo.db.Model(&models.RiskMatrixRiskRating{}).
		Where("type = ? AND risk_matrix_id = ?", ratingType, RiskMatrixID).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)

	}
	// Return true if either title exists, otherwise false
	return count > 0, nil
}

func (repo *RiskMatrixRiskRatingRepository) ExistsByTypeExcludingID(id uint, ratingType string, RiskMatrixID uint) (bool, error) {
	var count int64

	err := repo.db.Model(&models.RiskMatrixRiskRating{}).
		Where("id <> ? AND type = ? AND risk_matrix_id = ?", id, ratingType, RiskMatrixID).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
