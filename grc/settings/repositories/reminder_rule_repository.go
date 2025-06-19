package repositories

import (
	"jethur/grc/models"
	orgDTOs "jethur/grc/settings/dtos"

	"gorm.io/gorm"
)

// Interface
type ReminderRuleRepositoryInterface interface {
	GetAll(params *orgDTOs.ReminderRuleListQueryParamsDTO) ([]*orgDTOs.ReminderRuleListDTO, error)
	GetByID(id uint) (*models.ReminderRule, error)
	GetFullDetailsByID(id uint) (*models.ReminderRule, error)
	Create(item *models.ReminderRule) (uint, error)
	Update(updatedReminderRule *models.ReminderRule) error
	Delete(id uint) error

	GetReminderRulesByReminderID(reminderID uint) ([]models.ReminderRule, error)
}

// Implementation
type ReminderRuleRepository struct {
	db *gorm.DB
}

func NewReminderRuleRepository(db *gorm.DB) *ReminderRuleRepository {
	return &ReminderRuleRepository{
		db: db,
	}
}

func (repo *ReminderRuleRepository) GetAll(params *orgDTOs.ReminderRuleListQueryParamsDTO) ([]*orgDTOs.ReminderRuleListDTO, error) {
	var result []*orgDTOs.ReminderRuleListDTO

	query := repo.db.Table("reminder_rules").
		Select(`
			reminder_rules.*,
			reminders.title_en AS reminder_title_en,
			reminders.title_ar AS reminder_title_ar,
			reminders.type AS reminder_type,
			trigger_conditions.type AS trigger_condition_type,
			trigger_conditions.title_en AS trigger_condition_title_en,
			trigger_conditions.title_ar AS trigger_condition_title_ar,
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
		Joins(`LEFT JOIN reminders ON reminders.id = reminder_rules.reminder_id`).
		Joins(`LEFT JOIN trigger_conditions ON trigger_conditions.id = reminder_rules.trigger_condition_id`).
		Joins("LEFT JOIN users AS created_user ON created_user.id = reminder_rules.created_by_id").
		Joins("LEFT JOIN users AS updated_user ON updated_user.id = reminder_rules.updated_by_id")

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			reminder_rules.title_en LIKE ? OR
			reminder_rules.title_ar LIKE ? OR
			trigger_conditions.title_en LIKE ? OR
			trigger_conditions.title_ar LIKE ? OR
			reminders.title_en LIKE ? OR
			reminders.title_ar LIKE ?
		`, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	if params.SortBy == "" {
		params.SortBy = "reminder_rules.title_en"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	err := query.Find(&result).Error
	return result, err
}

func (repo *ReminderRuleRepository) GetByID(id uint) (*models.ReminderRule, error) {
	var reminderRule *models.ReminderRule
	err := repo.db.First(&reminderRule, id).Error
	return reminderRule, err
}

func (repo *ReminderRuleRepository) GetFullDetailsByID(id uint) (*models.ReminderRule, error) {

	var reminderRule *models.ReminderRule

	result := repo.db.
		Preload("Reminder").
		Preload("TriggerCondition").
		First(&reminderRule, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return reminderRule, nil
}

func (repo *ReminderRuleRepository) Create(reminderRule *models.ReminderRule) (uint, error) {
	err := repo.db.Create(reminderRule).Error
	return reminderRule.ID, err
}

func (repo *ReminderRuleRepository) Update(updatedReminderRule *models.ReminderRule) error {

	if err := repo.db.Save(updatedReminderRule).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ReminderRuleRepository) Delete(id uint) error {
	return repo.db.Delete(&models.ReminderRule{}, id).Error
}

func (repo *ReminderRuleRepository) GetReminderRulesByReminderID(reminderID uint) ([]models.ReminderRule, error) {
	var reminderRule []models.ReminderRule

	err := repo.db.
		Where("reminder_id = ?", reminderID).
		Find(&reminderRule).Error
	if err != nil {
		return nil, err
	}
	return reminderRule, nil
}
