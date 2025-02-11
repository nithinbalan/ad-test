package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type ComplianceFolderRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ComplianceFolder, error)
	GetByID(id uint) (*models.ComplianceFolder, error)
	Create(folder *models.ComplianceFolder) (*models.ComplianceFolder, error)
	Update(id uint, folder *models.ComplianceFolder) (*models.ComplianceFolder, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

// Implementation
type ComplianceFolderRepository struct {
	db *gorm.DB
}

func NewComplianceFolderRepository(db *gorm.DB) *ComplianceFolderRepository {
	return &ComplianceFolderRepository{
		db: db,
	}
}

func (repo *ComplianceFolderRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ComplianceFolder, error) {
	var folders []*models.ComplianceFolder

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
	result := query.Find(&folders)

	return folders, result.Error
}

func (repo *ComplianceFolderRepository) GetByID(id uint) (*models.ComplianceFolder, error) {
	var folder *models.ComplianceFolder
	result := repo.db.First(&folder, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return folder, result.Error
}

func (repo *ComplianceFolderRepository) Create(folder *models.ComplianceFolder) (*models.ComplianceFolder, error) {
	result := repo.db.Create(&folder)

	return folder, result.Error
}

func (repo *ComplianceFolderRepository) Update(id uint, updatedFolder *models.ComplianceFolder) (*models.ComplianceFolder, error) {
	var folder *models.ComplianceFolder

	if err := repo.db.First(&folder, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&folder).Updates(updatedFolder).Error; err != nil {
		return nil, err
	}

	return folder, nil
}

func (repo *ComplianceFolderRepository) Delete(id uint) error {
	var folder *models.ComplianceFolder
	if err := repo.db.First(&folder, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.ComplianceFolder{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *ComplianceFolderRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.ComplianceFolder{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.ComplianceFolder{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil
}

func (repo *ComplianceFolderRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'active' and within the provided 'ids' slice
	if err := repo.db.Model(&models.ComplianceFolder{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.ComplianceFolder{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}

	// Return the IDs of the records that were deactivated
	return idsToUpdate, nil
}

func (repo *ComplianceFolderRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64

	// Check for the existence of titleEn or titleAr in a single query
	err := repo.db.Model(&models.ComplianceFolder{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	// Return true if either title exists, otherwise false
	return count > 0, nil
}

func (repo *ComplianceFolderRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.ComplianceFolder{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
