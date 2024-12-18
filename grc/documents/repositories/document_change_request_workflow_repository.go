package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type DocumentChangeRequestWorkflowRepositoryInterface interface {
	CreateDocumentCRWorkflow(documentCRWorkflow *models.DocumentChangeRequestWorkflow, tx ...utility.DatabaseTransaction) (uint, error)
	CreateDocumentCRWorkflowParticipant(docWorkflowCRParticipant *models.DocumentChangeRequestWorkflowParticipant, tx ...utility.DatabaseTransaction) error

	DeleteDocumentCRWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteDocumentCRWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteDocumentCRWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error
}

type DocumentChangeRequestWorkflowRepository struct {
	db *gorm.DB
}

func NewDocumentChangeRequestWorkflowRepository(db *gorm.DB) *DocumentChangeRequestWorkflowRepository {
	return &DocumentChangeRequestWorkflowRepository{
		db: db,
	}
}

func (repo *DocumentChangeRequestWorkflowRepository) CreateDocumentCRWorkflow(documentCRWorkflow *models.DocumentChangeRequestWorkflow, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(documentCRWorkflow).Error; err != nil {
		return 0, err
	}

	return documentCRWorkflow.ID, nil
}

func (repo *DocumentChangeRequestWorkflowRepository) CreateDocumentCRWorkflowParticipant(docWorkflowCRParticipant *models.DocumentChangeRequestWorkflowParticipant, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(docWorkflowCRParticipant).Error; err != nil {
		return err
	}

	return nil
}

func (repo *DocumentChangeRequestWorkflowRepository) DeleteDocumentCRWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", participantIds).Delete(&models.DocumentChangeRequestWorkflowParticipant{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *DocumentChangeRequestWorkflowRepository) DeleteDocumentCRWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", reviewIds).Delete(&models.DocumentChangeRequestWorkflowReview{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *DocumentChangeRequestWorkflowRepository) DeleteDocumentCRWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", workflowIds).Delete(&models.DocumentChangeRequestWorkflow{}).Error; err != nil {
		return err
	}

	return nil
}
