package repositories

import (
	"errors"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/grc/risk/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type RiskMatrixImpactRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*riskDTOs.RiskMatrixImpactDTO, error)
	GetByID(id uint) (*models.RiskMatrixImpact, error)
	Create(impact *models.RiskMatrixImpact) (uint, error)
	Update(impact *models.RiskMatrixImpact) error
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string, RiskMatrixID uint) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string, RiskMatrixID uint) (bool, error)
}

// implementation
type RiskMatrixImpactRepository struct {
	db *gorm.DB
}

func NewRiskMatrixImpactRepository(db *gorm.DB) *RiskMatrixImpactRepository {
	return &RiskMatrixImpactRepository{
		db: db,
	}
}

func (repo *RiskMatrixImpactRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*riskDTOs.RiskMatrixImpactDTO, error) {
	var impacts []*riskDTOs.RiskMatrixImpactDTO

	query := repo.db.Table("risk_matrix_impacts").
		Select(`
            risk_matrix_impacts.*, risk_matrix.title as risk_matrix_title`).
		Joins(`
		    LEFT JOIN risk_matrix ON risk_matrix.id = risk_matrix_impacts.risk_matrix_id
		`)

	if keyword != "" {
		query = query.Where(`
				(risk_matrix_impacts.title_en LIKE ? OR
				risk_matrix_impacts.title_ar LIKE ? OR
				risk_matrix_impacts.description_en LIKE ? OR
				risk_matrix_impacts.description_ar LIKE ?)`,
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	if sortBy == "" {
		sortBy = "risk_matrix_impacts.title_en"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}
	query = query.Order(sortBy + " " + sort)
	result := query.Find(&impacts)

	return impacts, result.Error

}

func (repo *RiskMatrixImpactRepository) GetByID(id uint) (*models.RiskMatrixImpact, error) {
	var impact *models.RiskMatrixImpact
	result := repo.db.Preload("RiskMatrix").First(&impact, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return impact, result.Error
}

func (repo *RiskMatrixImpactRepository) Create(impact *models.RiskMatrixImpact) (uint, error) {
	if err := repo.db.Create(impact).Error; err != nil {
		return 0, err
	}

	return impact.ID, nil
}

func (repo *RiskMatrixImpactRepository) Update(updatedImpact *models.RiskMatrixImpact) error {
	if err := repo.db.Save(updatedImpact).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskMatrixImpactRepository) Delete(id uint) error {
	var impact *models.RiskMatrixImpact
	if err := repo.db.First(&impact, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.RiskMatrixImpact{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *RiskMatrixImpactRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.RiskMatrixImpact{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.RiskMatrixImpact{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil
}

func (repo *RiskMatrixImpactRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.RiskMatrixImpact{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.RiskMatrixImpact{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil

}

func (repo *RiskMatrixImpactRepository) ExistsByTitle(titleEn, titleAr string, RiskMatrixID uint) (bool, error) {
	var count int64

	// Check for the existence of titleEn or titleAr in a single query
	err := repo.db.Model(&models.RiskMatrixImpact{}).
		Where("(title_en = ? OR title_ar = ?) AND risk_matrix_id = ?", titleEn, titleAr, RiskMatrixID).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)

	}
	// Return true if either title exists, otherwise false
	return count > 0, nil
}

func (repo *RiskMatrixImpactRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string, RiskMatrixID uint) (bool, error) {
	var count int64

	err := repo.db.Model(&models.RiskMatrixImpact{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?) AND risk_matrix_id = ?", id, titleEn, titleAr, RiskMatrixID).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
