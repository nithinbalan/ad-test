package repositories

import (
	// complianceDTOs "jethur/grc/compliance/dtos"
	"errors"
	"jethur/grc/compliance/models"
	masterModels "jethur/grc/masters/models"

	"gorm.io/gorm"
)

type SLACheckQuestionRepositoryInterface interface {
	GetQuestionByID(id uint) (*models.SLACheckQuestion, error)
	CreateQuestion(slaCheckQuestion *models.SLACheckQuestion) (uint, error)
	UpdateQuestion(SLACheckQuestion *models.SLACheckQuestion) error
	DeleteQuestion(id uint) error
	GetQuestionGroupByID(groupID uint) (*models.SLACheckQuestionGroup, error)
	CreateQuestionGroup(slaCheckQuestionGroup *models.SLACheckQuestionGroup) (uint, error)
	UpdateQuestionGroup(group *models.SLACheckQuestionGroup) error
	DeleteQuestionGroup(id uint) error
	GetSimilarQuestion(checkID uint, question *masterModels.ChecklistQuestion) (*models.SLACheckQuestion, error)
	GetSimilarQuestionGroup(titleEn string, titleAr string, checkID uint) (*models.SLACheckQuestionGroup, error)
	GetQuestionsByGroupID(groupID uint) ([]models.SLACheckQuestion, error)
}

type SLACheckQuestionRepository struct {
	db *gorm.DB
}

func NewSLACheckQuestionRepository(db *gorm.DB) *SLACheckQuestionRepository {
	return &SLACheckQuestionRepository{
		db: db,
	}
}

func (repo *SLACheckQuestionRepository) GetQuestionByID(id uint) (*models.SLACheckQuestion, error) {
	var SLACheckQuestion *models.SLACheckQuestion

	result := repo.db.First(&SLACheckQuestion, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return SLACheckQuestion, nil

}

func (repo *SLACheckQuestionRepository) CreateQuestion(slaCheckQuestion *models.SLACheckQuestion) (uint, error) {

	if err := repo.db.Create(slaCheckQuestion).Error; err != nil {
		return 0, err
	}

	return slaCheckQuestion.ID, nil
}

func (repo *SLACheckQuestionRepository) UpdateQuestion(SLACheckQuestion *models.SLACheckQuestion) error {
	if err := repo.db.Save(SLACheckQuestion).Error; err != nil {
		return err
	}
	return nil
}

func (repo *SLACheckQuestionRepository) DeleteQuestion(id uint) error {

	var slaCheckQuestion *models.SLACheckQuestion
	if err := repo.db.First(&slaCheckQuestion, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.SLACheckQuestion{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *SLACheckQuestionRepository) GetQuestionGroupByID(groupID uint) (*models.SLACheckQuestionGroup, error) {
	var group models.SLACheckQuestionGroup
	if err := repo.db.First(&group, groupID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &group, nil
}

func (repo *SLACheckQuestionRepository) CreateQuestionGroup(slaCheckQuestionGroup *models.SLACheckQuestionGroup) (uint, error) {

	if err := repo.db.Create(slaCheckQuestionGroup).Error; err != nil {
		return 0, err
	}

	return slaCheckQuestionGroup.ID, nil
}

func (repo *SLACheckQuestionRepository) UpdateQuestionGroup(group *models.SLACheckQuestionGroup) error {
	if err := repo.db.Save(group).Error; err != nil {
		return err
	}
	return nil
}

func (repo *SLACheckQuestionRepository) DeleteQuestionGroup(id uint) error {

	var slaCheckQuestionGroup *models.SLACheckQuestionGroup
	if err := repo.db.First(&slaCheckQuestionGroup, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.SLACheckQuestionGroup{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *SLACheckQuestionRepository) GetSimilarQuestion(checkID uint, question *masterModels.ChecklistQuestion) (*models.SLACheckQuestion, error) {
	var slaCheckQuestion models.SLACheckQuestion
	if err := repo.db.
		Where("sla_check_question_group_id IN (SELECT id FROM sla_check_question_groups WHERE sla_check_id = ?) AND (title_en = ? OR title_ar = ?)", checkID, question.QuestionEn, question.QuestionAr).
		First(&slaCheckQuestion).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &slaCheckQuestion, nil
}

func (repo *SLACheckQuestionRepository) GetSimilarQuestionGroup(titleEn string, titleAr string, checkID uint) (*models.SLACheckQuestionGroup, error) {
	var group models.SLACheckQuestionGroup
	if err := repo.db.Where("(title_en = ? OR title_ar = ?) AND sla_check_id = ?", titleEn, titleAr, checkID).First(&group).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &group, nil
}

func (repo *SLACheckQuestionRepository) GetQuestionsByGroupID(groupID uint) ([]models.SLACheckQuestion, error) {
	var questions []models.SLACheckQuestion
	if err := repo.db.Where("sla_check_question_group_id = ?", groupID).Find(&questions).Error; err != nil {
		return nil, err
	}
	return questions, nil
}
