package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type DocumentChangeRequestWorkflowRepositoryInterface interface {
	CreateDocumentCRWorkflow(documentCRWorkflow *models.DocumentChangeRequestWorkflow, tx ...utility.DatabaseTransaction) (uint, error)
	UpdateDocumentCRWorkflow(documentCRWorkflow *models.DocumentChangeRequestWorkflow, tx ...utility.DatabaseTransaction) error
	CreateDocumentCRWorkflowParticipant(docWorkflowCRParticipant *models.DocumentChangeRequestWorkflowParticipant, tx ...utility.DatabaseTransaction) error

	GetDocumentCRWorkflowByLevel(changeRequestID uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.DocumentChangeRequestWorkflow, error)
	GetDocumentCRWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.DocumentChangeRequestWorkflow, error)
	GetNextLevel(id uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.DocumentChangeRequestWorkflow, error)
	CountApprovedParticipants(workflowID uint, tx ...utility.DatabaseTransaction) (int, error)

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

func (repo *DocumentChangeRequestWorkflowRepository) UpdateDocumentCRWorkflow(documentCRWorkflow *models.DocumentChangeRequestWorkflow, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(documentCRWorkflow).Error; err != nil {
		return err
	}

	return nil
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

func (repo *DocumentChangeRequestWorkflowRepository) GetDocumentCRWorkflowByLevel(changeRequestID uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.DocumentChangeRequestWorkflow, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var documentCRWorkflow *models.DocumentChangeRequestWorkflow
	result := db.
		Preload("WorkflowLevelType").
		Preload("DocumentChangeRequestWorkflowParticipants").
		Preload("DocumentChangeRequestWorkflowReviews").
		Preload("DocumentChangeRequestWorkflowParticipants.User").
		Where("document_change_request_id = ? AND level = ?", changeRequestID, currentApprovalLevel).
		Find(&documentCRWorkflow)

	if result.Error != nil {
		return nil, result.Error
	}

	return documentCRWorkflow, nil
}

func (repo *DocumentChangeRequestWorkflowRepository) GetDocumentCRWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.DocumentChangeRequestWorkflow, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var documentCRWorkflow *models.DocumentChangeRequestWorkflow

	result := db.First(&documentCRWorkflow, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return documentCRWorkflow, nil
}

func (repo *DocumentChangeRequestWorkflowRepository) GetNextLevel(id uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.DocumentChangeRequestWorkflow, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var nextWorkflow models.DocumentChangeRequestWorkflow
	result := db.Where("document_change_request_id = ? AND level > ?", id, currentLevel).
		Order("level ASC").  // Ensure we are ascending by level
		First(&nextWorkflow) // Retrieve the first result only which will be the next level

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil // No further levels available
	} else if result.Error != nil {
		return nil, result.Error // Other error occurred
	}

	return &nextWorkflow, nil
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

func (repo *DocumentChangeRequestWorkflowRepository) CountApprovedParticipants(workflowID uint, tx ...utility.DatabaseTransaction) (int, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var resubmitID *uint // Pointer to handle NULL values
	err := db.Raw(`
        SELECT MAX(id)
        FROM document_change_request_workflow_reviews
        WHERE document_change_request_workflow_id = ?
          AND document_change_request_workflow_status_id = (
              SELECT id
              FROM document_change_request_workflow_statuses
              WHERE type = 'resubmit'
          )
    `, workflowID).Scan(&resubmitID).Error

	if err != nil {
		return 0, err
	}

	if resubmitID == nil {
		resubmitID = new(uint) // Initialize to zero if needed
	}

	var count int64
	query := db.Model(&models.DocumentChangeRequestWorkflowReview{}).
		Where("document_change_request_workflow_id = ? AND document_change_request_workflow_status_id = (SELECT id FROM document_change_request_workflow_statuses WHERE type = 'approve')", workflowID).
		Distinct("created_by_id")

	if *resubmitID > 0 {
		query = query.Where("id > ?", *resubmitID)
	}

	err = query.Count(&count).Error
	if err != nil {
		return 0, err
	}

	return int(count), nil
}
