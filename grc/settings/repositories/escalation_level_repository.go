package repositories

import (
	"jethur/grc/models"
	orgDTOs "jethur/grc/settings/dtos"

	"gorm.io/gorm"
)

// Interface
type EscalationLevelRepositoryInterface interface {
	GetAll(params *orgDTOs.EscalationLevelListQueryParamsDTO) ([]*orgDTOs.EscalationLevelListDTO, error)
	GetByID(id uint) (*models.EscalationLevel, error)
	GetFullDetailsByID(id uint) (*models.EscalationLevel, error)
	Create(item *models.EscalationLevel) (uint, error)
	Update(updatedEscalationLevel *models.EscalationLevel) error
	Delete(id uint) error

	GetEscalationLevelsByEscalationID(escalationID uint) ([]models.EscalationLevel, error)
}

// Implementation
type EscalationLevelRepository struct {
	db *gorm.DB
}

func NewEscalationLevelRepository(db *gorm.DB) *EscalationLevelRepository {
	return &EscalationLevelRepository{
		db: db,
	}
}

func (repo *EscalationLevelRepository) GetAll(params *orgDTOs.EscalationLevelListQueryParamsDTO) ([]*orgDTOs.EscalationLevelListDTO, error) {
	var result []*orgDTOs.EscalationLevelListDTO

	query := repo.db.Table("escalation_levels").
		Select(`
			escalation_levels.*,
			escalations.title_en AS escalation_title_en,
			escalations.title_ar AS escalation_title_ar,
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
		Joins(`LEFT JOIN escalations ON escalations.id = escalation_levels.escalation_id`).
		Joins("LEFT JOIN users AS created_user ON created_user.id = escalation_levels.created_by_id").
		Joins("LEFT JOIN users AS updated_user ON updated_user.id = escalation_levels.updated_by_id")

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			escalations.title_en LIKE ? OR
			escalations.title_ar LIKE ?
		`, keyword, keyword, keyword, keyword)
	}

	if params.SortBy == "" {
		params.SortBy = "escalation_levels.day_offset"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	err := query.Find(&result).Error
	return result, err
}

func (repo *EscalationLevelRepository) GetByID(id uint) (*models.EscalationLevel, error) {
	var escalationLevel *models.EscalationLevel
	err := repo.db.First(&escalationLevel, id).Error
	return escalationLevel, err
}

func (repo *EscalationLevelRepository) GetFullDetailsByID(id uint) (*models.EscalationLevel, error) {

	var escalationLevel *models.EscalationLevel

	result := repo.db.
		Preload("Escalation").
		Preload("Designation").
		First(&escalationLevel, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return escalationLevel, nil
}

func (repo *EscalationLevelRepository) Create(escalationLevel *models.EscalationLevel) (uint, error) {
	err := repo.db.Create(escalationLevel).Error
	return escalationLevel.ID, err
}

func (repo *EscalationLevelRepository) Update(updatedEscalationLevel *models.EscalationLevel) error {

	if err := repo.db.Save(updatedEscalationLevel).Error; err != nil {
		return err
	}

	return nil
}

func (repo *EscalationLevelRepository) Delete(id uint) error {
	return repo.db.Delete(&models.EscalationLevel{}, id).Error
}

func (repo *EscalationLevelRepository) GetEscalationLevelsByEscalationID(escalationID uint) ([]models.EscalationLevel, error) {
	var escalationLevel []models.EscalationLevel

	err := repo.db.
		Where("escalation_id = ?", escalationID).
		Find(&escalationLevel).Error
	if err != nil {
		return nil, err
	}
	return escalationLevel, nil
}
