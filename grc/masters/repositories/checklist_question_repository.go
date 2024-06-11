package repositories

import (
	"errors"
	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type ChecklistQuestionRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, page int, pageSize int, filters map[string][]string) (*masterDTOs.ChecklistQuestionsGroupByQuestionGroupListResponseDTO, error)
	GetByID(id uint) (*models.ChecklistQuestion, error)
	Create(checklistQuestion *models.ChecklistQuestion) (uint, error)
	Update(id uint, checklistQuestion *models.ChecklistQuestion) error
	Delete(id uint) error
	ExistsByTitle(QuestionEn string, QuestionAr string, ChecklistID uint, QuestionGroupID uint) (bool, error)
	ExistsQuestionLibrary(QuestionLibraryID uint, QuestionEn, QuestionAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, QuestionEn, QuestionAr string, ChecklistID uint, QuestionGroupID uint) (bool, error)
	ExistsQuestionLibraryExcludingID(id uint, QuestionEn, QuestionAr string, QuestionLibraryID uint) (bool, error)
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

func (repo *ChecklistQuestionRepository) GetAll(sort string, sortBy string, keyword string, page int, pageSize int, filters map[string][]string) (*masterDTOs.ChecklistQuestionsGroupByQuestionGroupListResponseDTO, error) {

	var response []*masterDTOs.ChecklistQuestionsGroupByQuestionGroupDTO
	var totalCount int64

	query := repo.db.Table("checklist_questions").
		Select(
			"checklist_questions.id",
			"checklist_questions.checklist_id",
			"checklist_questions.question_en",
			"checklist_questions.question_ar",
			"checklists.title_en AS checklist_title_en",
			"checklists.title_ar AS checklist_title_ar",
			"question_groups.title_en AS question_group_title_en",
			"question_groups.title_ar AS question_group_title_ar",
			"question_groups.id AS question_group_id",
		).
		Joins("INNER JOIN checklists ON checklist_questions.checklist_id = checklists.id").
		Joins("INNER JOIN question_groups ON checklist_questions.question_group_id = question_groups.id")

	if keyword != "" {
		query = query.Where("(checklist_questions.question_en LIKE ? OR checklist_questions.question_ar LIKE ?)", "%"+keyword+"%", "%"+keyword+"%")
	}
	if sortBy == "" {
		sortBy = "checklist_questions.question_en"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	// Apply dynamic filters to the query
	for columnName, values := range filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	query = query.Order(sortBy + " " + sort)

	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)

	rows, err := query.Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	checklistQuestions := make(map[uint]*masterDTOs.ChecklistQuestionsGroupByQuestionGroupDTO)

	for rows.Next() {
		var questions masterDTOs.ChecklistQuestionListDTO
		err := repo.db.ScanRows(rows, &questions)
		if err != nil {
			return nil, err
		}

		if _, ok := checklistQuestions[questions.QuestionGroupID]; !ok {
			checklistQuestions[questions.QuestionGroupID] = &masterDTOs.ChecklistQuestionsGroupByQuestionGroupDTO{
				QuestionGroupID:      questions.QuestionGroupID,
				QuestionGroupTitleEn: questions.QuestionGroupTitleEn,
				QuestionGroupTitleAr: questions.QuestionGroupTitleAr,
			}
		}

		checklistQuestions[questions.QuestionGroupID].Questions = append(checklistQuestions[questions.QuestionGroupID].Questions, &questions)
	}

	for _, v := range checklistQuestions {
		response = append(response, v)
	}

	// Ensure response is not null
	if response == nil {
		response = []*masterDTOs.ChecklistQuestionsGroupByQuestionGroupDTO{}
	}

	// Calculate total pages for pagination
	totalPages := int((totalCount + int64(pageSize) - 1) / int64(pageSize))

	// Create pagination info
	paginationInfo := commonDTOs.PaginationInfoDTO{
		CurrentPage: page,
		PageSize:    pageSize,
		TotalPages:  totalPages,
		TotalCount:  totalCount,
	}

	// Wrap response and pagination info in a single DTO
	return &masterDTOs.ChecklistQuestionsGroupByQuestionGroupListResponseDTO{
		Items:      response,
		Pagination: paginationInfo,
	}, nil

}

func (repo *ChecklistQuestionRepository) GetByID(id uint) (*models.ChecklistQuestion, error) {
	var checklistQuestion *models.ChecklistQuestion
	result := repo.db.Preload("Checklist").Preload("QuestionGroup").First(&checklistQuestion, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return checklistQuestion, result.Error
}

func (repo *ChecklistQuestionRepository) Create(checklistQuestion *models.ChecklistQuestion) (uint, error) {

	if err := repo.db.Create(checklistQuestion).Error; err != nil {
		return 0, err
	}

	return checklistQuestion.ID, nil

}

func (repo *ChecklistQuestionRepository) Update(id uint, checklistQuestion *models.ChecklistQuestion) error {

	if err := repo.db.Save(checklistQuestion).Error; err != nil {
		return err
	}

	return nil

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

func (repo *ChecklistQuestionRepository) ExistsByTitle(QuestionEn, QuestionAr string, ChecklistID uint, QuestionGroupID uint) (bool, error) {
	var count int64

	// Check for the existence of QuestionEn or QuestionAr in a single query
	err := repo.db.Model(&models.ChecklistQuestion{}).
		Where("(question_en = ? OR question_ar = ?) AND checklist_id = ? AND question_group_id = ?", QuestionEn, QuestionAr, ChecklistID, QuestionGroupID).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)

	}
	// Return true if either title exists, otherwise false
	return count > 0, nil
}

func (repo *ChecklistQuestionRepository) ExistsQuestionLibrary(QuestionLibraryID uint, QuestionEn, QuestionAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.ChecklistQuestion{}).
		Where("question_library_id = ? AND (question_en = ? OR question_ar = ?)", QuestionLibraryID, QuestionEn, QuestionAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)

	}
	return count > 0, nil
}

func (repo *ChecklistQuestionRepository) ExistsByTitleExcludingID(id uint, QuestionEn, QuestionAr string, ChecklistID uint, QuestionGroupID uint) (bool, error) {
	var count int64

	err := repo.db.Model(&models.ChecklistQuestion{}).
		Where("id <> ? AND (question_en = ? OR question_ar = ?) AND checklist_id = ? AND question_group_id = ?", id, QuestionEn, QuestionAr, ChecklistID, QuestionGroupID).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}

func (repo *ChecklistQuestionRepository) ExistsQuestionLibraryExcludingID(id uint, QuestionEn, QuestionAr string, QuestionLibraryID uint) (bool, error) {
	var count int64

	err := repo.db.Model(&models.ChecklistQuestion{}).
		Where("id <> ? AND (question_en = ? OR question_ar = ?) AND question_library_id = ?", id, QuestionEn, QuestionAr, QuestionLibraryID).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)

	}
	return count > 0, nil
}
