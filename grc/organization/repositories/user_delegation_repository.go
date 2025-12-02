package repositories

import (
	"jethur/grc/models"
	orgDTOs "jethur/grc/organization/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interface
type UserDelegationRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(userID uint) ([]*orgDTOs.UserDelegationHistoryListDTO, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.UserDelegationHistory, error)
	Create(item *models.UserDelegationHistory, tx ...utility.DatabaseTransaction) (uint, error)

	GetLatestUserDelegation(userID uint) (*orgDTOs.UserDelegationHistoryListDTO, error)
	UpdateActionPlanDelegationUser(responsibleUserID uint, tables []string, value *uint, tx ...utility.DatabaseTransaction) error
}

// Implementation
type UserDelegationRepository struct {
	db *gorm.DB
}

func NewUserDelegationRepository(db *gorm.DB) *UserDelegationRepository {
	return &UserDelegationRepository{
		db: db,
	}
}

func (repo *UserDelegationRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *UserDelegationRepository) GetAll(userID uint) ([]*orgDTOs.UserDelegationHistoryListDTO, error) {
	var userDelegations []*orgDTOs.UserDelegationHistoryListDTO

	query := repo.db.Table("user_delegation_history").
		Select(`
		user_delegation_history.*,
		CONVERT(VARCHAR, user_delegation_history.assigned_date, 23) AS assigned_date,
		delegatee_user.id AS delegatee_user_id,
		delegatee_user.first_name_en AS delegatee_user_first_name_en,
		delegatee_user.first_name_ar AS delegatee_user_first_name_ar,
		delegatee_user.last_name_en AS delegatee_user_last_name_en,
		delegatee_user.last_name_ar AS delegatee_user_last_name_ar,
		delegatee_user.initial AS delegatee_user_initial,
		delegatee_user.color AS delegatee_user_color,
		delegatee_user.email AS delegatee_user_email,
		delegatee_user.profile_image_path AS delegatee_user_profile_image_path,
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
		Joins("LEFT JOIN users AS delegatee_user ON delegatee_user.id = user_delegation_history.delegatee_user_id").
		Joins("LEFT JOIN users AS created_user ON created_user.id = user_delegation_history.created_by_id").
		Joins("LEFT JOIN users AS updated_user ON updated_user.id = user_delegation_history.updated_by_id").
		Where("user_delegation_history.user_id = ?", userID)

	err := query.Find(&userDelegations).Error
	if err != nil {
		return nil, err
	}

	return userDelegations, nil
}

func (repo *UserDelegationRepository) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.UserDelegationHistory, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var userDelegationHistory *models.UserDelegationHistory
	result := db.First(&userDelegationHistory, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return userDelegationHistory, nil
}

func (repo *UserDelegationRepository) Create(userDelegationHistory *models.UserDelegationHistory, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(userDelegationHistory).Error; err != nil {
		return 0, err
	}
	return userDelegationHistory.ID, nil
}

func (repo *UserDelegationRepository) GetLatestUserDelegation(userID uint) (*orgDTOs.UserDelegationHistoryListDTO, error) {
	var userDelegations *orgDTOs.UserDelegationHistoryListDTO

	query := repo.db.Table("user_delegation_history").
		Select(`
		user_delegation_history.*,
		CONVERT(VARCHAR, user_delegation_history.assigned_date, 23) AS assigned_date,
		delegatee_user.id AS delegatee_user_id,
		delegatee_user.first_name_en AS delegatee_user_first_name_en,
		delegatee_user.first_name_ar AS delegatee_user_first_name_ar,
		delegatee_user.last_name_en AS delegatee_user_last_name_en,
		delegatee_user.last_name_ar AS delegatee_user_last_name_ar,
		delegatee_user.initial AS delegatee_user_initial,
		delegatee_user.color AS delegatee_user_color,
		delegatee_user.email AS delegatee_user_email,
		delegatee_user.profile_image_path AS delegatee_user_profile_image_path,
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
		Joins("LEFT JOIN users AS delegatee_user ON delegatee_user.id = user_delegation_history.delegatee_user_id").
		Joins("LEFT JOIN users AS created_user ON created_user.id = user_delegation_history.created_by_id").
		Joins("LEFT JOIN users AS updated_user ON updated_user.id = user_delegation_history.updated_by_id").
		Where("user_delegation_history.user_id = ?", userID).
		Order("user_delegation_history.created_at DESC").
		First(&userDelegations)

	if query.Error != nil {
		return nil, query.Error
	}

	return userDelegations, nil
}

func (repo *UserDelegationRepository) UpdateActionPlanDelegationUser(responsibleUserID uint, tables []string, value *uint, tx ...utility.DatabaseTransaction) error {

	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Transaction(func(tx *gorm.DB) error {
		for _, table := range tables {
			if table != "document_change_request_workflow_participants" && table != "document_workflow_participants" && table != "risk_workflow_participants" {
				if err := tx.Table(table).
					Where("responsible_user_id = ? ", responsibleUserID).
					Update("delegatee_user_id", value).Error; err != nil {
					return err
				}
			} else {
				if err := tx.Table(table).
					Where("user_id = ? ", responsibleUserID).
					Update("delegatee_user_id", value).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}
