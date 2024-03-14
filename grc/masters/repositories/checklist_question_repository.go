package repositories

import (
	"errors"
	"jethur/grc/masters/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type ChecklistQuestionRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ChecklistQuestion, error)
	GetByID(id uint) (*models.ChecklistQuestion, error)
	Create(article *models.ChecklistQuestion) (*models.ChecklistQuestion, error)
	Update(id uint, checklistQuestion *models.ChecklistQuestion) (*models.ChecklistQuestion, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

// implementation
type ChecklistQuestionRepository struct {
	db *gorm.DB
}

func NewChecklistQuestionRepository(db *gorm.DB) *ChecklistQuestionRepository {
	return &ChecklistQuestionRepository{
		db: db,
	}
}

func (repo *ChecklistQuestionRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ChecklistQuestion, error) {

	var checklistQuestions []*models.ChecklistQuestion

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
	result := query.Find(&checklistQuestions)

	return checklistQuestions, result.Error
}

func (repo *ChecklistQuestionRepository) GetByID(id uint) (*models.ChecklistQuestion, error) {
	var checklistQuestion *models.ChecklistQuestion
	result := repo.db.First(&checklistQuestion, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return checklistQuestion, result.Error
}

func (repo *ChecklistQuestionRepository) Create(checklistQuestion *models.ChecklistQuestion) (*models.ChecklistQuestion, error) {
	result := repo.db.Create(&checklistQuestion)

	return checklistQuestion, result.Error
}

func (repo *ChecklistQuestionRepository) Update(id uint, updatedChecklistQuestion *models.ChecklistQuestion) (*models.ChecklistQuestion, error) {
	var checklistQuestion *models.ChecklistQuestion

	if err := repo.db.First(&checklistQuestion, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&checklistQuestion).Updates(updatedChecklistQuestion).Error; err != nil {
		return nil, err
	}

	return checklistQuestion, nil
}

func (repo *ChecklistQuestionRepository) Delete(id uint) error {

	var checklistQuestion *models.ChecklistQuestion
	if err := repo.db.First(&checklistQuestion, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.ChecklistQuestion{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *ChecklistQuestionRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.ChecklistQuestion{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.ChecklistQuestion{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil
}

func (repo *ChecklistQuestionRepository) Deactivate(ids []uint) ([]uint, error) {

	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.ChecklistQuestion{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.ChecklistQuestion{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil

}

func (repo *ChecklistQuestionRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64

	// Check for the existence of titleEn or titleAr in a single query
	err := repo.db.Model(&models.ChecklistQuestion{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)

	}
	// Return true if either title exists, otherwise false
	return count > 0, nil
}

func (repo *ChecklistQuestionRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.ChecklistQuestion{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
