package repositories

import (
	"database/sql"
	"jethur/grc/models"
	settingsDTOs "jethur/grc/settings/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type WorkflowRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAllWorkflowsByType(workflowType string) ([]*settingsDTOs.WorkflowListDTO, error)
	GetWorkflowLevelByID(id uint, tx ...utility.DatabaseTransaction) (*models.WorkflowLevel, error)
	GetWorkflowLevels(workflowID uint) (*models.Workflow, error)
	CreateWorkflowLevel(workflowLevel *models.WorkflowLevel, tx ...utility.DatabaseTransaction) (uint, error)
	CreateParticipant(participant *models.WorkflowLevelParticipant, tx ...utility.DatabaseTransaction) error
	DeleteWorkflowLevelParticipants(workflowLevelID uint, tx ...utility.DatabaseTransaction) error
	DeleteWorkflowLevel(workflowLevelID uint, tx ...utility.DatabaseTransaction) error
	GetMaxLevelByWorkflowID(workflowID uint, tx ...utility.DatabaseTransaction) (int, error)
	DecrementSubsequentLevels(workflowID uint, deletedLevel int, tx ...utility.DatabaseTransaction) error
	IncrementSubsequentLevels(workflowID uint, startingLevel uint, tx ...utility.DatabaseTransaction) error
}

// Implementation
type WorkflowRepository struct {
	db *gorm.DB
}

func NewWorkflowRepository(db *gorm.DB) *WorkflowRepository {
	return &WorkflowRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *WorkflowRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *WorkflowRepository) GetAllWorkflowsByType(workflowType string) ([]*settingsDTOs.WorkflowListDTO, error) {
	var workflows []*settingsDTOs.WorkflowListDTO

	query := repo.db.Table("workflows").
		Select(`
			workflows.id,
			workflows.type,
			workflows.document_type_id,
			workflows.status,
			workflows.created_by_id,
			workflows.updated_by_id,
			workflows.created_at,
			workflows.updated_at,
			document_types.title_en AS document_type_title_en,
			document_types.title_ar AS document_type_title_ar,
			created_user.first_name_en AS created_user_first_name_en,
			created_user.first_name_ar AS created_user_first_name_ar,
			created_user.email AS created_user_email,
			updated_user.first_name_en AS updated_user_first_name_en,
			updated_user.first_name_ar AS updated_user_first_name_ar,
			updated_user.email AS updated_user_email
		`).
		Joins("LEFT JOIN document_types ON document_types.id = workflows.document_type_id").
		Joins("LEFT JOIN users AS created_user ON created_user.id = workflows.created_by_id").
		Joins("LEFT JOIN users AS updated_user ON updated_user.id = workflows.updated_by_id").
		Where("workflows.type = ?", workflowType)

	err := query.Find(&workflows).Error
	if err != nil {
		return nil, err
	}

	return workflows, nil
}

func (repo *WorkflowRepository) GetWorkflowLevels(workflowID uint) (*models.Workflow, error) {
	var workflow *models.Workflow

	if err := repo.db.
		Preload("WorkflowLevels", func(db *gorm.DB) *gorm.DB {
			return db.Order("workflow_levels.level ASC") // Ordering WorkflowLevels by level ascending
		}).
		Preload("WorkflowLevels.WorkflowLevelType").
		Preload("WorkflowLevels.Participants").
		Preload("WorkflowLevels.Participants.User").
		Preload("WorkflowLevels.Participants.Role").
		Where("id = ?", workflowID). // Assuming the primary key field is "id"
		Find(&workflow).Error; err != nil {
		return nil, err
	}
	return workflow, nil
}

func (repo *WorkflowRepository) GetWorkflowLevelByID(id uint, tx ...utility.DatabaseTransaction) (*models.WorkflowLevel, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var workflowLevel *models.WorkflowLevel
	result := db.First(&workflowLevel, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return workflowLevel, nil
}

func (repo *WorkflowRepository) CreateWorkflowLevel(workflowLevel *models.WorkflowLevel, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(workflowLevel).Error; err != nil {
		return 0, err
	}

	return workflowLevel.ID, nil
}

func (repo *WorkflowRepository) CreateParticipant(participant *models.WorkflowLevelParticipant, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(participant).Error; err != nil {
		return err
	}

	return nil
}

func (repo *WorkflowRepository) DeleteWorkflowLevelParticipants(workflowLevelID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("workflow_level_id = ?", workflowLevelID).Delete(&models.WorkflowLevelParticipant{}).Error
}

func (repo *WorkflowRepository) DeleteWorkflowLevel(workflowLevelID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Delete(&models.WorkflowLevel{}, workflowLevelID).Error
}

func (repo *WorkflowRepository) GetMaxLevelByWorkflowID(workflowID uint, tx ...utility.DatabaseTransaction) (int, error) {
	var maxLevel sql.NullInt64 // Use sql.NullInt64 to handle possible NULL values
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	err := db.Model(&models.WorkflowLevel{}).Where("workflow_id = ?", workflowID).Select("max(level)").Row().Scan(&maxLevel)
	if err != nil {
		return 0, err
	}
	if !maxLevel.Valid {
		return 0, nil // Return 0 if no levels exist (maxLevel is NULL)
	}
	return int(maxLevel.Int64), nil
}

func (repo *WorkflowRepository) DecrementSubsequentLevels(workflowID uint, deletedLevel int, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Model(&models.WorkflowLevel{}).
		Where("workflow_id = ? AND level > ?", workflowID, deletedLevel).
		UpdateColumn("level", gorm.Expr("level - 1")).Error
}

func (repo *WorkflowRepository) IncrementSubsequentLevels(workflowID uint, startingLevel uint, tx ...utility.DatabaseTransaction) error {
	var db *gorm.DB = repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Increment levels for all levels above the starting level
	return db.Model(&models.WorkflowLevel{}).
		Where("workflow_id = ? AND level >= ?", workflowID, startingLevel).
		UpdateColumn("level", gorm.Expr("level + 1")).Error
}
