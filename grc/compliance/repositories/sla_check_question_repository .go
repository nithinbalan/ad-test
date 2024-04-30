package repositories

import (
	"errors"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/models"

	"gorm.io/gorm"
)

type SLACheckQuestionRepositoryInterface interface {
	GetAllQuestionGroups(checkID uint) ([]*complianceDTOs.SLACheckQuestionGroupListDTO, error)
	GetQuestionGroupByID(groupID uint) (*models.SLACheckQuestionGroup, error)
	CreateQuestionGroup(slaCheckQuestionGroup *models.SLACheckQuestionGroup) (uint, error)
	UpdateQuestionGroup(group *models.SLACheckQuestionGroup) error
	DeleteQuestionGroup(id uint) error
	GetQuestionsByGroupID(groupID uint) ([]models.SLACheckQuestion, error)
	GetQuestionByID(id uint) (*models.SLACheckQuestion, error)
	CreateQuestion(slaCheckQuestion *models.SLACheckQuestion) (uint, error)
	UpdateQuestion(SLACheckQuestion *models.SLACheckQuestion) error
	DeleteQuestion(id uint) error
	GetSimilarQuestion(checkID uint, titleEn string, titleAr string) (*models.SLACheckQuestion, error)
	GetSimilarQuestionGroup(titleEn string, titleAr string, checkID uint) (*models.SLACheckQuestionGroup, error)
	IncrementTotalCount(groupID uint) error
	IncrementCompletedCount(groupID uint) error
}

type SLACheckQuestionRepository struct {
	db *gorm.DB
}

func NewSLACheckQuestionRepository(db *gorm.DB) *SLACheckQuestionRepository {
	return &SLACheckQuestionRepository{
		db: db,
	}
}

func (repo *SLACheckQuestionRepository) GetAllQuestionGroups(checkID uint) ([]*complianceDTOs.SLACheckQuestionGroupListDTO, error) {
	var slaCheckQuestionGroups []*complianceDTOs.SLACheckQuestionGroupListDTO

	if err := repo.db.Table("sla_check_question_groups").Where("sla_check_id = ?", checkID).Find(&slaCheckQuestionGroups).Error; err != nil {
		return nil, err
	}

	return slaCheckQuestionGroups, nil
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

func (repo *SLACheckQuestionRepository) GetSimilarQuestion(checkID uint, titleEn string, titleAr string) (*models.SLACheckQuestion, error) {
	var slaCheckQuestion models.SLACheckQuestion
	if err := repo.db.
		Where("sla_check_question_group_id IN (SELECT id FROM sla_check_question_groups WHERE sla_check_id = ?) AND (title_en = ? OR title_ar = ?)", checkID, titleEn, titleAr).
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
	if err := repo.db.Where("sla_check_question_group_id = ?", groupID).Preload("Answers").Find(&questions).Error; err != nil {
		return nil, err
	}
	return questions, nil
}

// IncrementTotalCount increments the total count of SLACheckQuestionGroup by 1
func (repo *SLACheckQuestionRepository) IncrementTotalCount(groupID uint) error {
	var group models.SLACheckQuestionGroup
	if err := repo.db.First(&group, groupID).Error; err != nil {
		return err
	}

	group.TotalCount++

	if err := repo.db.Save(&group).Error; err != nil {
		return err
	}

	return nil
}

// IncrementCompletedCount increments the total count of SLACheckQuestionGroup by 1
func (repo *SLACheckQuestionRepository) IncrementCompletedCount(groupID uint) error {
	var group models.SLACheckQuestionGroup
	if err := repo.db.First(&group, groupID).Error; err != nil {
		return err
	}

	group.CompletedCount++

	if err := repo.db.Save(&group).Error; err != nil {
		return err
	}

	return nil
}
