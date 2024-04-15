package repositories

import (
	// complianceDTOs "jethur/grc/compliance/dtos"
	"errors"
	"jethur/grc/compliance/models"
	masterModels "jethur/grc/masters/models"

	"gorm.io/gorm"
)

type SLACheckQuestionRepositoryInterface interface {
	// GetAll(slaCheckQuestionID uint, sort string, sortBy string, keyword string) ([]*complianceDTOs.SLACheckQuestionListDTO, error)
	// GetByID(id uint) (*models.SLACheckQuestion, error)
	CreateQuestionGroup(slaCheckQuestionGroup *models.SLACheckQuestionGroup) (uint, error)
	CreateQuestion(slaCheckQuestion *models.SLACheckQuestion) (uint, error)
	Update(SLACheckQuestion *models.SLACheckQuestion) error
	GetQuestionGroupByTitleAndCheckID(title string, checkID uint) (*models.SLACheckQuestionGroup, error)
	GetSimilarQuestion(checkID uint, question *masterModels.ChecklistQuestion) (*models.SLACheckQuestion, error)
}

type SLACheckQuestionRepository struct {
	db *gorm.DB
}

func NewSLACheckQuestionRepository(db *gorm.DB) *SLACheckQuestionRepository {
	return &SLACheckQuestionRepository{
		db: db,
	}
}

// func (repo *SLACheckQuestionRepository) GetAll(slaCheckQuestionID uint, sort string, sortBy string, keyword string) ([]*complianceDTOs.SLACheckQuestionListDTO, error) {

// 	var slaCheckQuestions []*complianceDTOs.SLACheckQuestionListDTO

// 	query := repo.db.Table("sla_check_questions").
// 		Select(`
//             sla_check_questions.*,
//             compliance_requirements.title AS compliance_requirement_title,
//             compliance_statuses.title_en AS compliance_status_title_en,
//             compliance_statuses.title_ar AS compliance_status_title_ar,
//             compliance_statuses.color_code AS compliance_status_color_code,
//             created_user.first_name_en AS created_user_first_name_en,
//             created_user.first_name_ar AS created_user_first_name_ar,
//             created_user.last_name_en AS created_user_last_name_en,
//             created_user.last_name_ar AS created_user_last_name_ar,
//             updated_user.first_name_en AS updated_user_first_name_en,
//             updated_user.first_name_ar AS updated_user_first_name_ar,
//             updated_user.last_name_en AS updated_user_last_name_en,
//             updated_user.last_name_ar AS updated_user_last_name_ar
//         `).
// 		Joins(`
//             LEFT JOIN compliance_requirements ON compliance_requirements.id = sla_check_questions.compliance_requirement_id
//             LEFT JOIN compliance_statuses ON compliance_statuses.id = sla_check_questions.compliance_status_id
//             LEFT JOIN users AS created_user ON created_user.id = sla_check_questions.created_by
//             LEFT JOIN users AS updated_user ON updated_user.id = sla_check_questions.updated_by
//         `).
// 		Where("sla_check_questions.compliance_requirement_id = ?", slaCheckQuestionID)

// 	if keyword != "" {
// 		query = query.Where(`
//             (compliance_requirements.title LIKE ? OR
//             sla_check_questions.comment LIKE ? OR
//             compliance_statuses.title_en LIKE ? OR
//             compliance_statuses.title_ar LIKE ?)`,
// 			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
// 	}
// 	if sortBy == "" {
// 		sortBy = "sla_check_questions.id"
// 	}
// 	if sort != "asc" && sort != "desc" {
// 		sort = "asc"
// 	}
// 	query = query.Order(sortBy + " " + sort)
// 	result := query.Find(&slaCheckQuestions)

// 	return slaCheckQuestions, result.Error
// }

// func (repo *SLACheckQuestionRepository) GetByID(id uint) (*models.SLACheckQuestion, error) {
// 	var slaCheckQuestion *models.SLACheckQuestion
// 	result := repo.db.First(&slaCheckQuestion, id)

// 	if result.Error != nil {
// 		return nil, result.Error
// 	}

// 	return slaCheckQuestion, nil
// }

func (repo *SLACheckQuestionRepository) CreateQuestionGroup(slaCheckQuestionGroup *models.SLACheckQuestionGroup) (uint, error) {

	if err := repo.db.Create(slaCheckQuestionGroup).Error; err != nil {
		return 0, err
	}

	return slaCheckQuestionGroup.ID, nil
}

func (repo *SLACheckQuestionRepository) CreateQuestion(slaCheckQuestion *models.SLACheckQuestion) (uint, error) {

	if err := repo.db.Create(slaCheckQuestion).Error; err != nil {
		return 0, err
	}

	return slaCheckQuestion.ID, nil
}

func (repo *SLACheckQuestionRepository) Update(SLACheckQuestion *models.SLACheckQuestion) error {
	if err := repo.db.Save(SLACheckQuestion).Error; err != nil {
		return err
	}
	return nil
}

func (repo *SLACheckQuestionRepository) GetQuestionGroupByTitleAndCheckID(title string, checkID uint) (*models.SLACheckQuestionGroup, error) {
	var group models.SLACheckQuestionGroup
	if err := repo.db.Where("title_en = ? AND sla_check_id = ?", title, checkID).First(&group).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Group not found
		}
		return nil, err
	}
	return &group, nil
}

func (repo *SLACheckQuestionRepository) GetSimilarQuestion(checkID uint, question *masterModels.ChecklistQuestion) (*models.SLACheckQuestion, error) {
	var slaCheckQuestion models.SLACheckQuestion
	if err := repo.db.
		Where("sla_check_question_group_id IN (SELECT id FROM sla_check_question_groups WHERE sla_check_id = ?) AND (title_en = ? OR title_ar = ?)", checkID, question.QuestionEn, question.QuestionAr).
		First(&slaCheckQuestion).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No similar question found
		}
		return nil, err
	}
	return &slaCheckQuestion, nil
}
