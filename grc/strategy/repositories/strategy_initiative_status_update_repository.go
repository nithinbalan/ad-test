package repositories

import (
	"jethur/grc/models"
	strategyDTOs "jethur/grc/strategy/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type StrategyInitiativeStatusUpdateRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(strategyInitiativeID uint, sort string, sortBy string, keyword string) ([]*strategyDTOs.StrategyInitiativeStatusUpdateListDTO, error)
	GetByID(id uint) (*models.StrategicInitiativeStatusUpdate, error)
	Create(strategyInitiativeStatusUpdate *models.StrategicInitiativeStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error)
	Update(strategyInitiativeStatusUpdate *models.StrategicInitiativeStatusUpdate, tx ...utility.DatabaseTransaction) error
	GetLatestStrategyInitiativeStatusUpdate(strategyInitiativeID uint) (*strategyDTOs.StrategyInitiativeStatusUpdateListDTO, error)
}

type StrategyInitiativeStatusUpdateRepository struct {
	db *gorm.DB
}

func NewStrategyInitiativeStatusUpdateRepository(db *gorm.DB) *StrategyInitiativeStatusUpdateRepository {
	return &StrategyInitiativeStatusUpdateRepository{
		db: db,
	}
}

func (repo *StrategyInitiativeStatusUpdateRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *StrategyInitiativeStatusUpdateRepository) GetAll(strategyInitiativeID uint, sort string, sortBy string, keyword string) ([]*strategyDTOs.StrategyInitiativeStatusUpdateListDTO, error) {
	var strategyInitiativeStatusUpdates []*strategyDTOs.StrategyInitiativeStatusUpdateListDTO

	query := repo.db.Table("strategic_initiative_status_updates").
		Select(`
		strategic_initiative_status_updates.*,
		strategic_intiatives.id AS strategic_initiative_id,
		strategic_intiatives.title_en AS strategic_initiative_title,
		strategic_initiative_statuses.title_en AS strategic_initiative_status_title_en,
		strategic_initiative_statuses.title_ar AS strategic_initiative_status_title_ar,
		strategic_initiative_statuses.color_code AS strategic_initiative_status_color_code,
		strategic_initiative_statuses.type AS strategic_initiative_status_type,
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
	`).Joins(`
		LEFT JOIN strategic_intiatives ON strategic_intiatives.id = strategic_initiative_status_updates.strategy_initiative_id
		LEFT JOIN strategic_initiative_statuses ON strategic_initiative_statuses.id = strategic_initiative_status_updates.strategic_initiative_status_id
		LEFT JOIN users AS created_user ON created_user.id = strategic_initiative_status_updates.created_by_id
		LEFT JOIN users AS updated_user ON updated_user.id = strategic_initiative_status_updates.updated_by_id
		`).
		Where("strategic_initiative_status_updates.strategy_initiative_id = ?", strategyInitiativeID)

	if keyword != "" {
		query = query.Where(`
			(strategic_intiatives.title_en LIKE ? OR
			strategic_intiatives.title_ar LIKE ? OR
			strategic_initiative_statuses.title_en LIKE ? OR
			strategic_initiative_statuses.title_ar LIKE ? OR
			strategic_initiative_status_updates.comment LIKE ?)`,
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if sortBy == "" {
		sortBy = "strategic_initiative_status_updates.id"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}
	query = query.Order(sortBy + " " + sort)
	result := query.Find(&strategyInitiativeStatusUpdates)

	return strategyInitiativeStatusUpdates, result.Error
}

func (repo *StrategyInitiativeStatusUpdateRepository) GetByID(id uint) (*models.StrategicInitiativeStatusUpdate, error) {
	var strategyInitiativeStatusUpdate *models.StrategicInitiativeStatusUpdate
	result := repo.db.First(&strategyInitiativeStatusUpdate, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return strategyInitiativeStatusUpdate, nil
}

func (repo *StrategyInitiativeStatusUpdateRepository) Create(strategyInitiativeStatusUpdate *models.StrategicInitiativeStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(strategyInitiativeStatusUpdate).Error; err != nil {
		return 0, err
	}

	return strategyInitiativeStatusUpdate.ID, nil
}

func (repo *StrategyInitiativeStatusUpdateRepository) Update(strategyInitiativeStatusUpdate *models.StrategicInitiativeStatusUpdate, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(strategyInitiativeStatusUpdate).Error; err != nil {
		return err
	}
	return nil
}

func (repo *StrategyInitiativeStatusUpdateRepository) GetLatestStrategyInitiativeStatusUpdate(strategyInitiativeID uint) (*strategyDTOs.StrategyInitiativeStatusUpdateListDTO, error) {
	var strategyInitiativeStatusUpdate *strategyDTOs.StrategyInitiativeStatusUpdateListDTO

	query := repo.db.Table("strategic_initiative_status_updates").
		Select(`
		strategic_initiative_status_updates.*,
		strategic_intiatives.title_en AS strategic_initiative_title_en,
		strategic_intiatives.title_ar AS strategic_initiative_title_ar,
		strategic_initiative_statuses.title_en AS strategic_initiative_status_title_en,
		strategic_initiative_statuses.title_ar AS strategic_initiative_status_title_ar,
		strategic_initiative_statuses.color_code AS strategic_initiative_status_color_code,
		strategic_initiative_statuses.type AS strategic_initiative_status_type,
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
	`).Joins(`
		LEFT JOIN strategic_intiatives ON strategic_intiatives.id = strategic_initiative_status_updates.strategy_initiative_id
		LEFT JOIN strategic_initiative_statuses ON strategic_initiative_statuses.id = strategic_initiative_status_updates.strategic_initiative_status_id
		LEFT JOIN users AS created_user ON created_user.id = strategic_initiative_status_updates.created_by_id
		LEFT JOIN users AS updated_user ON updated_user.id = strategic_initiative_status_updates.updated_by_id
		`).
		Where("strategic_initiative_status_updates.strategy_initiative_id = ?", strategyInitiativeID).
		Order("strategic_initiative_status_updates.created_at DESC").
		First(&strategyInitiativeStatusUpdate)

	if query.Error != nil {
		return nil, query.Error
	}

	return strategyInitiativeStatusUpdate, nil
}
