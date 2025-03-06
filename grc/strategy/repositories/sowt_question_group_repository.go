package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type SwotQuestionGroupRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.SwotQuestionGroup, error)
	GetByID(id uint) (*models.SwotQuestionGroup, error)
	Create(swotQuestionGroup *models.SwotQuestionGroup) (*models.SwotQuestionGroup, error)
	Update(id uint, updatedSwotQuestionGroup *models.SwotQuestionGroup) (*models.SwotQuestionGroup, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

// implementation
type SwotQuestionGroupRepository struct {
	db *gorm.DB
}

func NewSwotQuestionGroupRepository(db *gorm.DB) *SwotQuestionGroupRepository {
	return &SwotQuestionGroupRepository{
		db: db,
	}
}

func (repo *SwotQuestionGroupRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.SwotQuestionGroup, error) {
	var swotQuestionGroups []*models.SwotQuestionGroup

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
	result := query.Find(&swotQuestionGroups)

	return swotQuestionGroups, result.Error
}

func (repo *SwotQuestionGroupRepository) GetByID(id uint) (*models.SwotQuestionGroup, error) {
	var swotQuestionGroup *models.SwotQuestionGroup
	result := repo.db.First(&swotQuestionGroup, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return swotQuestionGroup, nil
}

func (repo *SwotQuestionGroupRepository) Create(swotQuestionGroup *models.SwotQuestionGroup) (*models.SwotQuestionGroup, error) {
	result := repo.db.Create(&swotQuestionGroup)

	return swotQuestionGroup, result.Error
}

func (repo *SwotQuestionGroupRepository) Update(id uint, updatedSwotQuestionGroup *models.SwotQuestionGroup) (*models.SwotQuestionGroup, error) {
	var swotQuestionGroup *models.SwotQuestionGroup

	if err := repo.db.First(&swotQuestionGroup, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&swotQuestionGroup).Updates(updatedSwotQuestionGroup).Error; err != nil {
		return nil, err
	}

	return swotQuestionGroup, nil
}

func (repo *SwotQuestionGroupRepository) Delete(id uint) error {

	var swotQuestionGroup *models.SwotQuestionGroup
	if err := repo.db.First(&swotQuestionGroup, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.SwotQuestionGroup{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *SwotQuestionGroupRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.SwotQuestionGroup{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.SwotQuestionGroup{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *SwotQuestionGroupRepository) Deactivate(ids []uint) ([]uint, error) {

	var idsToUpdate []uint

	if err := repo.db.Model(&models.SwotQuestionGroup{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.SwotQuestionGroup{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	return idsToUpdate, nil

}

func (repo *SwotQuestionGroupRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.SwotQuestionGroup{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)

	}
	return count > 0, nil
}

func (repo *SwotQuestionGroupRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.SwotQuestionGroup{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
