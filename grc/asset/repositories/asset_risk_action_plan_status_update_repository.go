package repositories

import (
	assetDTOs "jethur/grc/asset/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AssetRiskActionPlanStatusUpdateRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(riskActionPlanID uint, sort string, sortBy string, keyword string) ([]*assetDTOs.AssetRiskActionPlanStatusUpdateListDTO, error)
	GetByID(id uint) (*models.RiskActionPlanStatusUpdate, error)
	Create(riskActionPlanStatusUpdate *models.RiskActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error)
	Update(riskActionPlanStatusUpdate *models.RiskActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) error
	GetLatestRiskActionPlanStatusUpdate(riskActionPlanID uint) (*assetDTOs.AssetRiskActionPlanStatusUpdateListDTO, error)

	// Document
	CreateDocuments(docs []*models.RiskActionPlanStatusUpdateDocument, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.RiskActionPlanStatusUpdateDocument, error)
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error

	RiskActionPlanUpdateStatus(riskActionPlan *models.RiskActionPlan, tx ...utility.DatabaseTransaction) error
}

type AssetRiskActionPlanStatusUpdateRepository struct {
	db *gorm.DB
}

func NewAssetRiskActionPlanStatusUpdateRepository(db *gorm.DB) *AssetRiskActionPlanStatusUpdateRepository {
	return &AssetRiskActionPlanStatusUpdateRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *AssetRiskActionPlanStatusUpdateRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *AssetRiskActionPlanStatusUpdateRepository) GetAll(riskActionPlanID uint, sort string, sortBy string, keyword string) ([]*assetDTOs.AssetRiskActionPlanStatusUpdateListDTO, error) {

	var riskActionPlanUpdates []*assetDTOs.AssetRiskActionPlanStatusUpdateListDTO

	query := repo.db.Table("risk_action_plan_status_updates").
		Select(`
            risk_action_plan_status_updates.*,
            risk_action_plans.title AS risk_action_plan_title,
            action_plan_statuses.title_en AS action_plan_status_title_en,
            action_plan_statuses.title_ar AS action_plan_status_title_ar,
            action_plan_statuses.color_code AS action_plan_status_color_code,
			action_plan_statuses.type AS action_plan_status_type,
            created_user.id AS created_user_id,
        	created_user.first_name_en AS created_user_first_name_en,
			created_user.first_name_ar AS created_user_first_name_ar,
			created_user.last_name_en AS created_user_last_name_en,
			created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
			created_user.color AS created_user_color,
			created_user.email AS created_user_email,
			created_user.profile_image_path AS created_user_profile_image_path,
			updated_user.id AS updated_user_id,
    		updated_user.first_name_en AS updated_user_first_name_en,
			updated_user.first_name_ar AS updated_user_first_name_ar,
			updated_user.last_name_en AS updated_user_last_name_en,
			updated_user.last_name_ar AS updated_user_last_name_ar,
			updated_user.initial AS updated_user_initial,
			updated_user.color AS updated_user_color,
			updated_user.email AS updated_user_email,
			updated_user.profile_image_path AS updated_user_profile_image_path
        `).
		Joins(`
            LEFT JOIN risk_action_plans ON risk_action_plans.id = risk_action_plan_status_updates.risk_action_plan_id
            LEFT JOIN action_plan_statuses ON action_plan_statuses.id = risk_action_plan_status_updates.action_plan_status_id
            LEFT JOIN users AS created_user ON created_user.id = risk_action_plan_status_updates.created_by_id
            LEFT JOIN users AS updated_user ON updated_user.id = risk_action_plan_status_updates.updated_by_id
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

	for _, statusUpdate := range riskActionPlanUpdates {
		var documents []*assetDTOs.AssetRiskActionPlanStatusUpdateDocumentDTO

		repo.db.Model(&models.RiskActionPlanStatusUpdateDocument{}).Where("risk_action_plan_status_update_id = ?", statusUpdate.ID).Find(&documents)
		statusUpdate.Documents = documents
	}

	return riskActionPlanUpdates, result.Error
}

func (repo *AssetRiskActionPlanStatusUpdateRepository) GetByID(id uint) (*models.RiskActionPlanStatusUpdate, error) {
	var riskActionPlanStatusUpdate *models.RiskActionPlanStatusUpdate
	result := repo.db.Preload("Documents").First(&riskActionPlanStatusUpdate, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return riskActionPlanStatusUpdate, nil
}

func (repo *AssetRiskActionPlanStatusUpdateRepository) Create(riskActionPlanStatusUpdate *models.RiskActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(riskActionPlanStatusUpdate).Error; err != nil {
		return 0, err
	}

	return riskActionPlanStatusUpdate.ID, nil
}

func (repo *AssetRiskActionPlanStatusUpdateRepository) Update(riskActionPlanStatusUpdate *models.RiskActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(riskActionPlanStatusUpdate).Error; err != nil {
		return err
	}
	return nil
}

func (repo *AssetRiskActionPlanStatusUpdateRepository) CreateDocuments(docs []*models.RiskActionPlanStatusUpdateDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&docs).Error; err != nil {
		return err
	}
	return nil
}

func (repo *AssetRiskActionPlanStatusUpdateRepository) GetDocumentByID(id uint) (*models.RiskActionPlanStatusUpdateDocument, error) {
	var document models.RiskActionPlanStatusUpdateDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (repo *AssetRiskActionPlanStatusUpdateRepository) DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", docIds).Delete(&models.RiskActionPlanStatusUpdateDocument{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AssetRiskActionPlanStatusUpdateRepository) GetLatestRiskActionPlanStatusUpdate(riskActionPlanID uint) (*assetDTOs.AssetRiskActionPlanStatusUpdateListDTO, error) {

	var riskActionPlanUpdate *assetDTOs.AssetRiskActionPlanStatusUpdateListDTO

	query := repo.db.Table("risk_action_plan_status_updates").
		Select(`
            risk_action_plan_status_updates.*,
            risk_action_plans.title AS risk_action_plan_title,
            action_plan_statuses.title_en AS action_plan_status_title_en,
            action_plan_statuses.title_ar AS action_plan_status_title_ar,
            action_plan_statuses.color_code AS action_plan_status_color_code,
			action_plan_statuses.type AS action_plan_status_type,
        	created_user.id AS created_user_id,
        		created_user.first_name_en AS created_user_first_name_en,
        	created_user.first_name_ar AS created_user_first_name_ar,
        	created_user.last_name_en AS created_user_last_name_en,
        	created_user.last_name_ar AS created_user_last_name_ar,
        	created_user.initial AS created_user_initial,
        	created_user.color AS created_user_color,
        	created_user.email AS created_user_email,
        	created_user.profile_image_path AS created_user_profile_image_path,
			updated_user.id AS updated_user_id,
    		updated_user.first_name_en AS updated_user_first_name_en,
        	updated_user.first_name_ar AS updated_user_first_name_ar,
        	updated_user.last_name_en AS updated_user_last_name_en,
        	updated_user.last_name_ar AS updated_user_last_name_ar,
			updated_user.initial AS updated_user_initial,
        	updated_user.color AS updated_user_color,
        	updated_user.email AS updated_user_email,
        	updated_user.profile_image_path AS updated_user_profile_image_path
        `).
		Joins(`
            LEFT JOIN risk_action_plans ON risk_action_plans.id = risk_action_plan_status_updates.risk_action_plan_id
            LEFT JOIN action_plan_statuses ON action_plan_statuses.id = risk_action_plan_status_updates.action_plan_status_id
            LEFT JOIN users AS created_user ON created_user.id = risk_action_plan_status_updates.created_by_id
            LEFT JOIN users AS updated_user ON updated_user.id = risk_action_plan_status_updates.updated_by_id
        `).
		Where("risk_action_plan_status_updates.risk_action_plan_id = ?", riskActionPlanID).
		Order("risk_action_plan_status_updates.created_at DESC").
		First(&riskActionPlanUpdate)

	if query.Error != nil {
		return nil, query.Error
	}

	return riskActionPlanUpdate, nil
}

func (repo *AssetRiskActionPlanStatusUpdateRepository) RiskActionPlanUpdateStatus(riskActionPlan *models.RiskActionPlan, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(riskActionPlan).Error; err != nil {
		return err
	}
	return nil
}
