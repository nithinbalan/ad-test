package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type KRIRegisterWorkflowRepositoryInterface interface {
	GetWorkflowList(kriRegisterID uint) (*models.KRIRegister, error)
	CreateKRIRegisterWorkflow(kriRegisterWorkflow *models.KRIRegisterWorkflow, tx ...utility.DatabaseTransaction) (uint, error)
	UpdateKRIRegisterWorkflow(kriRegisterWorkflow *models.KRIRegisterWorkflow, tx ...utility.DatabaseTransaction) error
	CreateKRIRegisterWorkflowParticipant(workflowParticipant *models.KRIRegisterWorkflowParticipant, tx ...utility.DatabaseTransaction) error
	DeleteKRIRegisterWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteKRIRegisterWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteKRIRegisterWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error

	GetKRIRegisterWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.KRIRegisterWorkflow, error)
	GetKRIRegisterWorkflowByLevel(kriRegisterID uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.KRIRegisterWorkflow, error)

	GetNextLevel(kriRegisterID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.KRIRegisterWorkflow, error)
	CountApprovedParticipants(workflowID uint, tx ...utility.DatabaseTransaction) (int, error)
}

type KRIRegisterWorkflowRepository struct {
	db *gorm.DB
}

func NewKRIRegisterWorkflowRepository(db *gorm.DB) *KRIRegisterWorkflowRepository {
	return &KRIRegisterWorkflowRepository{db: db}
}

func (repo *KRIRegisterWorkflowRepository) GetWorkflowList(kriRegisterID uint) (*models.KRIRegister, error) {
	var kriRegister models.KRIRegister
	result := repo.db.
		Preload("SubmittedBy").
		Preload("KRIRegisterWorkflows").
		Preload("KRIRegisterWorkflows.WorkflowStatus").
		Preload("KRIRegisterWorkflows.KRIRegisterWorkflowParticipants").
		Preload("KRIRegisterWorkflows.KRIRegisterWorkflowParticipants.User").
		Preload("KRIRegisterWorkflows.KRIRegisterWorkflowParticipants.Role").
		Preload("KRIRegisterWorkflows.KRIRegisterWorkflowReviews", "workflow_status_id IN (?)", repo.db.Model(&models.WorkflowStatus{}).Select("id").Where("type IN (?)", []string{"approve", "reject"})).
		Preload("KRIRegisterWorkflows.KRIRegisterWorkflowReviews.WorkflowStatus").
		Preload("KRIRegisterWorkflows.KRIRegisterWorkflowReviews.CreatedByUser").
		Where("id = ?", kriRegisterID).
		Find(&kriRegister)

	if result.Error != nil {
		return nil, result.Error
	}

	return &kriRegister, nil
}

func (repo *KRIRegisterWorkflowRepository) CreateKRIRegisterWorkflow(kriRegisterWorkflow *models.KRIRegisterWorkflow, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(kriRegisterWorkflow).Error; err != nil {
		return 0, err
	}

	return kriRegisterWorkflow.ID, nil
}

func (repo *KRIRegisterWorkflowRepository) UpdateKRIRegisterWorkflow(kriRegisterWorkflow *models.KRIRegisterWorkflow, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(kriRegisterWorkflow).Error; err != nil {
		return err
	}

	return nil
}

func (repo *KRIRegisterWorkflowRepository) CreateKRIRegisterWorkflowParticipant(workflowParticipant *models.KRIRegisterWorkflowParticipant, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(workflowParticipant).Error; err != nil {
		return err
	}

	return nil
}

func (repo *KRIRegisterWorkflowRepository) GetKRIRegisterWorkflowByLevel(kriRegisterID uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.KRIRegisterWorkflow, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var kriRegisterWorkflow *models.KRIRegisterWorkflow
	result := db.
		Preload("WorkflowLevelType").
		Preload("KRIRegisterWorkflowParticipants").
		Preload("KRIRegisterWorkflowReviews").
		Preload("KRIRegisterWorkflowReviews.WorkflowStatus").
		Preload("KRIRegisterWorkflowParticipants.User").
		Where("kri_register_id = ? AND level = ?", kriRegisterID, currentApprovalLevel).
		Find(&kriRegisterWorkflow)

	if result.Error != nil {
		return nil, result.Error
	}

	return kriRegisterWorkflow, nil
}

func (repo *KRIRegisterWorkflowRepository) GetKRIRegisterWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.KRIRegisterWorkflow, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var kriRegisterWorkflow *models.KRIRegisterWorkflow

	result := db.First(&kriRegisterWorkflow, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return kriRegisterWorkflow, nil
}

func (repo *KRIRegisterWorkflowRepository) GetNextLevel(kriRegisterID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.KRIRegisterWorkflow, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var nextWorkflow models.KRIRegisterWorkflow
	result := db.Where("kri_register_id = ? AND level > ?", kriRegisterID, currentLevel).
		Order("level ASC").
		First(&nextWorkflow)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if result.Error != nil {
		return nil, result.Error
	}

	return &nextWorkflow, nil
}

func (repo *KRIRegisterWorkflowRepository) CountApprovedParticipants(workflowID uint, tx ...utility.DatabaseTransaction) (int, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var resubmitID *uint
	err := db.Raw(`
        SELECT MAX(id)
        FROM kri_register_workflow_reviews
        WHERE kri_register_workflow_id = ?
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
	query := db.Model(&models.KRIRegisterWorkflowReview{}).
		Where("kri_register_workflow_id = ? AND workflow_status_id = (SELECT id FROM workflow_statuses WHERE type = 'approve')", workflowID).
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

func (repo *KRIRegisterWorkflowRepository) DeleteKRIRegisterWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", participantIds).Delete(&models.KRIRegisterWorkflowParticipant{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *KRIRegisterWorkflowRepository) DeleteKRIRegisterWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", reviewIds).Delete(&models.KRIRegisterWorkflowReview{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *KRIRegisterWorkflowRepository) DeleteKRIRegisterWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", workflowIds).Delete(&models.KRIRegisterWorkflow{}).Error; err != nil {
		return err
	}

	return nil
}
