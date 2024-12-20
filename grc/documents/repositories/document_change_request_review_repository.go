package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type DocumentChangeRequestReviewRepositoryInterface interface {
	GetComments(id uint, currentApprovalLevel int) (*models.DocumentChangeRequest, error)
	CreateDocumentWorkflowReview(documentWorkflowReview *models.DocumentWorkflowReview, tx ...utility.DatabaseTransaction) error
}

type DocumentChangeRequestReviewRepository struct {
	db *gorm.DB
}

func NewDocumentChangeRequestReviewRepository(db *gorm.DB) *DocumentChangeRequestReviewRepository {
	return &DocumentChangeRequestReviewRepository{db: db}
}

func (repo *DocumentChangeRequestReviewRepository) GetComments(id uint, currentApprovalLevel int) (*models.DocumentChangeRequest, error) {
	var documentChangeRequest models.DocumentChangeRequest
	result := repo.db.
		Preload("SubmittedByUser").
		Preload("DocumentChangeRequestWorkflows", "level <= ?", currentApprovalLevel).
		Preload("DocumentChangeRequestWorkflows.DocumentChangeRequestWorkflowStatus").
		Preload("DocumentChangeRequestWorkflows.DocumentChangeRequestWorkflowParticipants").
		Preload("DocumentChangeRequestWorkflows.DocumentChangeRequestWorkflowParticipants.User").
		Preload("DocumentChangeRequestWorkflows.DocumentChangeRequestWorkflowReviews").
		Preload("DocumentChangeRequestWorkflows.DocumentChangeRequestWorkflowReviews.DocumentChangeRequestWorkflowStatus").
		Preload("DocumentChangeRequestWorkflows.DocumentChangeRequestWorkflowReviews.CreatedByUser").
		Preload("DocumentChangeRequestWorkflows.DocumentChangeRequestWorkflowReviews.UpdatedByUser").
		Where("id = ?", id).
		Find(&documentChangeRequest)

	if result.Error != nil {
		return nil, result.Error
	}

	return &documentChangeRequest, nil
}

func (repo *DocumentChangeRequestReviewRepository) CreateDocumentWorkflowReview(documentWorkflowReview *models.DocumentWorkflowReview, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(documentWorkflowReview).Error; err != nil {
		return err
	}
	return nil
}
