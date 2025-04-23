package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type DocumentWorkflowRepositoryInterface interface {
	GetWorkflowList(docVersionId uint) (*models.DocumentVersion, error)
	CreateDocumentWorkflow(documentWorkflow *models.DocumentWorkflow, tx ...utility.DatabaseTransaction) (uint, error)
	UpdateDocumentWorkflow(documentWorkflow *models.DocumentWorkflow, tx ...utility.DatabaseTransaction) error
	CreateDocumentWorkflowParticipant(docWorkflowParticipant *models.DocumentWorkflowParticipant, tx ...utility.DatabaseTransaction) error
	DeleteDocumentWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteDocumentWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteDocumentWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error

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

func (repo *DocumentWorkflowRepository) GetWorkflowList(docVersionId uint) (*models.DocumentVersion, error) {
	var documentVersion models.DocumentVersion
	result := repo.db.
		Preload("SubmittedByUser").
		Preload("DocumentWorkflows").
		Preload("DocumentWorkflows.DocumentWorkflowStatus").
		Preload("DocumentWorkflows.DocumentWorkflowParticipants").
		Preload("DocumentWorkflows.DocumentWorkflowParticipants.User").
		Preload("DocumentWorkflows.DocumentWorkflowParticipants.Role").
		Preload("DocumentWorkflows.DocumentWorkflowReviews", "document_workflow_status_id IN (?)", repo.db.Model(&models.DocumentWorkflowStatus{}).Select("id").Where("type = ?", "approve")).
		Preload("DocumentWorkflows.DocumentWorkflowReviews.CreatedByUser").
		Where("id = ?", docVersionId).
		Find(&documentVersion)

	if result.Error != nil {
		return nil, result.Error
	}

	return &documentVersion, nil
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
		Preload("DocumentWorkflowReviews.DocumentWorkflowStatus").
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

// func (repo *DocumentWorkflowRepository) CountApprovedParticipants(workflowID uint, tx ...utility.DatabaseTransaction) (int, error) {
// 	db := repo.db
// 	if len(tx) > 0 && tx[0] != nil {
// 		db = tx[0].DB()
// 	}
// 	var count int64
// 	err := db.Model(&models.DocumentWorkflowReview{}).
// 		Where("document_workflow_id = ? AND document_workflow_status_id = (SELECT id FROM document_workflow_statuses WHERE type = 'approve')", workflowID).
// 		Distinct("created_by_id").
// 		Count(&count).Error

// 	if err != nil {
// 		return 0, err
// 	}

//		return int(count), nil
//	}

func (repo *DocumentWorkflowRepository) CountApprovedParticipants(workflowID uint, tx ...utility.DatabaseTransaction) (int, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var resubmitID *uint // Pointer to handle NULL values
	err := db.Raw(`
        SELECT MAX(id)
        FROM document_workflow_reviews
        WHERE document_workflow_id = ?
          AND document_workflow_status_id = (
              SELECT id
              FROM document_workflow_statuses
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
	query := db.Model(&models.DocumentWorkflowReview{}).
		Where("document_workflow_id = ? AND document_workflow_status_id = (SELECT id FROM document_workflow_statuses WHERE type = 'approve')", workflowID).
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

func (repo *DocumentWorkflowRepository) DeleteDocumentWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", participantIds).Delete(&models.DocumentWorkflowParticipant{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *DocumentWorkflowRepository) DeleteDocumentWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", reviewIds).Delete(&models.DocumentWorkflowReview{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *DocumentWorkflowRepository) DeleteDocumentWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", workflowIds).Delete(&models.DocumentWorkflow{}).Error; err != nil {
		return err
	}

	return nil
}
