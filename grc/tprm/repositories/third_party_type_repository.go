package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type ThirdPartyTypeRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ThirdPartyType, error)
	GetByID(id uint) (*models.ThirdPartyType, error)
	Create(article *models.ThirdPartyType) (*models.ThirdPartyType, error)
	Update(id uint, thirdPartyType *models.ThirdPartyType) (*models.ThirdPartyType, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

// implementation
type ThirdPartyTypeRepository struct {
	db *gorm.DB
}

func NewThirdPartyTypeRepository(db *gorm.DB) *ThirdPartyTypeRepository {
	return &ThirdPartyTypeRepository{
		db: db,
	}
}

func (repo *ThirdPartyTypeRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ThirdPartyType, error) {
	var thirdPartyTypes []*models.ThirdPartyType
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
	result := query.Find(&thirdPartyTypes)

	return thirdPartyTypes, result.Error
}

func (repo *ThirdPartyTypeRepository) GetByID(id uint) (*models.ThirdPartyType, error) {
	var thirdPartyType *models.ThirdPartyType
	result := repo.db.First(&thirdPartyType, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return thirdPartyType, nil
}

func (repo *ThirdPartyTypeRepository) Create(thirdPartyType *models.ThirdPartyType) (*models.ThirdPartyType, error) {
	result := repo.db.Create(&thirdPartyType)

	return thirdPartyType, result.Error
}

func (repo *ThirdPartyTypeRepository) Update(id uint, updatedType *models.ThirdPartyType) (*models.ThirdPartyType, error) {
	var thirdPartyType *models.ThirdPartyType
	if err := repo.db.First(&thirdPartyType, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&thirdPartyType).Updates(updatedType).Error; err != nil {
		return nil, err
	}

	return thirdPartyType, nil
}

func (repo *ThirdPartyTypeRepository) Delete(id uint) error {
	var thirdPartyType *models.ThirdPartyType
	if err := repo.db.First(&thirdPartyType, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.ThirdPartyType{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *ThirdPartyTypeRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.ThirdPartyType{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.ThirdPartyType{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}
	// Return the IDs of the records that were activated
	return idsToUpdate, nil
}

func (repo *ThirdPartyTypeRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.ThirdPartyType{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.ThirdPartyType{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}
	return idsToUpdate, nil
}

func (repo *ThirdPartyTypeRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64

	// Check for the existence of titleEn or titleAr in a single query
	err := repo.db.Model(&models.ThirdPartyType{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)

	}
	return count > 0, nil
}

func (repo *ThirdPartyTypeRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.ThirdPartyType{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
