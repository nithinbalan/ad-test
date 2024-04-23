package repositories

import (
	"errors"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/grc/risk/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type RiskMatrixLikelihoodRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*riskDTOs.RiskMatrixLikelihoodDTO, error)
	GetByID(id uint) (*models.RiskMatrixLikelihood, error)
	Create(likelihood *models.RiskMatrixLikelihood) (uint, error)
	Update(likelihood *models.RiskMatrixLikelihood) error
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string, RiskMatrixID uint) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string, RiskMatrixID uint) (bool, error)
}

// implementation
type RiskMatrixLikelihoodRepository struct {
	db *gorm.DB
}

func NewRiskMatrixLikelihoodRepository(db *gorm.DB) *RiskMatrixLikelihoodRepository {
	return &RiskMatrixLikelihoodRepository{
		db: db,
	}
}

func (repo *RiskMatrixLikelihoodRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*riskDTOs.RiskMatrixLikelihoodDTO, error) {

	var likelihoods []*riskDTOs.RiskMatrixLikelihoodDTO

	query := repo.db.Table("risk_matrix_likelihoods").
		Select(`
            risk_matrix_likelihoods.*, risk_matrix.title as risk_matrix_title`).
		Joins(`
		    LEFT JOIN risk_matrix ON risk_matrix.id = risk_matrix_likelihoods.risk_matrix_id
		`)

	if keyword != "" {
		query = query.Where(`
				(risk_matrix_likelihoods.title_en LIKE ? OR
				risk_matrix_likelihoods.title_ar LIKE ? OR
				risk_matrix_likelihoods.description_en LIKE ? OR
				risk_matrix_likelihoods.description_ar LIKE ?)`,
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	if sortBy == "" {
		sortBy = "risk_matrix_likelihoods.title_en"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}
	query = query.Order(sortBy + " " + sort)
	result := query.Find(&likelihoods)

	return likelihoods, result.Error

}

func (repo *RiskMatrixLikelihoodRepository) GetByID(id uint) (*models.RiskMatrixLikelihood, error) {
	var likelihood *models.RiskMatrixLikelihood
	result := repo.db.Preload("RiskMatrix").First(&likelihood, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return likelihood, result.Error
}

func (repo *RiskMatrixLikelihoodRepository) Create(likelihood *models.RiskMatrixLikelihood) (uint, error) {

	if err := repo.db.Create(likelihood).Error; err != nil {
		return 0, err
	}

	return likelihood.ID, nil
}

func (repo *RiskMatrixLikelihoodRepository) Update(updatedLikelihood *models.RiskMatrixLikelihood) error {

	if err := repo.db.Save(updatedLikelihood).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskMatrixLikelihoodRepository) Delete(id uint) error {

	var likelihood *models.RiskMatrixLikelihood
	if err := repo.db.First(&likelihood, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.RiskMatrixLikelihood{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *RiskMatrixLikelihoodRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.RiskMatrixLikelihood{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.RiskMatrixLikelihood{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil
}

func (repo *RiskMatrixLikelihoodRepository) Deactivate(ids []uint) ([]uint, error) {

	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.RiskMatrixLikelihood{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.RiskMatrixLikelihood{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil

}

func (repo *RiskMatrixLikelihoodRepository) ExistsByTitle(titleEn, titleAr string, RiskMatrixID uint) (bool, error) {
	var count int64

	// Check for the existence of titleEn or titleAr in a single query
	err := repo.db.Model(&models.RiskMatrixLikelihood{}).
		Where("(title_en = ? OR title_ar = ?) AND risk_matrix_id = ?", titleEn, titleAr, RiskMatrixID).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)

	}
	// Return true if either title exists, otherwise false
	return count > 0, nil
}

func (repo *RiskMatrixLikelihoodRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string, RiskMatrixID uint) (bool, error) {
	var count int64

	err := repo.db.Model(&models.RiskMatrixLikelihood{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?) AND risk_matrix_id = ?", id, titleEn, titleAr, RiskMatrixID).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
