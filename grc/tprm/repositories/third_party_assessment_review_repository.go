package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ThirdPartyAssessmentReviewRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	CreateTprmReview(tprmReview *models.ThirdPartyAssessmentWorkflowReview, tx ...utility.DatabaseTransaction) error
	UpdateAssessment(assessment *models.ThirdPartyAssessment, tx ...utility.DatabaseTransaction) error
	UpdateThirdPartyRegister(thirdPartyRegiste *models.ThirdPartyRegister, tx ...utility.DatabaseTransaction) error
	GetAssessmentDetails(id uint, tx ...utility.DatabaseTransaction) (*models.ThirdPartyAssessment, error)
	DeleteTprmAssessmentWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error
}

type ThirdPartyAssessmentReviewRepository struct {
	db *gorm.DB
}

func NewThirdPartyAssessmentReviewRepository(db *gorm.DB) ThirdPartyAssessmentReviewRepositoryInterface {
	return &ThirdPartyAssessmentReviewRepository{db: db}
}

func (repo *ThirdPartyAssessmentReviewRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *ThirdPartyAssessmentReviewRepository) CreateTprmReview(tprmReview *models.ThirdPartyAssessmentWorkflowReview, tx ...utility.DatabaseTransaction) error {
	db := repo.db

	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(tprmReview).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ThirdPartyAssessmentReviewRepository) UpdateAssessment(assessment *models.ThirdPartyAssessment, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Save(assessment).Error
}

func (repo *ThirdPartyAssessmentReviewRepository) UpdateThirdPartyRegister(thirdPartyRegister *models.ThirdPartyRegister, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Save(thirdPartyRegister).Error
}

func (repo *ThirdPartyAssessmentReviewRepository) GetAssessmentDetails(id uint, tx ...utility.DatabaseTransaction) (*models.ThirdPartyAssessment, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var entity *models.ThirdPartyAssessment
	if err := db.Preload("ThirdPartyRegister").
		Preload("ThirdPartyAssessmentWorkflows").
		Preload("ThirdPartyAssessmentWorkflows.ThirdPartyAssessmentWorkflowParticipants").
		Preload("ThirdPartyAssessmentWorkflows.ThirdPartyAssessmentWorkflowReviews").
		First(&entity, id).Error; err != nil {
		return nil, err
	}

	return entity, nil
}

func (repo *ThirdPartyAssessmentReviewRepository) DeleteTprmAssessmentWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", reviewIds).Delete(&models.ThirdPartyAssessmentWorkflowReview{}).Error; err != nil {
		return err
	}

	return nil
}
