package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type DocumentWorkflowRepositoryInterface interface {
	CreateDocumentWorkflow(documentWorkflow *models.DocumentWorkflow, tx ...utility.DatabaseTransaction) (uint, error)
	UpdateDocumentWorkflow(documentWorkflow *models.DocumentWorkflow, tx ...utility.DatabaseTransaction) error
	CreateDocumentWorkflowParticipant(docWorkflowParticipant *models.DocumentWorkflowParticipant, tx ...utility.DatabaseTransaction) error

	GetDocumentWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.DocumentWorkflow, error)
	GetDocumentWorkflowByLevel(docVersionId uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.DocumentWorkflow, error)

	GetNextLevel(documentVersionID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.DocumentWorkflow, error)
	CountApprovedParticipants(workflowID uint, tx ...utility.DatabaseTransaction) (int, error)
}

type DocumentWorkflowRepository struct {
	db *gorm.DB
}

func NewDocumentWorkflowRepository(db *gorm.DB) *DocumentWorkflowRepository {
	return &DocumentWorkflowRepository{db: db}
}

func (repo *DocumentWorkflowRepository) CreateDocumentWorkflow(documentWorkflow *models.DocumentWorkflow, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(documentWorkflow).Error; err != nil {
		return 0, err
	}

	return documentWorkflow.ID, nil
}

func (repo *DocumentWorkflowRepository) UpdateDocumentWorkflow(documentWorkflow *models.DocumentWorkflow, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(documentWorkflow).Error; err != nil {
		return err
	}

	return nil
}

func (repo *DocumentWorkflowRepository) CreateDocumentWorkflowParticipant(docWorkflowParticipant *models.DocumentWorkflowParticipant, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(docWorkflowParticipant).Error; err != nil {
		return err
	}

	return nil
}

func (repo *DocumentWorkflowRepository) GetDocumentWorkflowByLevel(docVersionId uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.DocumentWorkflow, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var documentWorkflow *models.DocumentWorkflow
	result := db.
		Preload("WorkflowLevelType").
		Preload("DocumentWorkflowParticipants").
		Preload("DocumentWorkflowReviews").
		Preload("DocumentWorkflowParticipants.User").
		Where("document_version_id = ? AND level = ?", docVersionId, currentApprovalLevel).
		Find(&documentWorkflow)

	if result.Error != nil {
		return nil, result.Error
	}

	return documentWorkflow, nil
}

func (repo *DocumentWorkflowRepository) GetDocumentWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.DocumentWorkflow, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var documentWorkflow *models.DocumentWorkflow

	result := db.First(&documentWorkflow, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return documentWorkflow, nil
}

func (repo *DocumentWorkflowRepository) GetNextLevel(documentVersionID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.DocumentWorkflow, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var nextWorkflow models.DocumentWorkflow
	result := db.Where("document_version_id = ? AND level > ?", documentVersionID, currentLevel).
		Order("level ASC").  // Ensure we are ascending by level
		First(&nextWorkflow) // Retrieve the first result only which will be the next level

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil // No further levels available
	} else if result.Error != nil {
		return nil, result.Error // Other error occurred
	}

	return &nextWorkflow, nil
}

func (repo *DocumentWorkflowRepository) CountApprovedParticipants(workflowID uint, tx ...utility.DatabaseTransaction) (int, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var count int64
	err := db.Model(&models.DocumentWorkflowReview{}).
		Where("document_workflow_id = ? AND document_workflow_status_id = (SELECT id FROM document_workflow_statuses WHERE type = 'approve')", workflowID).
		Distinct("created_by_id").
		Count(&count).Error

	if err != nil {
		return 0, err
	}

	return int(count), nil
}
