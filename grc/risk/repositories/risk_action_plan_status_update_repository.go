package repositories

import (
	"jethur/grc/models"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type RiskActionPlanStatusUpdateRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(riskActionPlanID uint, sort string, sortBy string, keyword string) ([]*riskDTOs.RiskActionPlanStatusUpdateListDTO, error)
	GetByID(id uint) (*models.RiskActionPlanStatusUpdate, error)
	Create(RiskActionPlanStatusUpdate *models.RiskActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error)
	Update(RiskActionPlanStatusUpdate *models.RiskActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) error

	//Document
	CreateDocuments(docs []*models.RiskActionPlanStatusUpdateDocument, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.RiskActionPlanStatusUpdateDocument, error)
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error
}

type RiskActionPlanStatusUpdateRepository struct {
	db *gorm.DB
}

func NewRiskActionPlanStatusUpdateRepository(db *gorm.DB) *RiskActionPlanStatusUpdateRepository {
	return &RiskActionPlanStatusUpdateRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *RiskActionPlanStatusUpdateRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *RiskActionPlanStatusUpdateRepository) GetAll(riskActionPlanID uint, sort string, sortBy string, keyword string) ([]*riskDTOs.RiskActionPlanStatusUpdateListDTO, error) {

	var riskActionPlanUpdates []*riskDTOs.RiskActionPlanStatusUpdateListDTO

	query := repo.db.Table("risk_action_plan_status_updates").
		Select(`
            risk_action_plan_status_updates.*,
            risk_action_plans.title AS risk_action_plan_title,
            action_plan_statuses.title_en AS action_plan_status_title_en,
            action_plan_statuses.title_ar AS action_plan_status_title_ar,
            action_plan_statuses.color_code AS action_plan_status_color_code,
            created_user.first_name_en AS created_user_first_name_en,
            created_user.first_name_ar AS created_user_first_name_ar,
            created_user.last_name_en AS created_user_last_name_en,
            created_user.last_name_ar AS created_user_last_name_ar,
            updated_user.first_name_en AS updated_user_first_name_en,
            updated_user.first_name_ar AS updated_user_first_name_ar,
            updated_user.last_name_en AS updated_user_last_name_en,
            updated_user.last_name_ar AS updated_user_last_name_ar 
        `).
		Joins(`
            LEFT JOIN risk_action_plans ON risk_action_plans.id = risk_action_plan_status_updates.risk_action_plan_id
            LEFT JOIN action_plan_statuses ON action_plan_statuses.id = risk_action_plan_status_updates.action_plan_status_id
            LEFT JOIN users AS created_user ON created_user.id = risk_action_plan_status_updates.created_by
            LEFT JOIN users AS updated_user ON updated_user.id = risk_action_plan_status_updates.updated_by
        `).
		Where("risk_action_plan_status_updates.risk_action_plan_id = ?", riskActionPlanID)

	if keyword != "" {
		query = query.Where(`
            (risk_action_plans.title LIKE ? OR
            risk_action_plan_status_updates.comment LIKE ? OR
            action_plan_statuses.title_en LIKE ? OR
            action_plan_statuses.title_ar LIKE ?)`,
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if sortBy == "" {
		sortBy = "risk_action_plan_status_updates.id"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}
	query = query.Order(sortBy + " " + sort)
	result := query.Find(&riskActionPlanUpdates)

	return riskActionPlanUpdates, result.Error
}

func (repo *RiskActionPlanStatusUpdateRepository) GetByID(id uint) (*models.RiskActionPlanStatusUpdate, error) {
	var riskActionPlanStatusUpdate *models.RiskActionPlanStatusUpdate
	result := repo.db.First(&riskActionPlanStatusUpdate, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return riskActionPlanStatusUpdate, nil
}

func (repo *RiskActionPlanStatusUpdateRepository) Create(RiskActionPlanStatusUpdate *models.RiskActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(RiskActionPlanStatusUpdate).Error; err != nil {
		return 0, err
	}

	return RiskActionPlanStatusUpdate.ID, nil
}

func (repo *RiskActionPlanStatusUpdateRepository) Update(RiskActionPlanStatusUpdate *models.RiskActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(RiskActionPlanStatusUpdate).Error; err != nil {
		return err
	}
	return nil
}

func (repo *RiskActionPlanStatusUpdateRepository) CreateDocuments(docs []*models.RiskActionPlanStatusUpdateDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Save all documents in a single transaction
	if err := db.Create(&docs).Error; err != nil {
		return err
	}
	return nil
}

func (repo *RiskActionPlanStatusUpdateRepository) GetDocumentByID(id uint) (*models.RiskActionPlanStatusUpdateDocument, error) {
	var document models.RiskActionPlanStatusUpdateDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (repo *RiskActionPlanStatusUpdateRepository) DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete documents by IDs
	if err := db.Where("id IN (?)", docIds).Delete(&models.RiskActionPlanStatusUpdateDocument{}).Error; err != nil {
		return err
	}

	return nil
}
