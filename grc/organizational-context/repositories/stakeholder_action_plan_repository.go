package repositories

import (
	"jethur/grc/models"
	organizationalContextDTOs "jethur/grc/organizational-context/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type StakeholderActionPlanRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(stakeholderID uint, params *organizationalContextDTOs.StakeholderActionPlanQueryParamsDTO) ([]*organizationalContextDTOs.StakeholderActionPlanListDTO, error)
	GetByID(id uint) (*models.StakeholderActionPlan, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.StakeholderActionPlan, error)
	Create(actionPlan *models.StakeholderActionPlan, tx ...utility.DatabaseTransaction) (uint, error)
	Update(actionPlan *models.StakeholderActionPlan, tx ...utility.DatabaseTransaction) error
	Delete(actionPlanID uint, tx ...utility.DatabaseTransaction) error

	GetDocumentByID(id uint) (*models.StakeholderActionPlanDocument, error)
	CreateDocuments(docs []*models.StakeholderActionPlanDocument, tx ...utility.DatabaseTransaction) error
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error
}

type StakeholderActionPlanRepository struct {
	db *gorm.DB
}

func NewStakeholderActionPlanRepository(db *gorm.DB) *StakeholderActionPlanRepository {
	return &StakeholderActionPlanRepository{
		db: db,
	}
}

func (repo *StakeholderActionPlanRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *StakeholderActionPlanRepository) GetAll(stakeholderID uint, params *organizationalContextDTOs.StakeholderActionPlanQueryParamsDTO) ([]*organizationalContextDTOs.StakeholderActionPlanListDTO, error) {
	var stakeholderActionPlans []*organizationalContextDTOs.StakeholderActionPlanListDTO

	query := repo.db.Table("stakeholder_action_plans").
		Select(`
            stakeholder_action_plans.*,
            severity_levels.title_en AS severity_level_title_en, 
            severity_levels.title_ar AS severity_level_title_ar, 
			severity_levels.color_code AS severity_level_color_code,
            organization_units.title_en AS responsible_unit_title_en, 
            organization_units.title_ar AS responsible_unit_title_ar,
            responsible_user.first_name_en AS responsible_user_first_name_en,
            responsible_user.last_name_en AS responsible_user_last_name_en, 
            responsible_user.first_name_ar AS responsible_user_first_name_ar,
            responsible_user.last_name_ar AS responsible_user_last_name_ar,
            responsible_user.initial AS responsible_user_initial, 
            responsible_user.email AS responsible_user_email, 
            responsible_user.color AS responsible_user_color, 
            responsible_user.profile_image_path AS responsible_user_profile_image_path, 
			assigned_owner.id AS assigned_owner_id,
			assigned_owner.first_name_en AS assigned_owner_first_name_en,
			assigned_owner.first_name_ar AS assigned_owner_first_name_ar,
			assigned_owner.last_name_en AS assigned_owner_last_name_en,
			assigned_owner.last_name_ar AS assigned_owner_last_name_ar,
			assigned_owner.initial AS assigned_owner_initial,
			assigned_owner.color AS assigned_owner_color,
			assigned_owner.email AS assigned_owner_email,
			assigned_owner.profile_image_path AS assigned_owner_profile_image_path,
            CONVERT(VARCHAR, stakeholder_action_plans.target_date, 23) AS target_date,
            action_plan_statuses.title_en AS action_plan_status_title_en, 
            action_plan_statuses.title_ar AS action_plan_status_title_ar, 
			action_plan_statuses.color_code AS action_plan_status_color_code, 
			action_plan_statuses.type AS action_plan_status_type, 
            created_user.first_name_en AS created_user_first_name_en,
            created_user.last_name_en AS created_user_last_name_en, 
            created_user.first_name_ar AS created_user_first_name_ar,
            created_user.last_name_ar AS created_user_last_name_ar,
            created_user.profile_image_path AS created_user_profile_image_path, 
            created_user.email AS created_user_email, 
            created_user.initial AS created_user_initial, 
            created_user.color AS created_user_color, 
            updated_user.first_name_en AS updated_user_first_name_en, 
            updated_user.last_name_en AS updated_user_last_name_en,
            updated_user.first_name_ar AS updated_user_first_name_ar,
            updated_user.last_name_ar AS updated_user_last_name_ar,
            updated_user.profile_image_path AS updated_user_profile_image_path, 
            updated_user.email AS updated_user_email,
            updated_user.initial AS updated_user_initial,
            updated_user.color AS updated_user_color 
        `).
		Joins(`LEFT JOIN severity_levels ON severity_levels.id = stakeholder_action_plans.severity_level_id`).
		Joins(`LEFT JOIN users AS responsible_user ON responsible_user.id = stakeholder_action_plans.responsible_user_id`).
		Joins("LEFT JOIN users AS assigned_owner ON assigned_owner.id = stakeholder_action_plans.assigned_owner_id").
		Joins(`LEFT JOIN organization_units ON organization_units.id = stakeholder_action_plans.responsible_unit_id`).
		Joins(`LEFT JOIN action_plan_statuses ON action_plan_statuses.id = stakeholder_action_plans.action_plan_status_id`).
		Joins(`LEFT JOIN users AS created_user ON created_user.id = stakeholder_action_plans.created_by_id`).
		Joins(`LEFT JOIN users AS updated_user ON updated_user.id = stakeholder_action_plans.updated_by_id`).
		Where("stakeholder_action_plans.stakeholder_id = ?", stakeholderID)

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}
	// Apply quick filters
	for _, filter := range params.QuickFilter {
		switch filter {
		case "is_recovery":
			query = query.Where("stakeholder_action_plans.is_recovery = ?", true)
		case "is_not_recovery":
			query = query.Where("stakeholder_action_plans.is_recovery = ?", false)
		}
	}
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(organization_units.title_en LIKE ? OR
			organization_units.title_ar LIKE ? OR
			responsible_user.first_name_en LIKE ? OR
			responsible_user.last_name_en LIKE ? OR
			responsible_user.first_name_ar LIKE ? OR
			responsible_user.last_name_ar LIKE ? OR
			responsible_user.email LIKE ? OR
			action_plan_statuses.title_en LIKE ? OR
			action_plan_statuses.title_ar LIKE ? OR
			stakeholder_action_plans.title LIKE ? OR
			created_user.first_name_en LIKE ? OR
			created_user.last_name_en LIKE ? OR
			created_user.first_name_ar LIKE ? OR
			created_user.last_name_ar LIKE ? OR
			created_user.email LIKE ? OR
			updated_user.first_name_en LIKE ? OR
			updated_user.last_name_en LIKE ? OR
			updated_user.first_name_ar LIKE ? OR
			updated_user.last_name_ar LIKE ? OR
			updated_user.email LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword,
		)
	}

	if params.SortBy != "" && params.Sort != "" {
		query = query.Order(params.SortBy + " " + params.Sort)
	}

	if err := query.Find(&stakeholderActionPlans).Error; err != nil {
		return nil, err
	}

	return stakeholderActionPlans, nil
}

func (repo *StakeholderActionPlanRepository) GetByID(id uint) (*models.StakeholderActionPlan, error) {
	var actionPlan models.StakeholderActionPlan
	if err := repo.db.First(&actionPlan, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &actionPlan, nil
}

func (repo *StakeholderActionPlanRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.StakeholderActionPlan, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var actionPlan *models.StakeholderActionPlan
	result := db.
		Preload("Stakeholder").
		Preload("SeverityLevel").
		Preload("ActionPlanStatus").
		Preload("ResponsibleUser").
		Preload("ResponsibleUnit").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("Documents").
		Preload("AssignedOwner").
		Where("id = ?", id).
		First(&actionPlan)

	if result.Error != nil {
		return nil, result.Error
	}

	return actionPlan, nil
}

func (repo *StakeholderActionPlanRepository) Create(actionPlan *models.StakeholderActionPlan, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(actionPlan).Error; err != nil {
		return 0, err
	}

	return actionPlan.ID, nil
}

func (repo *StakeholderActionPlanRepository) Update(actionPlan *models.StakeholderActionPlan, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Model(&models.StakeholderActionPlan{}).Where("id = ?", actionPlan.ID).Updates(actionPlan).Error; err != nil {
		return err
	}

	return nil
}

func (repo *StakeholderActionPlanRepository) Delete(actionPlanID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id = ?", actionPlanID).Delete(&models.StakeholderActionPlan{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *StakeholderActionPlanRepository) GetDocumentByID(id uint) (*models.StakeholderActionPlanDocument, error) {
	var document models.StakeholderActionPlanDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (repo *StakeholderActionPlanRepository) CreateDocuments(docs []*models.StakeholderActionPlanDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&docs).Error; err != nil {
		return err
	}
	return nil
}

func (repo *StakeholderActionPlanRepository) DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Where("id IN (?)", docIds).Delete(&models.StakeholderActionPlanDocument{}).Error; err != nil {
		return err
	}
	return nil
}
