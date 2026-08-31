package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ProjectClassificationRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ProjectClassification, error)
	GetByID(id uint) (*models.ProjectClassification, error)
	Create(projectClassification *models.ProjectClassification) (*models.ProjectClassification, error)
	Update(id uint, projectClassification *models.ProjectClassification) (*models.ProjectClassification, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

type ProjectClassificationRepository struct {
	db *gorm.DB
}

func NewProjectClassificationRepository(db *gorm.DB) *ProjectClassificationRepository {
	return &ProjectClassificationRepository{
		db: db,
	}
}

func (repo *ProjectClassificationRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ProjectClassification, error) {
	var projectClassifications []*models.ProjectClassification
	query := repo.db

	if keyword != "" {
		query = query.Where(
			"(title_en LIKE ? OR title_ar LIKE ?)",
			"%"+keyword+"%",
			"%"+keyword+"%",
		)
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

	result := query.Order(sortBy + " " + sort).Find(&projectClassifications)

	return projectClassifications, result.Error
}

func (repo *ProjectClassificationRepository) GetByID(id uint) (*models.ProjectClassification, error) {
	var projectClassification models.ProjectClassification

	if err := repo.db.First(&projectClassification, id).Error; err != nil {
		return nil, err
	}

	return &projectClassification, nil
}

func (repo *ProjectClassificationRepository) Create(projectClassification *models.ProjectClassification) (*models.ProjectClassification, error) {
	if err := repo.db.Create(projectClassification).Error; err != nil {
		return nil, err
	}

	return projectClassification, nil
}

func (repo *ProjectClassificationRepository) Update(id uint, updatedClassification *models.ProjectClassification) (*models.ProjectClassification, error) {
	var projectClassification models.ProjectClassification

	if err := repo.db.First(&projectClassification, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&projectClassification).Updates(updatedClassification).Error; err != nil {
		return nil, err
	}

	return &projectClassification, nil
}

func (repo *ProjectClassificationRepository) Delete(id uint) error {
	var projectClassification models.ProjectClassification

	if err := repo.db.First(&projectClassification, id).Error; err != nil {
		return err
	}

	return repo.db.Delete(&projectClassification).Error
}

func (repo *ProjectClassificationRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.ProjectClassification{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.ProjectClassification{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *ProjectClassificationRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.ProjectClassification{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.ProjectClassification{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *ProjectClassificationRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.ProjectClassification{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}

func (repo *ProjectClassificationRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.ProjectClassification{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}