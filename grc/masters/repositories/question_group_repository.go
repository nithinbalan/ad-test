package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type QuestionGroupRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.QuestionGroup, error)
	GetByID(id uint) (*models.QuestionGroup, error)
	Create(article *models.QuestionGroup) (*models.QuestionGroup, error)
	Update(id uint, questionGroup *models.QuestionGroup) (*models.QuestionGroup, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

// implementation
type QuestionGroupRepository struct {
	db *gorm.DB
}

func NewQuestionGroupRepository(db *gorm.DB) *QuestionGroupRepository {
	return &QuestionGroupRepository{
		db: db,
	}
}

func (repo *QuestionGroupRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.QuestionGroup, error) {

	var questionGroups []*models.QuestionGroup

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
	result := query.Find(&questionGroups)

	return questionGroups, result.Error
}

func (repo *QuestionGroupRepository) GetByID(id uint) (*models.QuestionGroup, error) {
	var questionGroup *models.QuestionGroup
	result := repo.db.First(&questionGroup, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return questionGroup, result.Error
}

func (repo *QuestionGroupRepository) Create(questionGroup *models.QuestionGroup) (*models.QuestionGroup, error) {
	result := repo.db.Create(&questionGroup)

	return questionGroup, result.Error
}

func (repo *QuestionGroupRepository) Update(id uint, updatedQuestionGroup *models.QuestionGroup) (*models.QuestionGroup, error) {
	var questionGroup *models.QuestionGroup

	if err := repo.db.First(&questionGroup, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&questionGroup).Updates(updatedQuestionGroup).Error; err != nil {
		return nil, err
	}

	return questionGroup, nil
}

func (repo *QuestionGroupRepository) Delete(id uint) error {

	var questionGroup *models.QuestionGroup
	if err := repo.db.First(&questionGroup, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.QuestionGroup{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *QuestionGroupRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.QuestionGroup{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.QuestionGroup{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil
}

func (repo *QuestionGroupRepository) Deactivate(ids []uint) ([]uint, error) {

	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.QuestionGroup{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.QuestionGroup{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil

}

func (repo *QuestionGroupRepository) ExistsByTitle(titleEn string, titleAr string) (bool, error) {
	var count int64

	// Check for the existence of titleEn or titleAr in a single query
	err := repo.db.Model(&models.QuestionGroup{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)

	}
	// Return true if either title exists, otherwise false
	return count > 0, nil
}

func (repo *QuestionGroupRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.QuestionGroup{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
