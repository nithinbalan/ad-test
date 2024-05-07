package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type RiskSourceRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.RiskSource, error)
	GetByID(id uint) (*models.RiskSource, error)
	Create(article *models.RiskSource) (*models.RiskSource, error)
	Update(id uint, source *models.RiskSource) (*models.RiskSource, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

// implementation
type RiskSourceRepository struct {
	db *gorm.DB
}

func NewRiskSourceRepository(db *gorm.DB) *RiskSourceRepository {
	return &RiskSourceRepository{
		db: db,
	}
}

func (repo *RiskSourceRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.RiskSource, error) {

	var riskSources []*models.RiskSource

	query := repo.db

	if keyword != "" {
		query = query.Where("(title_en LIKE ? OR title_ar LIKE ?)", "%"+keyword+"%", "%"+keyword+"%")
	}
	if sortBy == "" {
		sortBy = "title_en"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	if status == "active" || status == "inactive" {
		query = query.Where("status = ?", status)
	}

	query = query.Order(sortBy + " " + sort)
	result := query.Find(&riskSources)

	return riskSources, result.Error
}

func (repo *RiskSourceRepository) GetByID(id uint) (*models.RiskSource, error) {
	var source *models.RiskSource

	result := repo.db.First(&source, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return source, nil
}

func (repo *RiskSourceRepository) Create(source *models.RiskSource) (*models.RiskSource, error) {
	result := repo.db.Create(&source)

	return source, result.Error
}

func (repo *RiskSourceRepository) Update(id uint, updatedSource *models.RiskSource) (*models.RiskSource, error) {
	var source *models.RiskSource

	if err := repo.db.First(&source, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&source).Updates(updatedSource).Error; err != nil {
		return nil, err
	}

	return source, nil
}

func (repo *RiskSourceRepository) Delete(id uint) error {

	var source *models.RiskSource
	if err := repo.db.First(&source, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.RiskSource{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *RiskSourceRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.RiskSource{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.RiskSource{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil
}

func (repo *RiskSourceRepository) Deactivate(ids []uint) ([]uint, error) {

	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.RiskSource{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.RiskSource{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil

}

func (repo *RiskSourceRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64

	// Check for the existence of titleEn or titleAr in a single query
	err := repo.db.Model(&models.RiskSource{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)

	}
	// Return true if either title exists, otherwise false
	return count > 0, nil
}

func (repo *RiskSourceRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.RiskSource{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
