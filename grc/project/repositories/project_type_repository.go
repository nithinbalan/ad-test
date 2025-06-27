package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interface
type ProjectTypeRepositoryInterface interface {
	GetAll(sort, sortBy, keyword, status string) ([]*models.ProjectType, error)
	GetByID(id uint) (*models.ProjectType, error)
	Create(projectType *models.ProjectType) (*models.ProjectType, error)
	Update(id uint, projectType *models.ProjectType) (*models.ProjectType, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

// Implementation
type ProjectTypeRepository struct {
	db *gorm.DB
}

func NewProjectTypeRepository(db *gorm.DB) *ProjectTypeRepository {
	return &ProjectTypeRepository{db: db}
}

func (repo *ProjectTypeRepository) GetAll(sort, sortBy, keyword, status string) ([]*models.ProjectType, error) {
	var projectTypes []*models.ProjectType
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
	result := query.Find(&projectTypes)

	return projectTypes, result.Error
}

func (repo *ProjectTypeRepository) GetByID(id uint) (*models.ProjectType, error) {
	var pt *models.ProjectType
	result := repo.db.First(&pt, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return pt, result.Error
}

func (repo *ProjectTypeRepository) Create(pt *models.ProjectType) (*models.ProjectType, error) {
	result := repo.db.Create(&pt)
	return pt, result.Error
}

func (repo *ProjectTypeRepository) Update(id uint, updatedPT *models.ProjectType) (*models.ProjectType, error) {
	var pt *models.ProjectType
	if err := repo.db.First(&pt, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&pt).Updates(updatedPT).Error; err != nil {
		return nil, err
	}

	return pt, nil
}

func (repo *ProjectTypeRepository) Delete(id uint) error {
	var pt *models.ProjectType
	if err := repo.db.First(&pt, id).Error; err != nil {
		return err
	}
	return repo.db.Delete(&models.ProjectType{}, id).Error
}

func (repo *ProjectTypeRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	if err := repo.db.Model(&models.ProjectType{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.ProjectType{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *ProjectTypeRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	if err := repo.db.Model(&models.ProjectType{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.ProjectType{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *ProjectTypeRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.ProjectType{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error
	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}

func (repo *ProjectTypeRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.ProjectType{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error
	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}
