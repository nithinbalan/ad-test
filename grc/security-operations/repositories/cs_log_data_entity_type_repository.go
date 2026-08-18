package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type CSLogDataEntityTypeRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.CSLogDataEntityType, error)
	GetByID(id uint) (*models.CSLogDataEntityType, error)
	Create(entityType *models.CSLogDataEntityType) (*models.CSLogDataEntityType, error)
	Update(id uint, entityType *models.CSLogDataEntityType) (*models.CSLogDataEntityType, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

type CSLogDataEntityTypeRepository struct {
	db *gorm.DB
}

func NewCSLogDataEntityTypeRepository(db *gorm.DB) *CSLogDataEntityTypeRepository {
	return &CSLogDataEntityTypeRepository{
		db: db,
	}
}

func (repo *CSLogDataEntityTypeRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.CSLogDataEntityType, error) {
	var entityTypes []*models.CSLogDataEntityType
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
	result := query.Find(&entityTypes)

	return entityTypes, result.Error
}

func (repo *CSLogDataEntityTypeRepository) GetByID(id uint) (*models.CSLogDataEntityType, error) {
	var entityType *models.CSLogDataEntityType
	result := repo.db.First(&entityType, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return entityType, nil
}

func (repo *CSLogDataEntityTypeRepository) Create(entityType *models.CSLogDataEntityType) (*models.CSLogDataEntityType, error) {
	result := repo.db.Create(&entityType)

	return entityType, result.Error
}

func (repo *CSLogDataEntityTypeRepository) Update(id uint, updated *models.CSLogDataEntityType) (*models.CSLogDataEntityType, error) {
	var entityType *models.CSLogDataEntityType
	if err := repo.db.First(&entityType, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&entityType).Updates(updated).Error; err != nil {
		return nil, err
	}

	return entityType, nil
}

func (repo *CSLogDataEntityTypeRepository) Delete(id uint) error {
	var entityType *models.CSLogDataEntityType
	if err := repo.db.First(&entityType, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.CSLogDataEntityType{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *CSLogDataEntityTypeRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	if err := repo.db.Model(&models.CSLogDataEntityType{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.CSLogDataEntityType{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}
	return idsToUpdate, nil
}

func (repo *CSLogDataEntityTypeRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	if err := repo.db.Model(&models.CSLogDataEntityType{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.CSLogDataEntityType{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}
	return idsToUpdate, nil
}

func (repo *CSLogDataEntityTypeRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.CSLogDataEntityType{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}

func (repo *CSLogDataEntityTypeRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.CSLogDataEntityType{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}
