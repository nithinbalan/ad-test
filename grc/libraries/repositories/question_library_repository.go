package repositories

import (
	"errors"
	"jethur/grc/libraries/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type QuestionLibraryRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string) ([]*models.QuestionLibraryWithGroupTitle, error)
	GetByID(id uint) (*models.QuestionLibrary, error)
	Create(article *models.QuestionLibrary) (*models.QuestionLibrary, error)
	Update(id uint, questionLibrary *models.QuestionLibrary) (*models.QuestionLibrary, error)
	Delete(id uint) error
	ExistsByTitle(QuestionEn string, QuestionAr string, QuestionGroupID uint) (bool, error)
	ExistsByTitleExcludingID(id uint, QuestionEn, QuestionAr string, QuestionGroupID uint) (bool, error)
}

// implementation
type QuestionLibraryRepository struct {
	db *gorm.DB
}

func NewQuestionLibraryRepository(db *gorm.DB) *QuestionLibraryRepository {
	return &QuestionLibraryRepository{
		db: db,
	}
}

func (repo *QuestionLibraryRepository) GetAll(sort string, sortBy string, keyword string) ([]*models.QuestionLibraryWithGroupTitle, error) {

	var questionLibrary []*models.QuestionLibraryWithGroupTitle

	query := repo.db.Table("question_library").
		Select(
			"question_library.id",
			"question_library.question_group_id",
			"question_library.question_en",
			"question_library.question_ar",
			"question_library.question_group_id",
			"question_groups.title_en AS question_group_title_en",
			"question_groups.title_ar AS question_group_title_ar",
		).
		Joins("INNER JOIN question_groups ON question_library.question_group_id = question_groups.id")

	if keyword != "" {
		query = query.Where("(question_en LIKE ? OR question_en LIKE ?)", "%"+keyword+"%", "%"+keyword+"%")
	}
	if sortBy == "" {
		sortBy = "question_en"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	// Order the results
	query = query.Order("question_library." + sortBy + " " + sort)

	// Execute the query and return result/error
	result := query.Find(&questionLibrary)

	return questionLibrary, result.Error

}

func (repo *QuestionLibraryRepository) GetByID(id uint) (*models.QuestionLibrary, error) {
	var questionLibrary *models.QuestionLibrary
	result := repo.db.Preload("QuestionGroup").First(&questionLibrary, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return questionLibrary, result.Error
}

func (repo *QuestionLibraryRepository) Create(questionLibrary *models.QuestionLibrary) (*models.QuestionLibrary, error) {
	result := repo.db.Create(&questionLibrary)

	return questionLibrary, result.Error
}

func (repo *QuestionLibraryRepository) Update(id uint, updatedQuestionLibrary *models.QuestionLibrary) (*models.QuestionLibrary, error) {
	var questionLibrary *models.QuestionLibrary

	if err := repo.db.First(&questionLibrary, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&questionLibrary).Updates(updatedQuestionLibrary).Error; err != nil {
		return nil, err
	}

	return questionLibrary, nil
}

func (repo *QuestionLibraryRepository) Delete(id uint) error {

	var questionLibrary *models.QuestionLibrary
	if err := repo.db.First(&questionLibrary, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.QuestionLibrary{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *QuestionLibraryRepository) ExistsByTitle(QuestionEn string, QuestionAr string, QuestionGroupID uint) (bool, error) {
	var count int64

	// Check for the existence of QuestionEn or QuestionAr in a single query
	err := repo.db.Model(&models.QuestionLibrary{}).
		Where("(question_en = ? OR question_ar = ?) AND question_group_id = ?", QuestionEn, QuestionAr, QuestionGroupID).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)

	}
	// Return true if either question exists, otherwise false
	return count > 0, nil
}

func (repo *QuestionLibraryRepository) ExistsByTitleExcludingID(id uint, QuestionEn, QuestionAr string, QuestionGroupID uint) (bool, error) {
	var count int64
	err := repo.db.Model(&models.QuestionLibrary{}).
		Where("id <> ? AND (question_en = ? OR question_ar = ?) AND question_group_id = ?", id, QuestionEn, QuestionAr, QuestionGroupID).
		Count(&count).Error
	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
