package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type InternalCommitteeWorkflowRepositoryInterface interface {
	GetWorkflowList(internalCommitteeID uint, latestVersionID uint, tx ...utility.DatabaseTransaction) (*models.InternalCommittee, error)
	CreateInternalCommitteeWorkflow(internalCommitteeWorkflow *models.InternalCommitteeWorkflow, tx ...utility.DatabaseTransaction) (uint, error)
	UpdateInternalCommitteeWorkflow(internalCommitteeWorkflow *models.InternalCommitteeWorkflow, tx ...utility.DatabaseTransaction) error
	CreateInternalCommitteeWorkflowParticipant(internalCommitteeWorkflowParticipant *models.InternalCommitteeWorkflowParticipant, tx ...utility.DatabaseTransaction) error
	DeleteInternalCommitteeWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteInternalCommitteeWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteInternalCommitteeWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error

	GetInternalCommitteeWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.InternalCommitteeWorkflow, error)
	GetInternalCommitteeWorkflowByLevel(internalCommitteeVersionID uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.InternalCommitteeWorkflow, error)

	GetNextLevel(internalCommitteeID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.InternalCommitteeWorkflow, error)
	CountApprovedParticipants(workflowID uint, tx ...utility.DatabaseTransaction) (int, error)
}

type InternalCommitteeWorkflowRepository struct {
	db *gorm.DB
}

func NewInternalCommitteeWorkflowRepository(db *gorm.DB) *InternalCommitteeWorkflowRepository {
	return &InternalCommitteeWorkflowRepository{db: db}
}

func (repo *InternalCommitteeWorkflowRepository) GetWorkflowList(committeeID uint, latestVersionID uint, tx ...utility.DatabaseTransaction) (*models.InternalCommittee, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var internalCommittee models.InternalCommittee
	result := db.
		Preload("InternalCommitteeVersions", "is_latest_version = ?", true).
		Preload("InternalCommitteeVersions.SubmittedByUser").
		Preload("InternalCommitteeWorkflows", "internal_committee_version_id = ?", latestVersionID).
		Preload("InternalCommitteeWorkflows.WorkflowStatus").
		Preload("InternalCommitteeWorkflows.InternalCommitteeWorkflowParticipants").
		Preload("InternalCommitteeWorkflows.InternalCommitteeWorkflowParticipants.User").
		Preload("InternalCommitteeWorkflows.InternalCommitteeWorkflowParticipants.Role").
		Preload("InternalCommitteeWorkflows.InternalCommitteeWorkflowReviews", "workflow_status_id IN (?)", db.Model(&models.WorkflowStatus{}).Select("id").Where("type IN (?)", []string{"approve", "revert"})).
		Preload("InternalCommitteeWorkflows.InternalCommitteeWorkflowReviews.WorkflowStatus").
		Preload("InternalCommitteeWorkflows.InternalCommitteeWorkflowReviews.CreatedByUser").
		Where("id = ?", committeeID).
		First(&internalCommittee)

	if result.Error != nil {
		return nil, result.Error
	}

	return &internalCommittee, nil
}

func (repo *InternalCommitteeWorkflowRepository) CreateInternalCommitteeWorkflow(internalCommitteeWorkflow *models.InternalCommitteeWorkflow, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(internalCommitteeWorkflow).Error; err != nil {
		return 0, err
	}

	return internalCommitteeWorkflow.ID, nil
}

func (repo *InternalCommitteeWorkflowRepository) UpdateInternalCommitteeWorkflow(internalCommitteeWorkflow *models.InternalCommitteeWorkflow, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Use Updates with Select to only update specific fields, avoiding preloaded association issues
	if err := db.Model(&models.InternalCommitteeWorkflow{}).
		Where("id = ?", internalCommitteeWorkflow.ID).
		Updates(map[string]interface{}{
			"workflow_status_id": internalCommitteeWorkflow.WorkflowStatusID,
		}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *InternalCommitteeWorkflowRepository) CreateInternalCommitteeWorkflowParticipant(internalCommitteeWorkflowParticipant *models.InternalCommitteeWorkflowParticipant, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(internalCommitteeWorkflowParticipant).Error; err != nil {
		return err
	}

	return nil
}

func (repo *InternalCommitteeWorkflowRepository) GetInternalCommitteeWorkflowByLevel(internalCommitteeVersionID uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.InternalCommitteeWorkflow, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var internalCommitteeWorkflow models.InternalCommitteeWorkflow
	result := db.
		Preload("WorkflowLevelType").
		Preload("InternalCommitteeWorkflowParticipants").
		Preload("InternalCommitteeWorkflowReviews").
		Preload("InternalCommitteeWorkflowReviews.WorkflowStatus").
		Preload("InternalCommitteeWorkflowParticipants.User").
		Where("internal_committee_version_id = ? AND level = ?", internalCommitteeVersionID, currentApprovalLevel).
		First(&internalCommitteeWorkflow)

	if result.Error != nil {
		return nil, result.Error
	}

	return &internalCommitteeWorkflow, nil
}

func (repo *InternalCommitteeWorkflowRepository) GetInternalCommitteeWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.InternalCommitteeWorkflow, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var internalCommitteeWorkflow *models.InternalCommitteeWorkflow

	result := db.First(&internalCommitteeWorkflow, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return internalCommitteeWorkflow, nil
}

func (repo *InternalCommitteeWorkflowRepository) GetNextLevel(internalCommitteeID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.InternalCommitteeWorkflow, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var nextWorkflow models.InternalCommitteeWorkflow
	result := db.Where("internal_committee_id = ? AND level > ?", internalCommitteeID, currentLevel).
		Order("level ASC").  // Ensure we are ascending by level
		First(&nextWorkflow) // Retrieve the first result only which will be the next level

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil // No further levels available
	} else if result.Error != nil {
		return nil, result.Error // Other error occurred
	}

	return &nextWorkflow, nil
}

// func (repo *RiskWorkflowRepository) CountApprovedParticipants(workflowID uint, tx ...utility.DatabaseTransaction) (int, error) {
// 	db := repo.db
// 	if len(tx) > 0 && tx[0] != nil {
// 		db = tx[0].DB()
// 	}
// 	var count int64
// 	err := db.Model(&models.RiskWorkflowReview{}).
// 		Where("risk_workflow_id = ? AND workflow_status_id = (SELECT id FROM workflow_statuses WHERE type = 'approve')", workflowID).
// 		Distinct("created_by_id").
// 		Count(&count).Error

// 	if err != nil {
// 		return 0, err
// 	}

//		return int(count), nil
//	}

func (repo *InternalCommitteeWorkflowRepository) CountApprovedParticipants(workflowID uint, tx ...utility.DatabaseTransaction) (int, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var resubmitID *uint // Pointer to handle NULL values
	err := db.Raw(`
        SELECT MAX(id)
        FROM internal_committee_workflow_reviews
        WHERE internal_committee_workflow_id = ?
          AND workflow_status_id = (
              SELECT id
              FROM workflow_statuses
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
	query := db.Model(&models.InternalCommitteeWorkflowReview{}).
		Where("internal_committee_workflow_id = ? AND workflow_status_id = (SELECT id FROM workflow_statuses WHERE type = 'approve')", workflowID).
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

func (repo *InternalCommitteeWorkflowRepository) DeleteInternalCommitteeWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", participantIds).Delete(&models.InternalCommitteeWorkflowParticipant{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *InternalCommitteeWorkflowRepository) DeleteInternalCommitteeWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", reviewIds).Delete(&models.InternalCommitteeWorkflowReview{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *InternalCommitteeWorkflowRepository) DeleteInternalCommitteeWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", workflowIds).Delete(&models.InternalCommitteeWorkflow{}).Error; err != nil {
		return err
	}

	return nil
}
