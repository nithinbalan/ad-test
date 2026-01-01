package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ThirdPartyAssessmentWorkflowRepositoryInterface interface {
	GetWorkflowList(assessmentID uint) (*models.ThirdPartyAssessment, error)
	CreateTprmAssessmentWorkflow(tprmAssessmentWorkflow *models.ThirdPartyAssessmentWorkflow, tx ...utility.DatabaseTransaction) (uint, error)
	UpdateTprmAssessmentWorkflow(tprmAssessmentWorkflow *models.ThirdPartyAssessmentWorkflow, tx ...utility.DatabaseTransaction) error
	CreateTprmAssessmentWorkflowParticipant(docWorkflowParticipant *models.ThirdPartyAssessmentWorkflowParticipant, tx ...utility.DatabaseTransaction) error
	DeleteTprmAssessmentWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteTprmAssessmentWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteTprmAssessmentWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error

	GetTprmAssessmentWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.ThirdPartyAssessmentWorkflow, error)
	GetTprmAssessmentWorkflowByLevel(assessmentID uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.ThirdPartyAssessmentWorkflow, error)

	GetNextLevel(assessmentID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.ThirdPartyAssessmentWorkflow, error)
	CountApprovedParticipants(workflowID uint, tx ...utility.DatabaseTransaction) (int, error)
}

type ThirdPartyAssessmentWorkflowRepository struct {
	db *gorm.DB
}

func NewThirdPartyAssessmentWorkflowRepository(db *gorm.DB) ThirdPartyAssessmentWorkflowRepositoryInterface {
	return &ThirdPartyAssessmentWorkflowRepository{db: db}
}

func (repo *ThirdPartyAssessmentWorkflowRepository) GetWorkflowList(assessmentID uint) (*models.ThirdPartyAssessment, error) {
	var tprmAssessment models.ThirdPartyAssessment
	result := repo.db.
		Preload("SubmittedByUser").
		Preload("ThirdPartyAssessmentWorkflows").
		Preload("ThirdPartyAssessmentWorkflows.WorkflowStatus").
		Preload("ThirdPartyAssessmentWorkflows.ThirdPartyAssessmentWorkflowParticipants").
		Preload("ThirdPartyAssessmentWorkflows.ThirdPartyAssessmentWorkflowParticipants.User").
		Preload("ThirdPartyAssessmentWorkflows.ThirdPartyAssessmentWorkflowParticipants.Role").
		Preload("ThirdPartyAssessmentWorkflows.ThirdPartyAssessmentWorkflowReviews", "workflow_status_id IN (?)", repo.db.Model(&models.WorkflowStatus{}).Select("id").Where("type IN (?)", []string{"approve", "reject"})).
		Preload("ThirdPartyAssessmentWorkflows.ThirdPartyAssessmentWorkflowReviews.WorkflowStatus").
		Preload("ThirdPartyAssessmentWorkflows.ThirdPartyAssessmentWorkflowReviews.CreatedByUser").
		Where("id = ?", assessmentID).
		Find(&tprmAssessment)

	if result.Error != nil {
		return nil, result.Error
	}

	return &tprmAssessment, nil
}

func (repo *ThirdPartyAssessmentWorkflowRepository) CreateTprmAssessmentWorkflow(tprmAssessmentWorkflow *models.ThirdPartyAssessmentWorkflow, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(tprmAssessmentWorkflow).Error; err != nil {
		return 0, err
	}

	return tprmAssessmentWorkflow.ID, nil
}

func (repo *ThirdPartyAssessmentWorkflowRepository) UpdateTprmAssessmentWorkflow(tprmAssessmentWorkflow *models.ThirdPartyAssessmentWorkflow, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(tprmAssessmentWorkflow).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ThirdPartyAssessmentWorkflowRepository) CreateTprmAssessmentWorkflowParticipant(docWorkflowParticipant *models.ThirdPartyAssessmentWorkflowParticipant, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(docWorkflowParticipant).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ThirdPartyAssessmentWorkflowRepository) GetTprmAssessmentWorkflowByLevel(assessmentID uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.ThirdPartyAssessmentWorkflow, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var workflow *models.ThirdPartyAssessmentWorkflow
	result := db.
		Preload("WorkflowLevelType").
		Preload("ThirdPartyAssessmentWorkflowParticipants").
		Preload("ThirdPartyAssessmentWorkflowReviews").
		Preload("ThirdPartyAssessmentWorkflowReviews.WorkflowStatus").
		Preload("ThirdPartyAssessmentWorkflowParticipants.User").
		Where("third_party_assessment_id = ? AND level = ?", assessmentID, currentApprovalLevel).
		Find(&workflow)

	if result.Error != nil {
		return nil, result.Error
	}

	return workflow, nil
}

func (repo *ThirdPartyAssessmentWorkflowRepository) GetTprmAssessmentWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.ThirdPartyAssessmentWorkflow, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var tprmWorkflow *models.ThirdPartyAssessmentWorkflow

	result := db.First(&tprmWorkflow, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return tprmWorkflow, nil
}

func (repo *ThirdPartyAssessmentWorkflowRepository) GetNextLevel(riskID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.ThirdPartyAssessmentWorkflow, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var nextWorkflow models.ThirdPartyAssessmentWorkflow
	result := db.Where("third_party_assessment_id = ? AND level > ?", riskID, currentLevel).
		Order("level ASC").
		First(&nextWorkflow)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if result.Error != nil {
		return nil, result.Error
	}

	return &nextWorkflow, nil
}

func (repo *ThirdPartyAssessmentWorkflowRepository) CountApprovedParticipants(workflowID uint, tx ...utility.DatabaseTransaction) (int, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var resubmitID *uint
	err := db.Raw(`
        SELECT MAX(id)
        FROM tprm_assessment_workflow_reviews
        WHERE third_party_assessment_id = ?
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
	query := db.Model(&models.ThirdPartyAssessmentWorkflowReview{}).
		Where("tprm_assessment_workflow_id = ? AND workflow_status_id = (SELECT id FROM workflow_statuses WHERE type = 'approve')", workflowID).
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

func (repo *ThirdPartyAssessmentWorkflowRepository) DeleteTprmAssessmentWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", participantIds).Delete(&models.ThirdPartyAssessmentWorkflowParticipant{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ThirdPartyAssessmentWorkflowRepository) DeleteTprmAssessmentWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", reviewIds).Delete(&models.ThirdPartyAssessmentWorkflowReview{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ThirdPartyAssessmentWorkflowRepository) DeleteTprmAssessmentWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", workflowIds).Delete(&models.ThirdPartyAssessmentWorkflow{}).Error; err != nil {
		return err
	}

	return nil
}
