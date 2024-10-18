package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type DocumentReviewRepositoryInterface interface {
	GetComments(docVersionId uint, currentApprovalLevel int) (*models.DocumentVersion, error)
	CreateDocumentWorkflowReview(documentWorkflowReview *models.DocumentWorkflowReview, tx ...utility.DatabaseTransaction) error
}

type DocumentReviewRepository struct {
	db *gorm.DB
}

func NewDocumentReviewRepository(db *gorm.DB) *DocumentReviewRepository {
	return &DocumentReviewRepository{db: db}
}

func (repo *DocumentReviewRepository) GetComments(docVersionId uint, currentApprovalLevel int) (*models.DocumentVersion, error) {
	var documentVersion models.DocumentVersion
	result := repo.db.
		Preload("SubmittedByUser").
		Preload("DocumentWorkflows", "level <= ?", currentApprovalLevel).
		Preload("DocumentWorkflows.DocumentWorkflowStatus").
		Preload("DocumentWorkflows.DocumentWorkflowParticipants").
		Preload("DocumentWorkflows.DocumentWorkflowParticipants.User").
		Preload("DocumentWorkflows.DocumentWorkflowReviews").
		Preload("DocumentWorkflows.DocumentWorkflowReviews.DocumentWorkflowStatus").
		Preload("DocumentWorkflows.DocumentWorkflowReviews.CreatedByUser").
		Preload("DocumentWorkflows.DocumentWorkflowReviews.UpdatedByUser").
		Where("id = ?", docVersionId).
		Find(&documentVersion)

	if result.Error != nil {
		return nil, result.Error
	}

	return &documentVersion, nil
}

func (repo *DocumentReviewRepository) CreateDocumentWorkflowReview(documentWorkflowReview *models.DocumentWorkflowReview, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(documentWorkflowReview).Error; err != nil {
		return err
	}
	return nil
}
