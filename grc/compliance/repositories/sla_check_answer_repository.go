package repositories

import (
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type SLACheckAnswerRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetByID(id uint) (*models.SLACheckAnswer, error)
	Create(slaCheckAnswer *models.SLACheckAnswer) (*models.SLACheckAnswer, error)
	Update(slaCheckAnswer *models.SLACheckAnswer) error
	Delete(id uint) error
	AssessmentSummary(id uint) (*complianceDTOs.AssessmentSummaryDTO, error)
}

type SLACheckAnswerRepository struct {
	db *gorm.DB
}

func NewSLACheckAnswerRepository(db *gorm.DB) *SLACheckAnswerRepository {
	return &SLACheckAnswerRepository{
		db: db,
	}
}

func (repo *SLACheckAnswerRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *SLACheckAnswerRepository) GetByID(id uint) (*models.SLACheckAnswer, error) {
	var slaCheckAnswer *models.SLACheckAnswer
	result := repo.db.First(&slaCheckAnswer, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return slaCheckAnswer, nil
}

func (repo *SLACheckAnswerRepository) Create(slaCheckAnswer *models.SLACheckAnswer) (*models.SLACheckAnswer, error) {

	if err := repo.db.Create(slaCheckAnswer).Error; err != nil {
		return nil, err
	}

	return slaCheckAnswer, nil
}

func (repo *SLACheckAnswerRepository) Update(slaCheckAnswer *models.SLACheckAnswer) error {

	if err := repo.db.Save(slaCheckAnswer).Error; err != nil {
		return err
	}

	return nil
}

func (repo *SLACheckAnswerRepository) Delete(id uint) error {

	var slaCheckAnswer *models.SLACheckAnswer
	if err := repo.db.First(&slaCheckAnswer, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.SLACheckAnswer{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *SLACheckAnswerRepository) AssessmentSummary(id uint) (*complianceDTOs.AssessmentSummaryDTO, error) {
	var questions []*models.SLACheckQuestion
	var answers []*models.SLACheckAnswer

	if err := repo.db.Where("sla_check_question_group_id IN (SELECT id FROM sla_check_question_groups WHERE sla_check_id IN (SELECT id FROM sla_checks WHERE sla_register_id = ?))", id).Find(&questions).Error; err != nil {
		return nil, err
	}

	// Collect the IDs of the questions
	var questionIDs []uint
	for _, question := range questions {
		questionIDs = append(questionIDs, question.ID)
	}
	if err := repo.db.Where("sla_check_question_id IN (?)", questionIDs).Find(&answers).Error; err != nil {
		return nil, err
	}

	// Calculate statistics
	totalQuestions := len(questions)
	totalAnswers := len(answers)
	answerCounts := make(map[uint]int)

	for _, answer := range answers {
		answerCounts[answer.ComplianceStatusID]++
	}

	statistics := &complianceDTOs.AssessmentSummaryDTO{
		TotalQuestions: totalQuestions,
		TotalAnswers:   totalAnswers,
		AnswerCounts:   answerCounts,
	}

	return statistics, nil
}
