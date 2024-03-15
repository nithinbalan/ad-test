package repositories

import (
	"errors"
	"jethur/grc/masters/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type ChecklistQuestionRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string) ([]*models.ChecklistQuestionWithChecklistTitle, error)
	GetByID(id uint) (*models.ChecklistQuestion, error)
	Create(article *models.ChecklistQuestion) (*models.ChecklistQuestion, error)
	Update(id uint, checklistQuestion *models.ChecklistQuestion) (*models.ChecklistQuestion, error)
	Delete(id uint) error
	ExistsByTitle(QuestionEn string, QuestionAr string, ChecklistID uint) (bool, error)
	ExistsQuestionLibrary(QuestionLibraryID uint) (bool, error)
	ExistsByTitleExcludingID(id uint, QuestionEn, QuestionAr string, ChecklistID uint) (bool, error)
	ExistsQuestionLibraryExcludingID(id uint, QuestionLibraryID uint) (bool, error)
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

func (repo *ChecklistQuestionRepository) GetAll(sort string, sortBy string, keyword string) ([]*models.ChecklistQuestionWithChecklistTitle, error) {

	var checklistQuestions []*models.ChecklistQuestionWithChecklistTitle

	query := repo.db.Table("checklist_questions").
		Select(
			"checklist_questions.id",
			"checklist_questions.checklist_id",
			"checklist_questions.question_library_id",
			"checklist_questions.question_en",
			"checklist_questions.question_ar",
			"checklists.title_en AS checklist_title_en",
			"checklists.title_ar AS checklist_title_ar",
		).
		Joins("INNER JOIN checklists ON checklist_questions.checklist_id = checklists.id")

	if keyword != "" {
		query = query.Where("(checklist_questions.question_en LIKE ? OR checklist_questions.question_ar LIKE ?)", "%"+keyword+"%", "%"+keyword+"%")
	}
	if sortBy == "" {
		sortBy = "checklist_questions.question_en"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	query = query.Order(sortBy + " " + sort)
	result := query.Find(&checklistQuestions)

	return checklistQuestions, result.Error
}

func (repo *ChecklistQuestionRepository) GetByID(id uint) (*models.ChecklistQuestion, error) {
	var checklistQuestion *models.ChecklistQuestion
	result := repo.db.Preload("Checklist").Preload("QuestionLibrary").First(&checklistQuestion, id)

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

func (repo *ChecklistQuestionRepository) ExistsByTitle(QuestionEn, QuestionAr string, ChecklistID uint) (bool, error) {
	var count int64

	// Check for the existence of QuestionEn or QuestionAr in a single query
	err := repo.db.Model(&models.ChecklistQuestion{}).
		Where("(question_en = ? OR question_ar = ?) AND checklist_id = ?", QuestionEn, QuestionAr, ChecklistID).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)

	}
	// Return true if either title exists, otherwise false
	return count > 0, nil
}

func (repo *ChecklistQuestionRepository) ExistsQuestionLibrary(QuestionLibraryID uint) (bool, error) {
	var count int64

	err := repo.db.Model(&models.ChecklistQuestion{}).
		Where("question_library_id = ?", QuestionLibraryID).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)

	}
	return count > 0, nil
}

func (repo *ChecklistQuestionRepository) ExistsByTitleExcludingID(id uint, QuestionEn, QuestionAr string, ChecklistID uint) (bool, error) {
	var count int64

	err := repo.db.Model(&models.ChecklistQuestion{}).
		Where("id <> ? AND (question_en = ? OR question_ar = ?) AND checklist_id = ?", id, QuestionEn, QuestionAr, ChecklistID).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}

func (repo *ChecklistQuestionRepository) ExistsQuestionLibraryExcludingID(id uint, QuestionLibraryID uint) (bool, error) {
	var count int64

	err := repo.db.Model(&models.ChecklistQuestion{}).
		Where("id <> ? AND question_library_id = ?", id, QuestionLibraryID).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)

	}
	return count > 0, nil
}
