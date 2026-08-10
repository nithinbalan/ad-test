package repositories

import (
	"errors"

	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type IdentityReviewCampaignWorkflowRepositoryInterface interface {
	GetWorkflowList(campaignID uint, tx ...utility.DatabaseTransaction) (*models.IdentityReviewCampaign, error)
	GetWorkflowHistories(campaignID uint, tx ...utility.DatabaseTransaction) ([]*models.IdentityReviewCampaignWorkflowHistory, error)
	CreateWorkflow(workflow *models.IdentityReviewCampaignWorkflow, tx ...utility.DatabaseTransaction) (uint, error)
	UpdateWorkflow(workflow *models.IdentityReviewCampaignWorkflow, tx ...utility.DatabaseTransaction) error
	CreateParticipant(participant *models.IdentityReviewCampaignWorkflowParticipant, tx ...utility.DatabaseTransaction) error
	UpdateParticipant(participant *models.IdentityReviewCampaignWorkflowParticipant, tx ...utility.DatabaseTransaction) error

	GetWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.IdentityReviewCampaignWorkflow, error)
	GetWorkflowByLevel(campaignID uint, level int, tx ...utility.DatabaseTransaction) (*models.IdentityReviewCampaignWorkflow, error)
	GetNextLevel(campaignID uint, level int, tx ...utility.DatabaseTransaction) (*models.IdentityReviewCampaignWorkflow, error)
	CountApprovedParticipants(workflowID uint, tx ...utility.DatabaseTransaction) (int, error)

	DeleteAllWorkflowsForCampaign(campaignID uint, tx ...utility.DatabaseTransaction) error
	DeleteReviewsByWorkflowID(workflowID uint, tx ...utility.DatabaseTransaction) error
	ResetParticipantStatuses(workflowID uint, tx ...utility.DatabaseTransaction) error
	UpdateWorkflowStatus(workflowID uint, statusID uint, tx ...utility.DatabaseTransaction) error
}

type IdentityReviewCampaignWorkflowRepository struct {
	db *gorm.DB
}

func NewIdentityReviewCampaignWorkflowRepository(db *gorm.DB) IdentityReviewCampaignWorkflowRepositoryInterface {
	return &IdentityReviewCampaignWorkflowRepository{
		db: db,
	}
}

func (repo *IdentityReviewCampaignWorkflowRepository) GetWorkflowList(campaignID uint, tx ...utility.DatabaseTransaction) (*models.IdentityReviewCampaign, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var campaign models.IdentityReviewCampaign
	result := db.
		Preload("CampaignStatus").
		Preload("SubmittedByUser").
		Preload("CampaignWorkflows").
		Preload("CampaignWorkflows.WorkflowStatus").
		Preload("CampaignWorkflows.IdentityReviewCampaignWorkflowParticipants").
		Preload("CampaignWorkflows.IdentityReviewCampaignWorkflowParticipants.User").
		Preload("CampaignWorkflows.IdentityReviewCampaignWorkflowParticipants.Role").
		Preload("CampaignWorkflows.IdentityReviewCampaignWorkflowReviews", "workflow_status_id IN (?)", repo.db.Model(&models.WorkflowStatus{}).Select("id").Where("type IN (?)", []string{"approve", "reject"})).
		Preload("CampaignWorkflows.IdentityReviewCampaignWorkflowReviews.WorkflowStatus").
		Preload("CampaignWorkflows.IdentityReviewCampaignWorkflowReviews.CreatedByUser").
		Where("id = ?", campaignID).
		First(&campaign)

	if result.Error != nil {
		return nil, result.Error
	}

	return &campaign, nil
}

func (repo *IdentityReviewCampaignWorkflowRepository) GetWorkflowHistories(campaignID uint, tx ...utility.DatabaseTransaction) ([]*models.IdentityReviewCampaignWorkflowHistory, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var histories []*models.IdentityReviewCampaignWorkflowHistory
	result := db.
		Preload("WorkflowStatus").
		Preload("CreatedByUser").
		Where("identity_review_campaign_id = ?", campaignID).
		Order("created_at DESC").
		Find(&histories)

	if result.Error != nil {
		return nil, result.Error
	}

	return histories, nil
}

func (repo *IdentityReviewCampaignWorkflowRepository) CreateWorkflow(workflow *models.IdentityReviewCampaignWorkflow, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(workflow).Error; err != nil {
		return 0, err
	}

	return workflow.ID, nil
}

func (repo *IdentityReviewCampaignWorkflowRepository) UpdateWorkflow(workflow *models.IdentityReviewCampaignWorkflow, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(workflow).Error; err != nil {
		return err
	}

	return nil
}

func (repo *IdentityReviewCampaignWorkflowRepository) CreateParticipant(participant *models.IdentityReviewCampaignWorkflowParticipant, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(participant).Error; err != nil {
		return err
	}

	return nil
}

func (repo *IdentityReviewCampaignWorkflowRepository) UpdateParticipant(participant *models.IdentityReviewCampaignWorkflowParticipant, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(participant).Error; err != nil {
		return err
	}

	return nil
}

func (repo *IdentityReviewCampaignWorkflowRepository) GetWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.IdentityReviewCampaignWorkflow, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var workflow *models.IdentityReviewCampaignWorkflow
	result := db.First(&workflow, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return workflow, nil
}

func (repo *IdentityReviewCampaignWorkflowRepository) GetWorkflowByLevel(campaignID uint, level int, tx ...utility.DatabaseTransaction) (*models.IdentityReviewCampaignWorkflow, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var workflow *models.IdentityReviewCampaignWorkflow
	result := db.
		Preload("WorkflowLevelType").
		Preload("IdentityReviewCampaignWorkflowParticipants").
		Preload("IdentityReviewCampaignWorkflowParticipants.User").
		Preload("IdentityReviewCampaignWorkflowReviews").
		Preload("IdentityReviewCampaignWorkflowReviews.WorkflowStatus").
		Where("identity_review_campaign_id = ? AND level = ?", campaignID, level).
		First(&workflow)

	if result.Error != nil {
		return nil, result.Error
	}

	return workflow, nil
}

func (repo *IdentityReviewCampaignWorkflowRepository) GetNextLevel(campaignID uint, level int, tx ...utility.DatabaseTransaction) (*models.IdentityReviewCampaignWorkflow, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var nextWorkflow models.IdentityReviewCampaignWorkflow
	result := db.
		Where("identity_review_campaign_id = ? AND level > ?", campaignID, level).
		Order("level ASC").
		First(&nextWorkflow)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if result.Error != nil {
		return nil, result.Error
	}

	return &nextWorkflow, nil
}

func (repo *IdentityReviewCampaignWorkflowRepository) CountApprovedParticipants(workflowID uint, tx ...utility.DatabaseTransaction) (int, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var resubmitID *uint
	err := db.Raw(`
		SELECT MAX(id)
		FROM identity_review_campaign_workflow_reviews
		WHERE identity_review_campaign_workflow_id = ?
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
		resubmitID = new(uint)
	}

	var count int64
	query := db.Model(&models.IdentityReviewCampaignWorkflowReview{}).
		Where("identity_review_campaign_workflow_id = ? AND workflow_status_id = (SELECT id FROM workflow_statuses WHERE type = 'approve')", workflowID).
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

func (repo *IdentityReviewCampaignWorkflowRepository) DeleteAllWorkflowsForCampaign(campaignID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var workflowIDs []uint
	if err := db.Model(&models.IdentityReviewCampaignWorkflow{}).
		Where("identity_review_campaign_id = ?", campaignID).
		Pluck("id", &workflowIDs).Error; err != nil {
		return err
	}

	if len(workflowIDs) == 0 {
		return nil
	}

	if err := db.Where("identity_review_campaign_workflow_id IN ?", workflowIDs).Delete(&models.IdentityReviewCampaignWorkflowReview{}).Error; err != nil {
		return err
	}

	if err := db.Where("identity_review_campaign_workflow_id IN ?", workflowIDs).Delete(&models.IdentityReviewCampaignWorkflowParticipant{}).Error; err != nil {
		return err
	}

	return db.Where("id IN ?", workflowIDs).Delete(&models.IdentityReviewCampaignWorkflow{}).Error
}

func (repo *IdentityReviewCampaignWorkflowRepository) DeleteReviewsByWorkflowID(workflowID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("identity_review_campaign_workflow_id = ?", workflowID).Delete(&models.IdentityReviewCampaignWorkflowReview{}).Error
}

func (repo *IdentityReviewCampaignWorkflowRepository) ResetParticipantStatuses(workflowID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Model(&models.IdentityReviewCampaignWorkflowParticipant{}).
		Where("identity_review_campaign_workflow_id = ?", workflowID).
		Update("workflow_status_id", nil).Error
}

func (repo *IdentityReviewCampaignWorkflowRepository) UpdateWorkflowStatus(workflowID uint, statusID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Model(&models.IdentityReviewCampaignWorkflow{}).
		Where("id = ?", workflowID).
		Update("workflow_status_id", statusID).Error
}
