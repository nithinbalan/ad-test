package repositories

import (
	"jethur/grc/models"
	"jethur/grc/organization/dtos"

	"gorm.io/gorm"
)

type ProcessExternalDependencyEntityTypeRepositoryInterface interface {
	GetAll(params *dtos.ProcessExternalDependencyEntityTypeListQueryParamsDTO) ([]*models.ProcessExternalDependencyEntityType, error)
	GetByID(id uint) (*models.ProcessExternalDependencyEntityType, error)
	Create(entityType *models.ProcessExternalDependencyEntityType) (uint, error)
	Update(entityType *models.ProcessExternalDependencyEntityType) error
	Delete(id uint) error
	ExistsByTitle(titleEn, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

type ProcessExternalDependencyEntityTypeRepository struct {
	db *gorm.DB
}

func NewProcessExternalDependencyEntityTypeRepository(db *gorm.DB) *ProcessExternalDependencyEntityTypeRepository {
	return &ProcessExternalDependencyEntityTypeRepository{
		db: db,
	}
}

func (repo *ProcessExternalDependencyEntityTypeRepository) GetAll(params *dtos.ProcessExternalDependencyEntityTypeListQueryParamsDTO) ([]*models.ProcessExternalDependencyEntityType, error) {
	var entityTypes []*models.ProcessExternalDependencyEntityType

	query := repo.db

	if params.Keyword != "" {
		query = query.Where("(title_en LIKE ? OR title_ar LIKE ?)", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}
	if params.SortBy == "" {
		params.SortBy = "created_at"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}
	if params.Status == "active" || params.Status == "inactive" {
		query = query.Where("status = ?", params.Status)
	}

	query = query.Order(params.SortBy + " " + params.Sort)
	result := query.Find(&entityTypes)

	return entityTypes, result.Error
}

func (repo *ProcessExternalDependencyEntityTypeRepository) GetByID(id uint) (*models.ProcessExternalDependencyEntityType, error) {
	var entityType *models.ProcessExternalDependencyEntityType

	result := repo.db.First(&entityType, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return entityType, nil
}

func (repo *ProcessExternalDependencyEntityTypeRepository) Create(entityType *models.ProcessExternalDependencyEntityType) (uint, error) {
	if err := repo.db.Create(entityType).Error; err != nil {
		return 0, err
	}

	return entityType.ID, nil
}

func (repo *ProcessExternalDependencyEntityTypeRepository) Update(entityType *models.ProcessExternalDependencyEntityType) error {
	return repo.db.Save(entityType).Error
}

func (repo *ProcessExternalDependencyEntityTypeRepository) Delete(id uint) error {
	return repo.db.Delete(&models.ProcessExternalDependencyEntityType{}, id).Error
}

func (repo *ProcessExternalDependencyEntityTypeRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.ProcessExternalDependencyEntityType{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	return count > 0, err
}

func (repo *ProcessExternalDependencyEntityTypeRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.ProcessExternalDependencyEntityType{}).
		Where("id != ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	return count > 0, err
}
